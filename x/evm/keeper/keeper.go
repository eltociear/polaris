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

package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/berachain/stargazer/eth/core"
	"github.com/berachain/stargazer/eth/params"
	"github.com/berachain/stargazer/x/evm/constants"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/tendermint/tendermint/libs/log"
)

// Compile-time interface assertions.
var _ core.StargazerHostChain = (*Keeper)(nil)

type Keeper struct {
	// The (unexposed) key used to access the store from the Context.
	storeKey storetypes.StoreKey

	stateProcessor *core.StateProcessor
	// It is used to retrieve infofrmation about the current / past
	// blocks and associated validator information.
	stakingKeeper StakingKeeper
}

// NewKeeper creates new instances of the stargazer Keeper.
func NewKeeper(
	storeKey storetypes.StoreKey,
	stakingKeeper StakingKeeper,
) *Keeper {
	k := &Keeper{
		storeKey:      storeKey,
		stakingKeeper: stakingKeeper,
	}
	config := params.EthChainConfig{}
	k.stateProcessor = core.NewStateProcessor(&config, k)
	return k
}

// `Logger` returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", constants.EvmNamespace)
}

// `CumulativeGasUsed` returns the cumulative gas used in the current block.
func (k *Keeper) CumulativeGasUsed(ctx context.Context, gasUsed uint64) uint64 {
	cumulativeGasUsed := gasUsed
	sCtx := sdk.UnwrapSDKContext(ctx)
	if sCtx.BlockGasMeter() != nil {
		limit := sCtx.BlockGasMeter().Limit()
		cumulativeGasUsed += sCtx.BlockGasMeter().GasConsumed()
		if cumulativeGasUsed > limit {
			cumulativeGasUsed = limit
		}
	}

	return cumulativeGasUsed
}
