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

package core

import (
	"github.com/berachain/stargazer/eth/api"
	"github.com/berachain/stargazer/eth/core/state"
	"github.com/berachain/stargazer/eth/core/vm"
)

// `blockchain` is the canonical, persistent object that operates the Stargazer EVM.
type blockchain struct {
	// `StateProcessor` is the canonical, persistent state processor that runs the EVM.
	*StateProcessor
	// `host` is the host chain that the Stargazer EVM is running on.
	host StargazerHostChain
}

// `NewChain` creates and returns a `api.Chain` with the given EVM chain configuration and host.
func NewChain(host StargazerHostChain) api.Chain {
	bc := &blockchain{
		host: host,
	}
	bc.StateProcessor = bc.buildStateProcessor(vm.Config{}, true)
	return bc
}

// `buildStateProcessor` builds and returns a `StateProcessor` with the given EVM configuration and
// commit flag.
func (bc *blockchain) buildStateProcessor(vmConfig vm.Config, commit bool) *StateProcessor {
	return NewStateProcessor(bc.host, state.NewStateDB(bc.host.GetStatePlugin()), vmConfig, commit)
}
