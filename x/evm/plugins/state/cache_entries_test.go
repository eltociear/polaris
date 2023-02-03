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
	coretypes "github.com/berachain/stargazer/eth/core/types"
	"github.com/berachain/stargazer/x/evm/plugins/state/store/journal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AddLogChange", func() {
	var (
		ce  *AddLogChange
		sdb *StateDB
	)
	BeforeEach(func() {
		sdb = &StateDB{
			logs: []*coretypes.Log{},
		}
		ce = &AddLogChange{
			sdb: sdb,
		}
	})
	It("implements journal.CacheEntry", func() {
		var _ journal.CacheEntry = ce
		Expect(ce).To(BeAssignableToTypeOf(&AddLogChange{}))
	})
	It("Revert should remove the last log", func() {
		sdb.logs = append(sdb.logs, &coretypes.Log{})
		ce.Revert()
		Expect(len(sdb.logs)).To(Equal(0))
	})
	It("Clone should return a new AddLogChange with the same sdb", func() {
		cloned, ok := ce.Clone().(*AddLogChange)
		Expect(ok).To(BeTrue())
		Expect(cloned.sdb).To(Equal(sdb))
		Expect(cloned).ToNot(BeIdenticalTo(ce))
	})
})

var _ = Describe("RefundChange", func() {
	var (
		ce  *RefundChange
		sdb *StateDB
	)

	BeforeEach(func() {
		sdb = &StateDB{
			refund: 0,
		}
		ce = &RefundChange{
			sdb:  sdb,
			prev: 0,
		}
	})
	It("implements journal.CacheEntry", func() {
		var _ journal.CacheEntry = ce
		Expect(ce).To(BeAssignableToTypeOf(&RefundChange{}))
	})
	It("Revert should restore the previous refund value", func() {
		sdb.refund = 100
		ce.prev = 50
		ce.Revert()
		Expect(sdb.refund).To(Equal(uint64(50)))
	})
	It("Clone should return a new RefundChange with the same sdb and prev", func() {
		cloned, ok := ce.Clone().(*RefundChange)
		Expect(ok).To(BeTrue())
		Expect(cloned.sdb).To(Equal(sdb))
		Expect(cloned.prev).To(Equal(ce.prev))
		Expect(cloned).ToNot(BeIdenticalTo(ce))
	})
})
