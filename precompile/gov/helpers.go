package gov

import (
	"context"
	"math/big"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"pkg.berachain.dev/stargazer/precompile/contracts/solidity/generated"
)

func (c *Contract) submitProposalHelper(
	ctx context.Context,
	messages []*codectypes.Any,
	initialDeposit []generated.IGovernanceModuleCoin,
	proposer sdk.AccAddress,
	metadata, title, summary string,
	expedited bool,
) ([]any, error) {
	var coins []sdk.Coin

	// Convert the initial deposit to sdk.Coin.
	for _, coin := range initialDeposit {
		coins = append(coins, sdk.NewCoin(coin.Denom, sdk.NewIntFromBigInt(coin.Amount)))
	}

	res, err := c.msgServer.SubmitProposal(ctx, &v1.MsgSubmitProposal{
		Messages:       messages,
		InitialDeposit: coins,
		Proposer:       proposer.String(),
		Metadata:       metadata,
		Title:          title,
		Summary:        summary,
		Expedited:      expedited,
	})
	if err != nil {
		return nil, err
	}

	return []any{big.NewInt(int64(res.ProposalId))}, nil
}

func (c *Contract) cancelProposalHelper(
	ctx context.Context,
	proposer sdk.AccAddress,
	proposalID *big.Int,
) ([]any, error) {
	res, err := c.msgServer.CancelProposal(ctx, &v1.MsgCancelProposal{
		ProposalId: proposalID.Uint64(),
		Proposer:   proposer.String(),
	})
	if err != nil {
		return nil, err
	}

	return []any{big.NewInt(int64(res.ProposalId))}, nil
}
