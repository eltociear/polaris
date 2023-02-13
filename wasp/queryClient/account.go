package queryClient

import (
	"context"

	"github.com/berachain/stargazer/wasp/models"
)

func (c *QueryClient) GetAccountErc20Balances(ctx context.Context, accountID uint64) ([]models.Erc20Balance, error) {
	return nil, nil
}
