package types

import (
	"google.golang.org/grpc"

	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type QueryClient struct {
	stakingTypes.QueryClient
}

func NewQueryClient(cc *grpc.ClientConn) *QueryClient {
	return &QueryClient{
		QueryClient: stakingTypes.NewQueryClient(cc),
	}
}
