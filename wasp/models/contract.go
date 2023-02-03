package models

import (
	"fmt"

	"github.com/berachain/stargazer/wasp"
)

var _ wasp.BaseModel = (*Contract)(nil)

type Contract struct {
	wasp.BasePersistenceModal `gorm:"type:int;auto_increment;not null;"`
	Address                   []byte `gorm:"type:bytea;column:contract_address" json:"contract_address"`
	Creator                   []byte `gorm:"type:bytea;column:creator" json:"creator"`
	DeployTxnHash             []byte `gorm:"type:bytea;column:deploy_txn_hash" json:"txnHash"`
	AbiId                     string `gorm:"type:int;not null;column:abi_id" json:"abi_id"`
	Abi                       Abi    `gorm:"foreignkey:id;references:abi_id" json:"contract"`
}

func (m *Contract) GetId() int64 {
	return m.ID
}

func (m *Contract) GetTable() string {
	return "acct"
}

func (m *Contract) GetRedisDb() int64 {
	return 1
}

func (m *Contract) GetRedisKey() string {
	key := fmt.Sprintf("%s:%s", m.GetTable(), m.Address)
	return key
}
