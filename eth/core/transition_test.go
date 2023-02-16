// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// See the file LICENSE for licensing terms.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package core_test

import (
	"context"
	"math/big"

	"github.com/berachain/stargazer/eth/common"
	"github.com/berachain/stargazer/eth/core"
	"github.com/berachain/stargazer/eth/core/mock"
	"github.com/berachain/stargazer/eth/core/vm"
	vmmock "github.com/berachain/stargazer/eth/core/vm/mock"
	"github.com/berachain/stargazer/eth/params"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StateTransition", func() {
	var (
		evm *vmmock.StargazerEVMMock
		sdb *vmmock.StargazerStateDBMock
		msg mock.MessageMock
		gp  = mock.NewGasPluginMock(0)
	)

	BeforeEach(func() {
		msg = *mock.NewEmptyMessage()
		evm = vmmock.NewStargazerEVM()
		sdb, _ = evm.StateDB().(*vmmock.StargazerStateDBMock)
		msg.FromFunc = func() common.Address {
			return common.Address{1}
		}

		msg.GasPriceFunc = func() *big.Int {
			return big.NewInt(123456789)
		}

		msg.ToFunc = func() *common.Address {
			return &common.Address{1}
		}

		gp = mock.NewGasPluginMock(0)
	})

	When("Contract Creation", func() {
		BeforeEach(func() {
			msg.ToFunc = func() *common.Address {
				return nil
			}
		})
		It("should call create", func() {
			msg.GasFunc = func() uint64 {
				return 53000 // exact intrinsic gas for create after homestead
			}
			res, err := core.ApplyMessage(evm, gp, &msg, true)
			Expect(len(evm.CreateCalls())).To(Equal(1))
			Expect(res.UsedGas).To(Equal(uint64(53000)))
			Expect(err).To(BeNil())
		})
		When("we have less than the intrinsic gas", func() {
			msg.GasFunc = func() uint64 {
				return 53000 - 1
			}
			It("should return error", func() {
				_, err := core.ApplyMessage(evm, gp, &msg, true)
				Expect(err).To(MatchError(core.ErrIntrinsicGas))
			})
		})

		It("should call create with commit", func() {
			msg.GasFunc = func() uint64 {
				return 53000
			}
			_, err := core.ApplyMessage(evm, gp, &msg, true)
			Expect(err).To(BeNil())
		})

		It("should handle transition error", func() {
			msg.GasFunc = func() uint64 {
				return 0
			}
			_, err := core.ApplyMessage(evm, gp, &msg, true)
			Expect(err).To(Not(BeNil()))
		})

		When("We call with a tracer", func() {
			var tracer *vmmock.EVMLoggerMock
			BeforeEach(func() {
				tracer = vmmock.NewEVMLoggerMock()
				evm.ConfigFunc = func() vm.Config {
					return vm.Config{
						Debug:  true,
						Tracer: tracer,
					}
				}
			})

			It("should call create with tracer", func() {
				msg.GasFunc = func() uint64 {
					return 53000 // exact intrinsic gas for create after homestead
				}
				_, err := core.ApplyMessage(evm, gp, &msg, false)
				Expect(len(tracer.CaptureTxStartCalls())).To(Equal(1))
				Expect(len(tracer.CaptureTxEndCalls())).To(Equal(1))
				Expect(err).To(BeNil())
			})
			It("should call create with tracer and commit", func() {
				msg.GasFunc = func() uint64 {
					return 53000 // exact intrinsic gas for create after homestead
				}
				sdb = vmmock.NewEmptyStateDB()
				evm.StateDBFunc = func() vm.StargazerStateDB {
					return sdb
				}
				_, err := core.ApplyMessage(evm, gp, &msg, true)
				Expect(err).To(BeNil())
				Expect(len(tracer.CaptureTxStartCalls())).To(Equal(1))
				Expect(len(tracer.CaptureTxEndCalls())).To(Equal(1))
				Expect(len(sdb.FinalizeCalls())).To(Equal(1))
			})
			It("should handle abort error", func() {
				msg.GasFunc = func() uint64 {
					return 0
				}
				_, err := core.ApplyMessage(evm, gp, &msg, false)
				Expect(len(tracer.CaptureTxStartCalls())).To(Equal(1))
				Expect(len(tracer.CaptureTxEndCalls())).To(Equal(1))
				Expect(err).To(Not(BeNil()))
			})
			It("should handle abort error with commit", func() {
				msg.GasFunc = func() uint64 {
					return 0
				}
				_, err := core.ApplyMessage(evm, gp, &msg, true)
				Expect(len(tracer.CaptureTxStartCalls())).To(Equal(1))
				Expect(len(tracer.CaptureTxEndCalls())).To(Equal(1))
				Expect(err).To(Not(BeNil()))
			})
		})
	})

	When("Contract Call", func() {
		BeforeEach(func() {
			msg.ToFunc = func() *common.Address {
				return &common.Address{1}
			}

			sdb.GetCodeHashFunc = func(addr common.Address) common.Hash {
				return common.Hash{1}
			}
			msg.GasFunc = func() uint64 {
				return 100000
			}
		})

		Context("Gas Refund", func() {
			BeforeEach(func() {
				sdb.GetRefundFunc = func() uint64 {
					return 20000
				}
				evm.StateDBFunc = func() vm.StargazerStateDB {
					return sdb
				}
				evm.CallFunc = func(caller vm.ContractRef, addr common.Address,
					input []byte, gas uint64, value *big.Int) ([]byte, uint64, error) {
					return []byte{}, 80000, nil
				}
			})

			When("we are in london", func() {
				It("should call call", func() {
					res, err := core.ApplyMessage(evm, gp, &msg, true)
					Expect(len(evm.CallCalls())).To(Equal(1))
					Expect(res.UsedGas).To(Equal(uint64(16000))) // refund is capped to 1/5th
					Expect(err).To(BeNil())
				})
			})

			When("we are not in london", func() {
				It("should call and cap refund properly", func() {
					evm.ChainConfigFunc = func() *params.ChainConfig {
						return &params.ChainConfig{
							LondonBlock:    big.NewInt(1000000000),
							HomesteadBlock: big.NewInt(0),
						}
					}
					res, err := core.ApplyMessage(evm, gp, &msg, true)
					Expect(len(evm.CallCalls())).To(Equal(1))
					Expect(res.UsedGas).To(Equal(uint64(10000))) // refund is capped to 1/2
					Expect(err).To(BeNil())
				})
			})
		})
		It("should check to ensure required funds are available", func() {
			msg.ValueFunc = func() *big.Int {
				return big.NewInt(1)
			}
			evm.ContextFunc = func() vm.BlockContext {
				return vm.BlockContext{
					CanTransfer: func(db vm.GethStateDB, addr common.Address, amount *big.Int) bool {
						return false
					},
				}
			}
			_, err := core.ApplyMessage(evm, gp, &msg, true)
			Expect(err).To(MatchError(core.ErrInsufficientFundsForTransfer))
		})
		When("the message has data", func() {
			It("should cost more gas", func() {
				msg.GasFunc = func() uint64 {
					return 6969699669
				}

				msg.DataFunc = func() []byte {
					return []byte{1, 2, 3}
				}

				// Call the intrinsic gas function with data
				st := core.NewStateTransition(evm, gp, &msg)
				Expect(gp.SetTxGasLimit(10000000)).To(BeNil())
				Expect(st.ConsumeEthIntrinsicGas(true, true, true)).To(BeNil())
				consumedWithData := gp.CumulativeGasUsed()

				// Reset the gas meter.
				gp.Prepare(context.Background())

				// Call the intrinsic gas function with no data
				msg.DataFunc = func() []byte {
					return []byte{}
				}
				Expect(st.ConsumeEthIntrinsicGas(true, true, true)).To(BeNil())

				// We expect that the call with Data will consume more gas.
				Expect(consumedWithData).To(BeNumerically(">", gp.CumulativeGasUsed()))
			})
		})
	})
})
