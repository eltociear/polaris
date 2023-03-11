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

package core

import (
	"pkg.berachain.dev/polaris/eth/common"
	"pkg.berachain.dev/polaris/eth/core/precompile"
	"pkg.berachain.dev/polaris/eth/core/state"
	"pkg.berachain.dev/polaris/eth/core/types"
	"pkg.berachain.dev/polaris/eth/params"
	libtypes "pkg.berachain.dev/polaris/lib/types"
)

// `PolarisHostChain` defines the plugins that the chain running Polaris EVM should implement.
type PolarisHostChain interface {
	// `GetBlockPlugin` returns the `BlockPlugin` of the Polaris host chain.
	GetBlockPlugin() BlockPlugin
	// `GetConfigurationPlugin` returns the `ConfigurationPlugin` of the Polaris host chain.
	GetConfigurationPlugin() ConfigurationPlugin
	// `GetGasPlugin` returns the `GasPlugin` of the Polaris host chain.
	GetGasPlugin() GasPlugin
	// `GetHistoricalPlugin` returns the OPTIONAL `HistoricalPlugin` of the Polaris host chain.
	GetHistoricalPlugin() HistoricalPlugin
	// `GetPrecompilePlugin` returns the OPTIONAL `PrecompilePlugin` of the Polaris host chain.
	GetPrecompilePlugin() PrecompilePlugin
	// `GetStatePlugin` returns the `StatePlugin` of the Polaris host chain.
	GetStatePlugin() StatePlugin
	// `GetTxPoolPlugin` returns the `TxPoolPlugin` of the Polaris host chain.
	GetTxPoolPlugin() TxPoolPlugin
}

// =============================================================================
// Mandatory Plugins
// =============================================================================

// The following plugins should be implemented by the chain running Polaris EVM and exposed via
// the `PolarisHostChain` interface. All plugins should be resettable with a given context.
type (
	// `BlockPlugin` defines the methods that the chain running Polaris EVM should implement to
	// support getting and setting block headers.
	BlockPlugin interface {
		// `BlockPlugin` implements `libtypes.Preparable`. Calling `Prepare` should reset the
		// `BlockPlugin` to a default state.
		libtypes.Preparable
		// `NewHeaderWithBlockNumber` returns a new block header with the given block number.
		NewHeaderWithBlockNumber(int64) *types.Header
		// `GetHeaderByNumber` returns the block header at the given block number.
		GetHeaderByNumber(int64) (*types.Header, error)
		// `SetHeaderByNumber` sets the block header at the given block number.
		SetHeaderByNumber(int64, *types.Header) error
		// `BaseFee` returns the base fee of the current block.
		BaseFee() uint64
	}

	// `ConfigurationPlugin` defines the methods that the chain running Polaris EVM should
	// implement in order to configuration the parameters of the Polaris EVM.
	ConfigurationPlugin interface {
		// `ConfigurationPlugin` implements `libtypes.Preparable`. Calling `Prepare` should reset
		// the `ConfigurationPlugin` to a default state.
		libtypes.Preparable
		// `ChainConfig` returns the current chain configuration of the Polaris EVM.
		ChainConfig() *params.ChainConfig
		// `ExtraEips` returns the extra EIPs that the Polaris EVM supports.
		ExtraEips() []int
		// `The fee collector is utilized on chains that have a fee collector account. This was added
		// specifically to support Cosmos-SDK chains, where we want the coinbase in the block header
		// to be the operator address of the proposer, but we want the coinbase in the BlockContext
		// to be the FeeCollectorAccount.
		FeeCollector() *common.Address
	}

	// `GasPlugin` is an interface that allows the Polaris EVM to consume gas on the host chain.
	GasPlugin interface {
		// `GasPlugin` implements `libtypes.Preparable`. Calling `Prepare` should reset the
		// `GasPlugin` to a default state.
		libtypes.Preparable
		// `GasPlugin` implements `libtypes.Resettable`. Calling `Reset` should reset the
		// `GasPlugin` to a default state
		libtypes.Resettable
		// `ConsumeGas` consumes the supplied amount of gas. It should not panic due to a
		// `GasOverflow` and should return `core.ErrOutOfGas` if the amount of gas remaining is
		// less than the amount requested. If the requested amount is greater than the amount of
		// gas remaining in the block, it should return core.ErrBlockOutOfGas.
		ConsumeGas(uint64) error
		// `GasRemaining` returns the amount of gas remaining for the current transaction.
		GasRemaining() uint64
		// `GasConsumed` returns the amount of gas used by the current transaction.
		GasConsumed() uint64
		// `BlockGasConsumed` returns the amount of gas used during the current block. The value
		// returned should NOT include any gas consumed during this transaction.
		// It should not panic.
		BlockGasConsumed() uint64
		// `BlockGasLimit` returns the gas limit of the current block. It should not panic.
		BlockGasLimit() uint64
	}

	// `StatePlugin` defines the methods that the chain running Polaris EVM should implement.
	StatePlugin interface {
		state.Plugin
		// `GetStateByNumber` returns the state at the given block height.
		GetStateByNumber(int64) (StatePlugin, error)
	}

	// `TxPoolPlugin` defines the methods that the chain running Polaris EVM should implement to
	// support the transaction pool.
	TxPoolPlugin interface {
		// `SendTx` submits the tx to the transaction pool.
		SendTx(tx *types.Transaction) error
		// `GetAllTransactions` returns all transactions in the transaction pool.
		GetAllTransactions() (types.Transactions, error)
		// `GetTransaction` returns the transaction from the pool with the given hash.
		GetTransaction(common.Hash) *types.Transaction
		// `GetNonce` returns the nonce of the given address in the transaction pool.
		GetNonce(common.Address) (uint64, error)
	}
)

// =============================================================================
// Optional Plugins
// =============================================================================

// `The following plugins are OPTIONAL to be implemented by the chain running Polaris EVM.
type (
	// `HistoricalPlugin` defines the methods that the chain running Polaris EVM should implement
	// in order to support storing historical blocks, receipts, and transactions. This plugin will
	// be used by the RPC backend to support certain methods on the Ethereum JSON RPC spec.
	// Implementing this plugin is optional.
	HistoricalPlugin interface {
		// `HistoricalPlugin` implements `libtypes.Preparable`.
		libtypes.Preparable
		// `GetBlockByNumber` returns the block at the given block number.
		GetBlockByNumber(int64) (*types.Block, error)
		// `GetBlockByHash` returns the block at the given block hash.
		GetBlockByHash(common.Hash) (*types.Block, error)
		// `GetTransactionByHash` returns the transaction lookup entry at the given transaction
		// hash.
		GetTransactionByHash(common.Hash) (*types.TxLookupEntry, error)
		// `GetReceiptByHash` returns the receipts at the given block hash.
		GetReceiptsByHash(common.Hash) (types.Receipts, error)
		// `StoreBlock` stores the given block.
		StoreBlock(*types.Block) error
		// `StoreReceipts` stores the receipts for the given block hash.
		StoreReceipts(common.Hash, types.Receipts) error
		// `StoreTransactions` stores the transactions for the given block hash.
		StoreTransactions(int64, common.Hash, types.Transactions) error
	}

	// `PrecompilePlugin` defines the methods that the chain running Polaris EVM should implement
	// in order to support running their own stateful precompiled contracts. Implementing this
	// plugin is optional.
	PrecompilePlugin = precompile.Plugin
)
