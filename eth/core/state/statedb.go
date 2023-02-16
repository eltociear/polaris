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

package state

import (
	"github.com/berachain/stargazer/eth/common"
	"github.com/berachain/stargazer/eth/core/state/journal"
	coretypes "github.com/berachain/stargazer/eth/core/types"
	"github.com/berachain/stargazer/eth/core/vm"
	"github.com/berachain/stargazer/eth/crypto"
	"github.com/berachain/stargazer/lib/snapshot"
	libtypes "github.com/berachain/stargazer/lib/types"
)

var (
	// `emptyCodeHash` is the Keccak256 Hash of empty code
	// 0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470.
	emptyCodeHash = crypto.Keccak256Hash(nil)
)

// `stateDB` is a struct that holds the plugins and controller to manage Ethereum state.
type stateDB struct {
	// Plugin is injected by the chain running the Stargazer EVM.
	Plugin

	// Journals built internally and required for the stateDB.
	LogsJournal
	RefundJournal

	// `ctrl` is used to manage snapshots and reverts across plugins and journals.
	ctrl libtypes.Controller[string, libtypes.Controllable[string]]

	// Dirty tracking of suicided accounts, we have to keep track of these manually, in order
	// for the code and state to still be accessible even after the account has been deleted.
	// We chose to keep track of them in a separate slice, rather than a map, because the
	// number of accounts that will be suicided in a single transaction is expected to be
	// very low.
	suicides []common.Address
}

// `NewStateDB` returns a `vm.StargazerStateDB` with the given `StatePlugin`.
func NewStateDB(sp Plugin) vm.StargazerStateDB {
	// Build the journals required for the stateDB
	lj := journal.NewLogs()
	rj := journal.NewRefund()

	// Build the controller and register the plugins and journals
	ctrl := snapshot.NewController[string, libtypes.Controllable[string]]()
	_ = ctrl.Register(lj)
	_ = ctrl.Register(rj)
	_ = ctrl.Register(sp)

	return &stateDB{
		Plugin:        sp,
		LogsJournal:   lj,
		RefundJournal: rj,
		ctrl:          ctrl,
		suicides:      make([]common.Address, 1), // very rare to suicide, so we alloc 1 slot.
	}
}

// =============================================================================
// Suicide
// =============================================================================

// Suicide implements the StargazerStateDB interface by marking the given address as suicided.
// This clears the account balance, but the code and state of the address remains available
// until after Commit is called.
func (sdb *stateDB) Suicide(addr common.Address) bool {
	// only smart contracts can commit suicide
	ch := sdb.GetCodeHash(addr)
	if (ch == common.Hash{}) || ch == emptyCodeHash {
		return false
	}

	// Reduce it's balance to 0.
	sdb.SubBalance(addr, sdb.GetBalance(addr))

	// Mark the underlying account for deletion in `Commit()`.
	sdb.suicides = append(sdb.suicides, addr)
	return true
}

// `HasSuicided` implements the `StargazerStateDB` interface by returning if the contract was suicided
// in current transaction.
func (sdb *stateDB) HasSuicided(addr common.Address) bool {
	for _, suicide := range sdb.suicides {
		if addr == suicide {
			return true
		}
	}
	return false
}

// `Empty` implements the `StargazerStateDB` interface by returning whether the state object
// is either non-existent or empty according to the EIP161 epecification
// (balance = nonce = code = 0)
// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-161.md
func (sdb *stateDB) Empty(addr common.Address) bool {
	ch := sdb.GetCodeHash(addr)
	return sdb.GetNonce(addr) == 0 &&
		(ch == emptyCodeHash || ch == common.Hash{}) &&
		sdb.GetBalance(addr).Sign() == 0
}

// =============================================================================
// Snapshot
// =============================================================================

// `Snapshot` implements `stateDB`.
func (sdb *stateDB) Snapshot() int {
	return sdb.ctrl.Snapshot()
}

// `RevertToSnapshot` implements `stateDB`.
func (sdb *stateDB) RevertToSnapshot(id int) {
	sdb.ctrl.RevertToSnapshot(id)
}

// =============================================================================
// Finalize
// =============================================================================

// `Finalize` deletes the suicided accounts, clears the suicides list, and finalizes all plugins.
func (sdb *stateDB) Finalize() {
	sdb.DeleteSuicides(sdb.suicides)
	sdb.suicides = make([]common.Address, 1)
	sdb.ctrl.Finalize()
}

// =============================================================================
// AccessList
// =============================================================================

func (sdb *stateDB) PrepareAccessList(
	sender common.Address,
	dst *common.Address,
	precompiles []common.Address,
	list coretypes.AccessList,
) {
	panic("not supported by Stargazer")
}

func (sdb *stateDB) AddAddressToAccessList(addr common.Address) {
	panic("not supported by Stargazer")
}

func (sdb *stateDB) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	panic("not supported by Stargazer")
}

func (sdb *stateDB) AddressInAccessList(addr common.Address) bool {
	return false
}

func (sdb *stateDB) SlotInAccessList(addr common.Address, slot common.Hash) (bool, bool) {
	return false, false
}

// =============================================================================
// PreImage
// =============================================================================

// AddPreimage implements the the `StateDB“ interface, but currently
// performs a no-op since the EnablePreimageRecording flag is disabled.
func (sdb *stateDB) AddPreimage(hash common.Hash, preimage []byte) {}
