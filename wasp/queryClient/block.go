package queryClient

import (
	"context"

	"github.com/berachain/stargazer/wasp"
	"github.com/berachain/stargazer/wasp/models"
)

func (c *QueryClient) GetLatestBlock(ctx context.Context) (*models.EthBlockModel, error) {
	// SHOULD BE SENT TO THE NODE
	blockClient := c.db.EthBlockModel
	blockModel, err := blockClient.WithContext(ctx).First()
	return blockModel, err
}

func (c *QueryClient) GetBlock(ctx context.Context, height string) (*models.EthBlockModel, error) {
	blockClient := c.db.EthBlockModel
	blockModel, err := blockClient.WithContext(ctx).Preload(blockClient.Txs).Where(blockClient.BlockNumber.Eq(height)).First()
	return blockModel, err
}
func (c *QueryClient) GetBlocks(ctx context.Context, pagination wasp.Pagination) ([]*models.EthBlockModel, error) {
	blockClient := c.db.EthBlockModel
	offset := (pagination.Page - 1) * pagination.Limit
	blockModel, err := blockClient.WithContext(ctx).Preload(blockClient.Txs).Offset(offset).Limit(pagination.Limit).Find()
	return blockModel, err
}
