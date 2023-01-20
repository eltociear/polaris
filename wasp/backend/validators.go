package backend

import (
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (b *Backend) GetValidators() ([]stakingtypes.Validator, error) {
	queryClient := b.queryClient
	req := &stakingtypes.QueryValidatorsRequest{
		Status: "BOND_STATUS_BONDED",
	}
	res, err := queryClient.Validators(b.ctx, req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var v = res.GetValidators()
	return v, nil
}
