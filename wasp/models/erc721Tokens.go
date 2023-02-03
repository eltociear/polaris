package models

import (
	"fmt"

	"github.com/berachain/stargazer/wasp"
)

var _ wasp.BaseModel = (*Erc721Tokens)(nil)

type Erc721Tokens struct {
	wasp.BasePersistenceModal `gorm:"type:int;auto_increment;not null;"`
	BalanceId                 int64  `gorm:"type:int;not null;column:balance_id" json:"balance_id"`
	TokenId                   int64  `gorm:"type:int;not null;column:token_id" json:"token_id"`
	Data                      []byte `gorm:"type:bytea;not null;column:data" json:"data"`
}

func (m *Erc721Tokens) GetId() int64 {
	return m.ID
}

func (m *Erc721Tokens) GetTable() string {
	return "acct"
}

func (m *Erc721Tokens) GetRedisDb() int64 {
	return 1
}

func (m *Erc721Tokens) GetRedisKey() string {
	key := fmt.Sprintf("%s:%s", m.GetTable(), m.BalanceId)
	return key
}
