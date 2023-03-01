// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package core_test

import (
	"context"
	"math/big"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"pkg.berachain.dev/stargazer/eth/common"
	"pkg.berachain.dev/stargazer/eth/core"
	"pkg.berachain.dev/stargazer/eth/core/mock"
	"pkg.berachain.dev/stargazer/eth/core/types"
	"pkg.berachain.dev/stargazer/eth/core/vm"
	vmmock "pkg.berachain.dev/stargazer/eth/core/vm/mock"
	"pkg.berachain.dev/stargazer/eth/crypto"
	"pkg.berachain.dev/stargazer/eth/params"
	"pkg.berachain.dev/stargazer/eth/testutil/contracts/solidity/generated"
)

var (
	dummyContract = common.HexToAddress("0x9fd0aA3B78277a1E717de9D3de434D4b812e5499")
	key, _        = crypto.GenerateEthKey()
	signer        = types.LatestSignerForChainID(params.DefaultChainConfig.ChainID)

	legacyTxData = &types.LegacyTx{
		Nonce:    0,
		To:       &dummyContract,
		Gas:      100000,
		GasPrice: big.NewInt(2),
		Data:     []byte("abcdef"),
	}
)

var _ = Describe("StateProcessor", func() {
	var (
		sdb           *vmmock.StargazerStateDBMock
		host          *mock.StargazerHostChainMock
		bp            *mock.BlockPluginMock
		gp            *mock.GasPluginMock
		cp            *mock.ConfigurationPluginMock
		pp            *mock.PrecompilePluginMock
		sp            *core.StateProcessor
		blockNumber   uint64
		blockGasLimit uint64
	)

	BeforeEach(func() {
		sdb = vmmock.NewEmptyStateDB()
		host = mock.NewMockHost()
		bp = mock.NewBlockPluginMock()
		gp = mock.NewGasPluginMock()
		cp = mock.NewConfigurationPluginMock()
		pp = mock.NewPrecompilePluginMock()
		host.GetBlockPluginFunc = func() core.BlockPlugin {
			return bp
		}
		host.GetGasPluginFunc = func() core.GasPlugin {
			return gp
		}
		host.GetConfigurationPluginFunc = func() core.ConfigurationPlugin {
			return cp
		}
		host.GetPrecompilePluginFunc = func() core.PrecompilePlugin {
			return pp
		}
		pp.RegisterFunc = func(pc vm.PrecompileContainer) error {
			return nil
		}
		sp = core.NewStateProcessor(host, sdb, vm.Config{}, true)
		Expect(sp).ToNot(BeNil())
		blockNumber = params.DefaultChainConfig.LondonBlock.Uint64() + 1
		blockGasLimit = 1000000

		bp.GetStargazerHeaderByNumberFunc = func(height int64) *types.StargazerHeader {
			header := types.NewEmptyStargazerHeader()
			header.GasLimit = blockGasLimit
			header.BaseFee = big.NewInt(1)
			header.Coinbase = common.BytesToAddress([]byte{2})
			header.Number = big.NewInt(int64(blockNumber))
			header.Time = uint64(3)
			header.Difficulty = new(big.Int)
			header.MixDigest = common.BytesToHash([]byte{})
			return header
		}
		pp.HasFunc = func(addr common.Address) bool {
			return false
		}

		sdb.PrepareFunc = func(rules params.Rules, sender common.Address,
			coinbase common.Address, dest *common.Address,
			precompiles []common.Address, txAccesses types.AccessList,
		) {
			// no-op
		}

		gp.SetBlockGasLimit(blockGasLimit)
		sp.Prepare(context.Background(), 0)
	})

	Context("Empty block", func() {
		It("should build a an empty block", func() {
			block, err := sp.Finalize(context.Background())
			Expect(err).To(BeNil())
			Expect(block).ToNot(BeNil())
			Expect(block.TxIndex()).To(Equal(0))
		})
	})

	Context("Block with transactions", func() {
		BeforeEach(func() {
			_, err := sp.Finalize(context.Background())
			Expect(err).To(BeNil())

			pp.ResetFunc = func(ctx context.Context) {
				// no-op
			}

			sp.Prepare(context.Background(), int64(blockNumber))
		})

		It("should error on an unsigned transaction", func() {
			receipt, err := sp.ProcessTransaction(context.Background(), types.NewTx(legacyTxData))
			Expect(err).ToNot(BeNil())
			Expect(receipt).To(BeNil())
			block, err := sp.Finalize(context.Background())
			Expect(err).To(BeNil())
			Expect(block).ToNot(BeNil())
			Expect(block.TxIndex()).To(Equal(0))
		})

		It("should not error on a signed transaction", func() {
			signedTx := types.MustSignNewTx(key, signer, legacyTxData)
			sdb.GetBalanceFunc = func(addr common.Address) *big.Int {
				return big.NewInt(200000)
			}
			result, err := sp.ProcessTransaction(context.Background(), signedTx)
			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())
			Expect(result.Status).To(Equal(types.ReceiptStatusSuccessful))
			Expect(result.BlockNumber).To(Equal(big.NewInt(int64(blockNumber))))
			Expect(result.TransactionIndex).To(Equal(uint(0)))
			Expect(result.TxHash.Hex()).To(Equal(signedTx.Hash().Hex()))
			Expect(result.GasUsed).ToNot(BeZero())
			block, err := sp.Finalize(context.Background())
			Expect(err).To(BeNil())
			Expect(block).ToNot(BeNil())
			Expect(block.TxIndex()).To(Equal(1))
		})

		It("should add a contract address to the receipt", func() {
			legacyTxDataCopy := *legacyTxData
			legacyTxDataCopy.To = nil
			sdb.GetBalanceFunc = func(addr common.Address) *big.Int {
				return big.NewInt(200000)
			}
			signedTx := types.MustSignNewTx(key, signer, &legacyTxDataCopy)
			result, err := sp.ProcessTransaction(context.Background(), signedTx)
			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())
			Expect(result.ContractAddress).ToNot(BeNil())
			block, err := sp.Finalize(context.Background())
			Expect(err).To(BeNil())
			Expect(block).ToNot(BeNil())
			Expect(block.TxIndex()).To(Equal(1))
		})

		It("should mark a receipt with a virtual machine error as failed", func() {
			sdb.GetBalanceFunc = func(addr common.Address) *big.Int {
				return big.NewInt(200000)
			}
			sdb.GetCodeFunc = func(addr common.Address) []byte {
				if addr != dummyContract {
					return nil
				}
				return common.Hex2Bytes(generated.NonRevertableTxMetaData.Bin)
			}
			sdb.GetCodeHashFunc = func(addr common.Address) common.Hash {
				if addr != dummyContract {
					return common.Hash{}
				}
				return crypto.Keccak256Hash(common.Hex2Bytes(generated.NonRevertableTxMetaData.Bin))
			}
			sdb.ExistFunc = func(addr common.Address) bool {
				return addr == dummyContract
			}
			legacyTxData.To = nil
			legacyTxData.Value = big.NewInt(0)
			signedTx := types.MustSignNewTx(key, signer, legacyTxData)
			result, err := sp.ProcessTransaction(context.Background(), signedTx)
			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())
			Expect(result.Status).To(Equal(types.ReceiptStatusSuccessful))
			block, err := sp.Finalize(context.Background())
			Expect(err).To(BeNil())
			Expect(block).ToNot(BeNil())
			Expect(block.TxIndex()).To(Equal(1))

			// Now try calling the contract
			legacyTxData.To = &dummyContract
			signedTx = types.MustSignNewTx(key, signer, legacyTxData)
			result, err = sp.ProcessTransaction(context.Background(), signedTx)
			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())
			Expect(result.Status).To(Equal(types.ReceiptStatusSuccessful))
		})
	})
})

