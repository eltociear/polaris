package repository

import (
	"context"
	"fmt"

	"github.com/berachain/stargazer/wasp/database"
	models "github.com/berachain/stargazer/wasp/models/block"
	stargazerproto "github.com/berachain/stargazer/wasp/proto"
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
		res := repo.db.Gorm.Create(&blockModel)
		fmt.Print("\n")

		fmt.Print(blockModel.ID)
		return res.Error
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
		response := &models.BlockModel{}
		err := repo.db.Gorm.Where("number = ?", blockModel.Number).Last(&response).Error
		if err != nil {
			return nil, err
		}
		block := response.ToGrpcBlockModel()
		data, err := proto.Marshal(block)

		return data, err
	})

	if err != nil {
		fmt.Printf("Unable to GET data")
		return nil
	}
	block := &stargazerproto.Block{}
	err = proto.Unmarshal(data, block)
	if err != nil {
		fmt.Print("Unable to unmarshal data", err)
		return nil

	}
	return &stargazerproto.ReadBlockResponse{
		Block: block,
	}
}
