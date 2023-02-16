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

package state_test

import (
	"math/big"
	"testing"

	"github.com/berachain/stargazer/eth/common"
	"github.com/berachain/stargazer/eth/core"
	"github.com/berachain/stargazer/testutil"
	"github.com/berachain/stargazer/x/evm/plugins/state"
)

var (
	numCalls           = 10000 // number of times snapshot is called
	numStoreOpsPerCall = 10    // number of read/write ops on stores during each call
	numReverts         = 2     // number of times an eth call is reverted in one tx
)

func GetNewStatePlugin() core.StatePlugin {
	ctx, ak, bk, _ := testutil.SetupMinimalKeepers()
	return state.NewPlugin(ctx, ak, bk, testutil.EvmKey, "abera", nil)
}

func BenchmarkArbitraryStateTransition(b *testing.B) {
	sp := GetNewStatePlugin()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var snapshots []int
		for c := 0; c < numCalls; c++ {
			sp.SetNonce(testutil.Bob, uint64(c+19)) // accStore set
			sp.SetState(                            // ethStore set
				testutil.Alice,
				common.BytesToHash([]byte{byte(c + 11)}),
				common.BytesToHash([]byte{byte(c + 22)}),
			)

			snapshots = append(snapshots, sp.Snapshot())
			for s := 0; s < numStoreOpsPerCall; s++ {
				sp.GetBalance(testutil.Alice)               // bankStore read
				sp.AddBalance(testutil.Bob, big.NewInt(10)) // bankStore write
				sp.GetCode(testutil.Alice)                  // ethStore read
			}
			if c > 0 && numReverts > 0 && c%(numCalls/numReverts) == 0 {
				sp.RevertToSnapshot(snapshots[len(snapshots)/2])
			}

			sp.DeleteSuicides([]common.Address{testutil.Alice}) // will invoke a delete
		}

		// commit only once
		sp.Finalize()
	}
}
