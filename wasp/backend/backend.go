package backend

import (
	"context"

	"github.com/berachain/stargazer/wasp/backend/postgres"
	"github.com/berachain/stargazer/wasp/backend/redis"

	"github.com/berachain/stargazer/wasp/types"

	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
)

type BackendI interface {
	GetValidators() ([]stakingTypes.Validator, error)
}

var _ BackendI = (*Backend)(nil)

type Backend struct {
	ctx            context.Context
	queryClient    *types.QueryClient
	redisClient    *redis.Redis
	postgresClient *postgres.Postgres
}

func NewBackend(
	cc *grpc.ClientConn,
	allowUnprotectedTxs bool,
) *Backend {
	postgresClient := postgres.NewFromEnv()
	redisClient := redis.NewFromEnv()
	return &Backend{
		ctx:            context.Background(),
		queryClient:    types.NewQueryClient(cc),
		postgresClient: &postgresClient,
		redisClient:    &redisClient,
	}
}
