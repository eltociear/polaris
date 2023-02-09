package repository

import (
	"context"
	"encoding/json"

	"github.com/berachain/stargazer/wasp/database"
	"github.com/berachain/stargazer/wasp/models"
	"github.com/berachain/stargazer/wasp/queryClient"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockRepo struct {
	db *database.Database
	qc *queryClient.QueryClient
}

func NewBlockRepo(db *database.Database, qc *queryClient.QueryClient) *BlockRepo {
	return &BlockRepo{
		db: db,
		qc: qc,
	}
}

func (repo *BlockRepo) CreateBlock(ctx context.Context, block *models.EthBlockModel) int {
	data, err := json.Marshal(block)
	if err != nil {
		panic(err)
	}

	req := &database.SetRequest{
		RedisDb: block.GetRedisDb(),
		Key:     block.GetRedisKey(),
		Value:   data,
	}

	err = repo.db.Set(req, func() error {
		res := repo.db.Gorm.Create(&block)
		return res.Error
	})

	if err != nil {
		return 1
	}

	return 0
}

func (repo *BlockRepo) BuildBlock(block *types.Block, txns []models.TransactionModel) *models.EthBlockModel {
	blockModel := models.GethToBlockModel(block, txns)
	return blockModel
}
