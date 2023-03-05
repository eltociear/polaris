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

	"pkg.berachain.dev/stargazer/eth/common"
	"pkg.berachain.dev/stargazer/eth/core"
	"pkg.berachain.dev/stargazer/eth/core/mock"
	"pkg.berachain.dev/stargazer/eth/core/types"
	"pkg.berachain.dev/stargazer/eth/core/vm"
	vmmock "pkg.berachain.dev/stargazer/eth/core/vm/mock"
	"pkg.berachain.dev/stargazer/eth/crypto"
	"pkg.berachain.dev/stargazer/eth/params"
	"pkg.berachain.dev/stargazer/eth/testutil/contracts/solidity/generated"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	dummyContract = common.HexToAddress("0x9fd0aA3B78277a1E717de9D3de434D4b812e5499")
	key, _        = crypto.GenerateEthKey()
	signer        = types.LatestSignerForChainID(params.DefaultChainConfig.ChainID)
	_             = key
	_             = signer
	dummyHeader   = &types.Header{
		Number:   big.NewInt(1),
		GasLimit: 1000000,
	}
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
		sp = core.NewStateProcessor(host, sdb, &vm.Config{})
		Expect(sp).ToNot(BeNil())
		blockNumber = params.DefaultChainConfig.LondonBlock.Uint64() + 1
		blockGasLimit = 1000000

		bp.NewHeaderWithBlockNumberFunc = func(height int64) *types.Header {
			header := dummyHeader
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
		sp.Prepare(context.Background(), nil, dummyHeader)
	})

	Context("Empty block", func() {
		It("should build a an empty block", func() {
			block, receipts, err := sp.Finalize(context.Background())
			Expect(err).ToNot(HaveOccurred())
			Expect(block).ToNot(BeNil())
			Expect(receipts).To(BeEmpty())
		})
	})

	Context("Block with transactions", func() {
		BeforeEach(func() {
			_, _, err := sp.Finalize(context.Background())
			Expect(err).ToNot(HaveOccurred())

			sp.Prepare(context.Background(), nil, dummyHeader)
		})

		It("should error on an unsigned transaction", func() {
			receipt, err := sp.ProcessTransaction(context.Background(), types.NewTx(legacyTxData))
			Expect(err).To(HaveOccurred())
			Expect(receipt).To(BeNil())
			block, receipts, err := sp.Finalize(context.Background())
			Expect(err).ToNot(HaveOccurred())
			Expect(block).ToNot(BeNil())
			Expect(receipts).To(BeEmpty())
		})

		It("should not error on a signed transaction", func() {
			// signedTx := types.MustSignNewTx(key, signer, legacyTxData)
			// sdb.GetBalanceFunc = func(addr common.Address) *big.Int {
			// 	return big.NewInt(200000)
			// }
			// result, err := sp.ProcessTransaction(context.Background(), signedTx)
			// Expect(err).ToNot(HaveOccurred())
			// Expect(result).ToNot(BeNil())
			// Expect(result.Err).ToNot(HaveOccurred())
			// Expect(result.UsedGas).ToNot(BeZero())
			// block, receipts, err := sp.Finalize(context.Background())
			// Expect(err).ToNot(HaveOccurred())
			// Expect(block).ToNot(BeNil())
			// Expect(len(receipts)).To(Equal(1))
		})

		It("should handle", func() {
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
			// legacyTxData.To = nil
			// legacyTxData.Value = big.NewInt(0)
			// signedTx := types.MustSignNewTx(key, signer, legacyTxData)
			// result, err := sp.ProcessTransaction(context.Background(), signedTx)
			// Expect(err).ToNot(HaveOccurred())
			// Expect(result).ToNot(BeNil())
			// Expect(result.Err).ToNot(HaveOccurred())
			// block, receipts, err := sp.Finalize(context.Background())
			// Expect(err).ToNot(HaveOccurred())
			// Expect(block).ToNot(BeNil())
			// Expect(len(receipts)).To(Equal(1))

			// // Now try calling the contract
			// legacyTxData.To = &dummyContract
			// signedTx = types.MustSignNewTx(key, signer, legacyTxData)
			// result, err = sp.ProcessTransaction(context.Background(), signedTx)
			// Expect(err).ToNot(HaveOccurred())
			// Expect(result).ToNot(BeNil())
			// Expect(result.Err).ToNot(HaveOccurred())
		})
	})
})

var _ = Describe("No precompile plugin provided", func() {
	It("should use the default plugin if none is provided", func() {
		host := mock.NewMockHost()
		bp := mock.NewBlockPluginMock()
		gp := mock.NewGasPluginMock()
		gp.SetBlockGasLimit(1000000)
		bp.NewHeaderWithBlockNumberFunc = func(height int64) *types.Header {
			header := dummyHeader
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
		sp := core.NewStateProcessor(host, vmmock.NewEmptyStateDB(), &vm.Config{})
		Expect(func() {
			sp.Prepare(context.Background(), nil, &types.Header{
				GasLimit: 1000000,
			})
		}).ToNot(Panic())
	})
})
