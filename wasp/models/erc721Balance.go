package models

import (
	"fmt"

	"github.com/berachain/stargazer/wasp"
)

var _ wasp.BaseModel = (*Erc721Balance)(nil)

type Erc721Balance struct {
	wasp.BasePersistenceModal `gorm:"type:int;auto_increment;not null;"`
	Address                   []byte         `gorm:"type:bytea;not null;column:owner" json:"owner"`
	ContractAddress           []byte         `gorm:"type:bytea;not null;column:contract_address" json:"contractAddress"`
	Id                        []Erc721Tokens `gorm:"foreignkey:balance_id;references:id" json:"tokenIds"`
}

func (m *Erc721Balance) GetId() int64 {
	return m.ID
}

func (m *Erc721Balance) GetTable() string {
	return "acct"
}

func (m *Erc721Balance) GetRedisDb() int64 {
	return 1
}

func (m *Erc721Balance) GetRedisKey() string {
	key := fmt.Sprintf("%s:%s", m.GetTable(), m.Address)
	return key
}
