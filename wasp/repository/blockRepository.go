package repository

import (
	"context"

	"github.com/berachain/stargazer/wasp/database"
	models "github.com/berachain/stargazer/wasp/models/block"
	stargazerproto "github.com/berachain/stargazer/wasp/proto"
	"github.com/berachain/stargazer/wasp/utils"
	"google.golang.org/protobuf/proto"
)

type BlockRepo struct {
	db *database.Database
}

func NewBlockRepo(db *database.Database) *BlockRepo {
	return &BlockRepo{
		db: db,
	}
}

func (repo *BlockRepo) CreateBlock(ctx context.Context, msg *stargazerproto.CreateBlockRequest) *stargazerproto.CreateBlockResponse {
	blockModel := models.FromGrpcBlockMessage(msg.Block)
	data, err := proto.Marshal(msg.Block)
	if err != nil {
		panic(err)
	}

	err = repo.db.RedisClient.Set(blockModel.GetRedisKey(), data)
	if err != nil {
		return &stargazerproto.CreateBlockResponse{
			Code: 1,
		}
	}

	err = repo.db.Gorm.Table(blockModel.GetTable()).Create(&blockModel).Error
	if err != nil {
		return &stargazerproto.CreateBlockResponse{
			Code: 1,
		}
	}
	return &stargazerproto.CreateBlockResponse{
		Code: 0,
	}
}

func (repo *BlockRepo) GetBlock(ctx context.Context, msg *stargazerproto.ReadBlockRequest) *stargazerproto.ReadBlockResponse {
	blockModel := models.FromGrpcBlockMessage(&stargazerproto.Block{
		Number: msg.Number,
	})
	data, err := repo.db.RedisClient.Get(blockModel.GetRedisKey())
	block := &stargazerproto.Block{}
	if err != nil {
		result := repo.db.Gorm.Where("number = ?", blockModel.Number).First(&block)
		if result.Error != nil {
			panic(result.Error)
		}
		return &stargazerproto.ReadBlockResponse{
			Block: block,
		}
	}
	byteData, err := utils.GetBytes(data)
	if err != nil {
		panic(err)
	}
	err = proto.Unmarshal(byteData, block)
	if err != nil {
		panic(err)
	}
	return &stargazerproto.ReadBlockResponse{
		Block: block,
	}

}
