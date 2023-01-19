package backend

import (
	"context"

	"github.com/berachain/stargazer/wasp/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
)

// BackendI implements the Cosmos and EVM backend.
type BackendI interface { //nolint: revive
	CosmosBackend
}

// CosmosBackend implements the functionality shared within cosmos namespaces
// as defined by Wallet Connect V2: https://docs.walletconnect.com/2.0/json-rpc/cosmos.
// Implemented by Backend.
type CosmosBackend interface { // TODO: define
	// GetAccounts()
	// SignDirect()
	// SignAmino()
	GetValidators() ([]stakingTypes.Validator, error)
}

// EVMBackend implements the functionality shared within ethereum namespaces
// as defined by EIP-1474: https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1474.md
// Implemented by Backend.
type EVMBackend interface {
	// Node specific queries
	// getValidators()
}

var _ BackendI = (*Backend)(nil)

// Backend implements the BackendI interface
type Backend struct {
	ctx                 context.Context
	queryClient         *types.QueryClient // gRPC query client
	allowUnprotectedTxs bool
}

// NewBackend creates a new Backend instance for cosmos and ethereum namespaces
func NewBackend(
	cc *grpc.ClientConn,
	allowUnprotectedTxs bool,
) *Backend {

	return &Backend{
		ctx:                 context.Background(),
		queryClient:         types.NewQueryClient(cc),
		allowUnprotectedTxs: allowUnprotectedTxs,
	}
}
