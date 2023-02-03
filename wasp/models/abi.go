package models

import (
	"fmt"

	"github.com/berachain/stargazer/wasp"
)

var _ wasp.BaseModel = (*Abi)(nil)

type Abi struct {
	wasp.BasePersistenceModal `gorm:"type:int;auto_increment;not null;"`
	Abi                       string `gorm:"type:string;not null;column:abi;unique;" json:"abi"`
	Tag                       string `gorm:"type:string;not null;column:tag" json:"tag"`
}

func (m *Abi) GetId() int64 {
	return m.ID
}

func (m *Abi) GetTable() string {
	return "acct"
}

func (m *Abi) GetRedisDb() int64 {
	return 1
}

func (m *Abi) GetRedisKey() string {
	key := fmt.Sprintf("%s:%s", m.GetTable(), m.Tag)
	return key
}
