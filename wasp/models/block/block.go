package models

import (
	"fmt"
	"strconv"

	"github.com/berachain/stargazer/wasp"
	proto "github.com/berachain/stargazer/wasp/proto"
)

var _ wasp.BaseModel = (*BlockModel)(nil)

type BlockModel struct {
	wasp.BasePersistenceModal `gorm:"type:int;auto_increment;not null;"`
	Number                    int64 `gorm:"type:bigint;not null;"`
}

func (m *BlockModel) ToGrpcBlockModel() *proto.Block {
	return &proto.Block{
		Number: m.Number,
	}
}

func FromGrpcBlockMessage(block *proto.Block) *BlockModel {
	return &BlockModel{
		Number: block.GetNumber(),
	}
}

func (m *BlockModel) GetId() int64 {
	return m.ID
}

func (m *BlockModel) GetTable() string {
	return "blocks"
}

func (m *BlockModel) GetRedisDb() int64 {
	return 1
}

func (m *BlockModel) GetRedisKey() string {
	key := fmt.Sprintf("%s:%s", m.GetTable(), strconv.FormatInt(m.Number, 10))
	return key
}
