package repository

import (
	"context"
	"encoding/json"

	"github.com/berachain/stargazer/wasp/database"
	"github.com/berachain/stargazer/wasp/models"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockRepo struct {
	db *database.Database
}

func NewBlockRepo(db *database.Database) *BlockRepo {
	return &BlockRepo{
		db: db,
	}
}

func (repo *BlockRepo) CreateBlock(ctx context.Context, msg *types.Block, signerType types.Signer) int {
	blockModel := models.GethToBlockModel(msg, &signerType)
	data, err := json.Marshal(blockModel)
	if err != nil {
		panic(err)
	}

	req := &database.SetRequest{
		RedisDb: blockModel.GetRedisDb(),
		Key:     blockModel.GetRedisKey(),
		Value:   data,
	}

	err = repo.db.Set(req, func() error {
		res := repo.db.Gorm.Create(&blockModel)
		return res.Error
	})

	if err != nil {
		return 1
	}

	return 0
}

// func (repo *BlockRepo) GetBlock(ctx context.Context, msg *stargazerproto.ReadBlockRequest) *stargazerproto.ReadBlockResponse {
// 	blockModel := models.FromGrpcBlockMessage(&stargazerproto.Block{
// 		Number: msg.Number,
// 	})
// 	req := &database.GetRequest{
// 		RedisDb: blockModel.GetRedisDb(),
// 		Key:     blockModel.GetRedisKey(),
// 	}

// 	data, err := repo.db.Get(req, func() ([]byte, error) {
// 		response := &models.BlockModel{}
// 		err := repo.db.Gorm.Where("number = ?", blockModel.Number).Last(&response).Error
// 		if err != nil {
// 			return nil, err
// 		}
// 		block := response.ToGrpcBlockModel()
// 		data, err := proto.Marshal(block)

// 		return data, err
// 	})

// 	if err != nil {
// 		fmt.Printf("Unable to GET data")
// 		return nil
// 	}
// 	block := &stargazerproto.Block{}
// 	err = proto.Unmarshal(data, block)
// 	if err != nil {
// 		fmt.Print("Unable to unmarshal data", err)
// 		return nil

// 	}
// 	return &stargazerproto.ReadBlockResponse{
// 		Block: block,
// 	}
// }
