package models

import (
	"fmt"

	"github.com/berachain/stargazer/wasp"
	"github.com/ethereum/go-ethereum/core/types"
)

var _ wasp.BaseModel = (*EthTxnReceipt)(nil)

type EthTxnReceipt struct {
	wasp.BasePersistenceModal `gorm:"type:int;auto_increment;not null;"`
	Type                      uint8    `gorm:"type:int;not null;" json:"type,omitempty"`
	PostState                 []byte   `gorm:"type:bytea;" json:"root"`
	Status                    uint64   `gorm:"type:int;not null;" json:"status"`
	CumulativeGasUsed         uint64   `gorm:"type:int;not null;" json:"cumulativeGasUsed"`
	Bloom                     []byte   `gorm:"type:bytea;not null;" json:"logsBloom"`
	Logs                      []EthLog `gorm:"foreignkey:tx_hash;references:tx_hash" json:"logs"`
	TxHash                    []byte   `gorm:"type:bytea;not null;column:tx_hash;unique" json:"transactionHash"`
	ContractAddress           []byte   `gorm:"type:bytea;not null;" json:"contractAddress"`
	GasUsed                   uint64   `gorm:"type:int;not null;" json:"gasUsed"`
	BlockHash                 []byte   `gorm:"type:bytea;not null;" json:"blockHash"`
	BlockNumber               string   `gorm:"type:string;not null;" json:"blockNumber"`
	TransactionIndex          uint     `gorm:"type:int;not null;" json:"transactionIndex"`
}

func GethToReceiptModel(r *types.Receipt, logs []EthLog) *EthTxnReceipt {
	return &EthTxnReceipt{
		Type:              r.Type,
		PostState:         r.PostState,
		Status:            r.Status,
		CumulativeGasUsed: r.CumulativeGasUsed,
		Bloom:             r.Bloom.Bytes(),
		Logs:              logs,
		TxHash:            r.TxHash.Bytes(),
		ContractAddress:   r.ContractAddress.Bytes(),
		GasUsed:           r.GasUsed,
		BlockHash:         r.BlockHash.Bytes(),
		BlockNumber:       r.BlockNumber.String(),
		TransactionIndex:  r.TransactionIndex,
	}
}
func (m *EthTxnReceipt) GetId() int64 {
	return m.ID
}

func (m *EthTxnReceipt) GetTable() string {
	return "block_models"
}

func (m *EthTxnReceipt) GetRedisDb() int64 {
	return 1
}

func (m *EthTxnReceipt) GetRedisKey() string {
	key := fmt.Sprintf("%s:%s", m.GetTable(), m.BlockNumber)
	return key
}
