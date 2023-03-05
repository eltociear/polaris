package gov

import (
	"context"
	"math/big"

	governancekeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	governancetypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"pkg.berachain.dev/stargazer/eth/accounts/abi"
	"pkg.berachain.dev/stargazer/eth/common"
	"pkg.berachain.dev/stargazer/eth/core/precompile"
	"pkg.berachain.dev/stargazer/lib/utils"
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
		msgServer:   governancekeeper.NewMsgServerImpl(*gk),
		querier:     *gk,
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

// `PrecompileMethods` implements the `precompile.StatefulImpl` interface.
func (c *Contract) PrecompileMethods() precompile.Methods {
	return precompile.Methods{
		&precompile.Method{
			AbiSig: "submitProposal(bytes,[]tuple,string,string,string,string,bool)",
		},
	}
}

// `SubmitProposal` is the method for the `submitProposal` method of the governance precompile contract.
func (c *Contract) SubmitProposalStringAddr(
	ctx context.Context,
	caller common.Address,
	value *big.Int,
	readonly bool,
	args ...any,
) ([]any, error) {
	message, ok := utils.GetAs[[]*codectypes.Any](args[0]) // TODO: check if this is the correct type.
	if !ok {
		return nil, ErrInvalidBytes
	}
	initialDeposit, ok := utils.GetAs[[]generated.IGovernanceModuleCoin](args[1])
	if !ok {
		return nil, ErrInvalidCoins
	}
	proposer, ok := utils.GetAs[string](args[2])
	if !ok {
		return nil, ErrInvalidString
	}
	metadata, ok := utils.GetAs[string](args[3])
	if !ok {
		return nil, ErrInvalidString
	}
	title, ok := utils.GetAs[string](args[4])
	if !ok {
		return nil, ErrInvalidString
	}
	summary, ok := utils.GetAs[string](args[5])
	if !ok {
		return nil, ErrInvalidString
	}
	expedited, ok := utils.GetAs[bool](args[6])
	if !ok {
		return nil, ErrInvalidBool
	}
	proposerAddr, err := sdk.AccAddressFromBech32(proposer)
	if err != nil {
		return nil, err
	}

	return c.submitProposalHelper(ctx, message, initialDeposit, proposerAddr, metadata, title, summary, expedited)
}
