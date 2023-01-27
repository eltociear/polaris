package wasp

import (
	"github.com/berachain/stargazer/wasp/proto"
)

type Database interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) (interface{}, error)
}

type DbClient interface {
	proto.BlockServiceServer
}

type BaseModel interface {
	GetTable() string
	GetId() int64
	GetRedisDb() int64
	GetRedisKey() string
}

type BasePersistenceModal struct {
	ID int64 `gorm:"type:int;primary_key"`
}
