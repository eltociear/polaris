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

	req := &database.SetRequest{
		RedisDb: blockModel.GetRedisDb(),
		Key:     blockModel.GetRedisKey(),
		Value:   data,
	}

	err = repo.db.Set(req, func() error {
		err = repo.db.Gorm.Table(blockModel.GetTable()).Create(&blockModel).Error
		return err
	})

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
	req := &database.GetRequest{
		RedisDb: blockModel.GetRedisDb(),
		Key:     blockModel.GetRedisKey(),
	}
	data, err := repo.db.Get(req, func() ([]byte, error) {
		block := &stargazerproto.Block{}
		result := repo.db.Gorm.Where("number = ?", blockModel.Number).First(&block)
		if result.Error != nil {
			return nil, result.Error
		}
		protomsg := &stargazerproto.ReadBlockResponse{
			Block: block,
		}
		data, err := proto.Marshal(protomsg)
		return data, err
	})
	if err != nil {
		panic(err)
	}
	byteData, err := utils.GetBytes(data)
	if err != nil {
		panic(err)
	}
	block := &stargazerproto.Block{}
	err = proto.Unmarshal(byteData, block)
	if err != nil {
		panic(err)
	}
	return &stargazerproto.ReadBlockResponse{
		Block: block,
	}
}
