package backend

import (
	"context"
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (b *Backend) GetValidators() ([]stakingtypes.Validator, error) {

	req := &stakingtypes.QueryValidatorsRequest{}
	res, err := b.queryClient.Validators(context.Background(), req)

	if err != nil {
		return nil, err
	}
	fmt.Println(res)
	var v = res.GetValidators()
	return v, nil
}
