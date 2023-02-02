package queryClient

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/ethclient"
)

type QueryClient struct {
	client *ethclient.Client
}

func NewQueryClient(client *ethclient.Client) *QueryClient {
	return &QueryClient{
		client: client,
	}
}

func (c *QueryClient) GetTransactionReceiptByHash(hash common.Hash) (*types.Receipt, error) {
	return c.client.TransactionReceipt(context.Background(), hash)
}

func (c *QueryClient) NetworkID(ctx context.Context) (*big.Int, error) {
	return c.client.NetworkID(ctx)
}
