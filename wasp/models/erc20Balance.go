package models

import (
	"fmt"

	"github.com/berachain/stargazer/wasp"
)

var _ wasp.BaseModel = (*Erc20Balance)(nil)

type Erc20Balance struct {
	wasp.BasePersistenceModal `gorm:"type:int;auto_increment;not null;"`
	Address                   []byte `gorm:"type:bytea;not null;column:owner" json:"owner"`
	ContractAddress           []byte `gorm:"type:bytea;not null;column:contract_address" json:"contractAddress"`
	Amount                    string `gorm:"type:string;not null;column:amount" json:"amount"`
}

func (m *Erc20Balance) GetId() int64 {
	return m.ID
}

func (m *Erc20Balance) GetTable() string {
	return "acct"
}

func (m *Erc20Balance) GetRedisDb() int64 {
	return 1
}

func (m *Erc20Balance) GetRedisKey() string {
	key := fmt.Sprintf("%s:%s", m.GetTable(), m.Address)
	return key
}
