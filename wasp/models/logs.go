package models

import (
	"fmt"
	"strconv"

	"github.com/berachain/stargazer/core/types"
	"github.com/berachain/stargazer/wasp"
)

var _ wasp.BaseModel = (*EthLog)(nil)

type EthLog struct {
	wasp.BasePersistenceModal `gorm:"type:int;auto_increment;not null;"`
	Address                   []byte   `gorm:"type:bytea;not null;" json:"address"`
	Topics                    [][]byte `gorm:"type:bytea;" json:"topics"`
	Data                      []byte   `gorm:"type:bytea;not null;" json:"data"`
	BlockNumber               string   `gorm:"type:string;not null;" json:"blockNumber"`
	TxHash                    []byte   `gorm:"type:bytea;not null;column:tx_hash;" json:"transactionHash"`
	TxIndex                   uint     `gorm:"type:int;not null;" json:"transactionIndex"`
	BlockHash                 []byte   `gorm:"type:bytea;not null;" json:"blockHash"`
	Index                     uint     `gorm:"type:int;not null;" json:"logIndex"`
	Removed                   bool     `gorm:"type:bool;not null;" json:"removed"`
}

func GethToEthLogModel(l *types.Log) *EthLog {
	return &EthLog{
		Address:     l.Address.Bytes(),
		Topics:      nil,
		Data:        l.Data,
		BlockNumber: strconv.FormatUint(l.BlockNumber, 10),
		TxHash:      l.TxHash.Bytes(),
		TxIndex:     l.TxIndex,
		BlockHash:   l.BlockHash.Bytes(),
		Index:       l.Index,
		Removed:     l.Removed,
	}
}
func (m *EthLog) GetId() int64 {
	return m.ID
}

func (m *EthLog) GetTable() string {
	return "block_models"
}

func (m *EthLog) GetRedisDb() int64 {
	return 1
}

func (m *EthLog) GetRedisKey() string {
	key := fmt.Sprintf("%s:%s", m.GetTable(), m.BlockNumber)
	return key
}
