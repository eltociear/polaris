package queryClient

import (
	"context"
	"fmt"
	"math/big"

	"github.com/berachain/stargazer/wasp/abi"
	"github.com/berachain/stargazer/wasp/query"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/ethclient"
)

type QueryClient struct {
	client *ethclient.Client
	db     *query.Query
}

func NewQueryClient(client *ethclient.Client, gormDb *gorm.DB) *QueryClient {
	return &QueryClient{
		client: client,
		db:     query.Use(gormDb),
	}
}

func (c *QueryClient) GetTransactionReceiptByHash(hash common.Hash) (*types.Receipt, error) {
	return c.client.TransactionReceipt(context.Background(), hash)
}

func (c *QueryClient) NetworkID(ctx context.Context) (*big.Int, error) {
	return c.client.NetworkID(ctx)
}

func (c *QueryClient) CodeAt(ctx context.Context, contractAddress []byte) ([]byte, error) {
	code, err := c.client.CodeAt(context.Background(), common.BytesToAddress(contractAddress), nil)
	if err != nil {
		fmt.Println("Error getting contract code:", err)
		return nil, err
	}
	return code, nil
}

func (c *QueryClient) GetEthClient() *ethclient.Client {
	return c.client
}

func (c *QueryClient) GetEthBalance(ctx context.Context, address []byte, blockNumber uint64) (*big.Int, error) {
	return c.client.BalanceAt(ctx, common.BytesToAddress(address), nil)

}

func (c *QueryClient) GetErc20Balance(ctx context.Context, contractAddress []byte, targetAddress []byte) (*big.Int, error) {
	contract := common.BytesToAddress(contractAddress)
	target := common.BytesToAddress(targetAddress)
	erc20, err := abi.NewERC20Caller(contract, c.client)
	if err != nil {
		return nil, err
	}

	balance, err := erc20.BalanceOf(nil, target)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
