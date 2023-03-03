package gov

import (
	governancekeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	governancetypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"pkg.berachain.dev/stargazer/eth/accounts/abi"
	"pkg.berachain.dev/stargazer/eth/common"
	"pkg.berachain.dev/stargazer/eth/core/precompile"
	"pkg.berachain.dev/stargazer/precompile/contracts/solidity/generated"
	evmutils "pkg.berachain.dev/stargazer/x/evm/utils"
)

// `Contract` is the precompile contract for the governance module.
type Contract struct {
	contractAbi *abi.ABI

	msgServer v1.MsgServer
	querier   v1.QueryServer
}

// `NewContract` is the constructor for the governance precompile contract.
func NewContract(gk **governancekeeper.Keeper) precompile.StatefulImpl {
	var contractAbi abi.ABI
	if err := contractAbi.UnmarshalJSON([]byte(generated.GovernanceModuleMetaData.ABI)); err != nil {
		panic(err)
	}
	return &Contract{
		contractAbi: &contractAbi,
		msgServer:   governancekeeper.NewLegacyMsgServerImpl(*gk),
		querier:     governancekeeper.Keeper{},
	}
}

// `RegistryKey` implements the `precompile.StatefulImpl` interface.
func (c *Contract) RegistryKey() common.Address {
	return evmutils.AccAddressToEthAddress(authtypes.NewModuleAddress(governancetypes.ModuleName))
}

// `ABIMethods` implements the `precompile.StatefulImpl` interface.
func (c *Contract) ABIMethods() map[string]abi.Method {
	return c.contractAbi.Methods
}

// `ABIEvents` implements the `precompile.StatefulImpl` interface.
func (c *Contract) ABIEvents() map[string]abi.Event {
	return c.contractAbi.Events
}

// `CustomValueDecoders` implements the `precompile.StatefulImpl` interface.
func (c *Contract) CustomValueDecoders() precompile.ValueDecoders {
	return nil
}
