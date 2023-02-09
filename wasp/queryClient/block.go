package queryClient

import (
	"context"

	"github.com/berachain/stargazer/wasp/models"
)

func (c *QueryClient) GetLatestBlock(ctx context.Context) (*models.EthBlockModel, error) {
	blockClient := c.db.EthBlockModel
	blockModel, err := blockClient.WithContext(ctx).First()
	return blockModel, err
}
