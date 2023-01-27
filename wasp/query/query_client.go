package query

import (
	"google.golang.org/grpc"

	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type QueryClient struct {
	stakingClient stakingTypes.QueryClient
}

func NewQueryClient(cc *grpc.ClientConn) *QueryClient {
	return &QueryClient{
		stakingClient: stakingTypes.NewQueryClient(cc),
	}
}
