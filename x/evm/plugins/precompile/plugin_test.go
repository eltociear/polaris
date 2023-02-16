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
	"context"
	"math/big"

	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/stargazer/eth/common"
	"github.com/berachain/stargazer/eth/core/vm"
	"github.com/berachain/stargazer/lib/utils"
	"github.com/berachain/stargazer/testutil"
	"github.com/berachain/stargazer/x/evm/plugins/state/events"
	"github.com/berachain/stargazer/x/evm/plugins/state/events/mock"

	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("plugin", func() {
	var p *plugin
	var sdb *mockSDB
	var ctx sdk.Context

	BeforeEach(func() {
		ctx = testutil.NewContext()
		ctx = ctx.WithEventManager(
			events.NewManagerFrom(ctx.EventManager(), mock.NewPrecompileLogFactory()),
		)
		p = utils.MustGetAs[*plugin](NewPluginFrom(ctx))
		sdb = &mockSDB{}
	})

	It("should use correctly consume gas", func() {
		_, remainingGas, err := p.Run(sdb, &mockStateless{}, []byte{}, addr, new(big.Int), 30, false)
		Expect(err).To(BeNil())
		Expect(remainingGas).To(Equal(uint64(10)))
	})

	It("should error on insufficient gas", func() {
		_, _, err := p.Run(sdb, &mockStateless{}, []byte{}, addr, new(big.Int), 5, true)
		Expect(err.Error()).To(Equal("out of gas"))
	})

	It("should plug in custom gas configs", func() {
		Expect(p.KVGasConfig().DeleteCost).To(Equal(uint64(1000)))
		Expect(p.TransientKVGasConfig().DeleteCost).To(Equal(uint64(100)))

		p.Context = p.WithKVGasConfig(storetypes.GasConfig{})
		Expect(p.KVGasConfig().DeleteCost).To(Equal(uint64(0)))
		p.Context = p.WithTransientKVGasConfig(storetypes.GasConfig{})
		Expect(p.TransientKVGasConfig().DeleteCost).To(Equal(uint64(0)))
	})
})

// MOCKS BELOW.

type mockSDB struct {
	vm.GethStateDB
}

type mockStateless struct{}

var addr = common.BytesToAddress([]byte{1})

func (ms *mockStateless) RegistryKey() common.Address {
	return addr
}

func (ms *mockStateless) Run(
	ctx context.Context, input []byte,
	caller common.Address, value *big.Int, readonly bool,
) ([]byte, error) {
	sdk.UnwrapSDKContext(ctx).GasMeter().ConsumeGas(10, "")
	return nil, nil
}

func (ms *mockStateless) RequiredGas(input []byte) uint64 {
	return 10
}

func (ms *mockStateless) WithStateDB(vm.GethStateDB) vm.PrecompileContainer {
	return ms
}
