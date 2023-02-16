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

package vm

import (
	coretypes "github.com/berachain/stargazer/eth/core/types"
	"github.com/berachain/stargazer/eth/params"
	libtypes "github.com/berachain/stargazer/lib/types"

	"math/big"

	"github.com/berachain/stargazer/eth/common"
)

type (
	// `StargazerEVM` defines an extension to the interface provided by Go-Ethereum to support
	// additional state transition functionalities.
	StargazerEVM interface {
		// `Create` creates a new contract using the given code and returns the contract address.
		Create(
			caller ContractRef, code []byte, gas uint64, value *big.Int,
		) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error)
		// `Call` executes the given contract with the given input data and returns the output.
		Call(
			caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int,
		) (ret []byte, leftOverGas uint64, err error)
		// `SetTxContext` sets the transaction context for the new transaction.
		SetTxContext(txCtx TxContext)
		// `StateDB` returns the `StateDB` attached to this `StargazerEVM`.
		StateDB() StargazerStateDB
		// `Config` returns the `Config` attached to this `StargazerEVM`.
		Config() Config
		// `Context` returns the `BlockContext` attached to this `StargazerEVM`.
		Context() BlockContext
		// `ChainConfig` returns the `ChainConfig` attached to this `StargazerEVM`.
		ChainConfig() *params.ChainConfig
	}

	// `StargazerStateDB` defines an extension to the interface provided by Go-Ethereum to support
	// additional state transition functionalities.
	StargazerStateDB interface {
		GethStateDB
		libtypes.Finalizeable
		// `Reset` resets the context for the new transaction.
		libtypes.Resettable
		// `TransferBalance` transfers the balance from one account to another
		TransferBalance(common.Address, common.Address, *big.Int)
		// `BuildLogsAndClear` builds the logs for the tx with the given metadata. NOTE: must be
		// called after `Finalize`.
		BuildLogsAndClear(
			txHash common.Hash, blockHash common.Hash, txIndex uint, logIndex uint,
		) []*coretypes.Log
	}

	// `RegistrablePrecompile` is a type for the base precompile implementation, which only needs to
	// provide an Ethereum address of where its contract is found.
	RegistrablePrecompile = libtypes.Registrable[common.Address]
)