var _ = Describe("No precompile plugin provided", func() {
	It("should use the default plugin if none is provided", func() {
		host := mock.NewMockHost()
		bp := mock.NewBlockPluginMock()
		gp := mock.NewGasPluginMock()
		gp.SetBlockGasLimit(1000000)
		bp.GetStargazerHeaderByNumberFunc = func(height int64) *types.Header {
			header := &types.Header{}
			header.GasLimit = 1000000
			header.Number = new(big.Int)
			header.Difficulty = new(big.Int)
			return header
		}
		host.GetBlockPluginFunc = func() core.BlockPlugin {
			return bp
		}
		host.GetGasPluginFunc = func() core.GasPlugin {
			return gp
		}
		host.GetConfigurationPluginFunc = func() core.ConfigurationPlugin {
			return mock.NewConfigurationPluginMock()
		}
		host.GetPrecompilePluginFunc = func() core.PrecompilePlugin {
			return nil
		}
		sp := core.NewStateProcessor(host, vmmock.NewEmptyStateDB(), vm.Config{}, true)
		Expect(func() { sp.Prepare(context.Background(), 0) }).ToNot(Panic())
	})
})

var _ = Describe("GetHashFn", func() {
	var (
		sdb           *vmmock.StargazerStateDBMock
		host          *mock.StargazerHostChainMock
		bp            *mock.BlockPluginMock
		gp            *mock.GasPluginMock
		cp            *mock.ConfigurationPluginMock
		pp            *mock.PrecompilePluginMock
		sp            *core.StateProcessor
		blockGasLimit uint64
	)

	BeforeEach(func() {
		sdb = vmmock.NewEmptyStateDB()
		host = mock.NewMockHost()
		bp = mock.NewBlockPluginMock()
		gp = mock.NewGasPluginMock()
		cp = mock.NewConfigurationPluginMock()
		pp = mock.NewPrecompilePluginMock()
		host.GetBlockPluginFunc = func() core.BlockPlugin {
			return bp
		}
		host.GetGasPluginFunc = func() core.GasPlugin {
			return gp
		}
		host.GetConfigurationPluginFunc = func() core.ConfigurationPlugin {
			return cp
		}
		host.GetPrecompilePluginFunc = func() core.PrecompilePlugin {
			return pp
		}
		pp.RegisterFunc = func(pc vm.PrecompileContainer) error {
			return nil
		}
		sp = core.NewStateProcessor(host, sdb, vm.Config{}, true)
		Expect(sp).ToNot(BeNil())
		blockGasLimit = 1000000

		bp.GetStargazerHeaderByNumberFunc = func(height int64) *types.StargazerHeader {
			return types.NewStargazerHeader(
				&types.Header{
					Number:     big.NewInt(height),
					BaseFee:    big.NewInt(69),
					GasLimit:   blockGasLimit,
					ParentHash: common.BytesToHash([]byte{uint8(height) - 1}),
					Time:       uint64(time.Now().Unix()),
					Difficulty: big.NewInt(0),
					MixDigest:  common.Hash{},
				},
				crypto.Keccak256Hash([]byte{byte(height)}),
			)
		}
		pp.HasFunc = func(addr common.Address) bool {
			return false
		}

		gp.SetBlockGasLimit(blockGasLimit)
	})

	It("should return empty hash", func() {
		sp.Prepare(context.Background(), 100)
		hashFn := sp.GetHashFn()
		Expect(hashFn(100)).To(Equal(common.Hash{}))

		_, err := sp.Finalize(context.Background())
		Expect(err).To(BeNil())

		sp.Prepare(context.Background(), 100)
		hashFn = sp.GetHashFn()
		Expect(hashFn(101)).To(Equal(common.Hash{}))
	})

	It("should return correct hash", func() {
		sp.Prepare(context.Background(), 100)
		hashFn := sp.GetHashFn()
		Expect(hashFn(99)).To(Equal(common.BytesToHash([]byte{99})))

		_, err := sp.Finalize(context.Background())
		Expect(err).To(BeNil())

		sp.Prepare(context.Background(), 101)
		hashFn = sp.GetHashFn()
		Expect(hashFn(99)).To(Equal(common.BytesToHash([]byte{99})))
	})
})
