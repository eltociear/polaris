package models

import (
	"fmt"

	"github.com/berachain/stargazer/wasp"
)

var _ wasp.BaseModel = (*EthAccount)(nil)

type EthAccount struct {
	wasp.BasePersistenceModal `gorm:"type:int;auto_increment;not null;"`
	Address                   []byte          `gorm:"type:bytea;not null;column:address" json:"address"`
	Alias                     string          `gorm:"type:bytea;column:alias" json:"alias"`
	Balance                   string          `gorm:"type:string;column:balance" json:"balance"`
	Erc721Balance             []Erc721Balance `gorm:"foreignkey:owner;references:address" json:"erc721Balances"`
	Erc20Balance              []Erc20Balance  `gorm:"foreignkey:owner;references:address" json:"erc20Balances"`
	IsContract                bool            `gorm:"type:boolean;not null;default:false" json:"isContract"`
	Contract                  Contract        `gorm:"foreignkey:contract_address;references:address" json:"contract"`
}

func (m *EthAccount) GetId() int64 {
	return m.ID
}

func (m *EthAccount) GetTable() string {
	return "acct"
}

func (m *EthAccount) GetRedisDb() int64 {
	return 1
}

func (m *EthAccount) GetRedisKey() string {
	key := fmt.Sprintf("%s:%s", m.GetTable(), m.Address)
	return key
}
