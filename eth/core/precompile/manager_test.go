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

package precompile

import (
	"math/big"

	"github.com/berachain/stargazer/eth/core/state"
	"github.com/berachain/stargazer/eth/core/vm"
	"github.com/berachain/stargazer/lib/common"
	"github.com/berachain/stargazer/lib/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("controller", func() {
	var c *manager
	var mr *mockRunner

	BeforeEach(func() {
		mr = &mockRunner{}
		c = utils.MustGetAs[*manager](NewManager(mr))
		err := c.Register(&mockStateless{})
		Expect(err).To(BeNil())
	})

	It("should find and run", func() {
		err := c.PrepareForStateTransition(&mockSdb{&state.StateDB{}})
		Expect(err).To(BeNil())

		pc := c.Get(addr)
		Expect(pc).ToNot(BeNil())

		_, _, err = c.Run(pc, []byte{}, addr, new(big.Int), 10, true)
		Expect(err).To(BeNil())
		Expect(mr.called).To(BeTrue())
		Expect(mr.calledWithStateDB).To(BeTrue())
	})

	It("should not find an unregistered", func() {
		found := c.Has(common.BytesToAddress([]byte{2}))
		Expect(found).To(BeFalse())
	})

	It("should error on incompatible statedb", func() {
		err := c.PrepareForStateTransition(badMockSdb{&state.StateDB{}})
		Expect(err.Error()).To(Equal("statedb is not compatible with Stargazer"))
	})
})

// MOCKS BELOW.

type mockSdb struct {
	vm.StargazerStateDB
}

type badMockSdb struct {
	vm.GethStateDB
}

type mockRunner struct {
	called            bool
	calledWithStateDB bool
}

func (mr *mockRunner) Run(
	pc vm.PrecompileContainer, statedb vm.StargazerStateDB, input []byte,
	caller common.Address, value *big.Int, suppliedGas uint64, readonly bool,
) ([]byte, uint64, error) {
	mr.called = true
	mr.calledWithStateDB = statedb != nil
	return nil, 0, nil
}
