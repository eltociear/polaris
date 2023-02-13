package models

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/berachain/stargazer/wasp"
	"github.com/ethereum/go-ethereum/core/types"
)

var _ wasp.BaseModel = (*TransactionModel)(nil)

type TransactionModel struct {
	wasp.BasePersistenceModal `gorm:"type:int;auto_increment;not null;"`
	Number                    string        `gorm:"type:varchar(64);not null;column:number;"`
	Hash                      string        `gorm:"type:varchar(128);not null;column:tx_hash"`
	Size                      string        `gorm:"type:varchar(64);not null"`
	Time                      uint64        `gorm:"type:int;not null;"`
	From                      []byte        `gorm:"type:bytea;not null;"`
	Type                      uint8         `gorm:"type:int;not null;"`
	ChainID                   uint64        `gorm:"type:int;not null;"`
	Data                      []byte        `gorm:"type:bytea;not null;"`
	Gas                       uint64        `gorm:"type:int;not null;"`
	GasPrice                  string        `gorm:"type:varchar(64);not null;"`
	GasTipCap                 string        `gorm:"type:varchar(64);not null;"`
	GasFeeCap                 string        `gorm:"type:varchar(64);not null;"`
	Value                     string        `gorm:"type:varchar(64);not null;"`
	Nonce                     uint64        `gorm:"type:int;not null;"`
	To                        []byte        `gorm:"type:bytea;"`
	Receipt                   EthTxnReceipt `gorm:"foreignkey:tx_hash;references:tx_hash" json:"receipt"`
}

func GethToTransactionModel(
	txn *types.Transaction,
	blockNumber string,
	time uint64,
	baseFee *big.Int,
	signer types.Signer,
	receipt EthTxnReceipt) *TransactionModel {

	txnMsg, _ := txn.AsMessage(signer, baseFee)
	from, _ := signer.Sender(txn)
	if txn.To() == nil {
		return &TransactionModel{
			Hash:      txn.Hash().Hex(),
			Size:      txn.Size().String(),
			Time:      time,
			From:      from.Bytes(),
			Type:      txn.Type(),
			ChainID:   txn.ChainId().Uint64(),
			Data:      txn.Data(),
			Gas:       txn.Gas(),
			GasPrice:  txn.GasPrice().String(),
			GasTipCap: txn.GasTipCap().String(),
			GasFeeCap: txn.GasFeeCap().String(),
			Value:     txn.Value().String(),
			Nonce:     txn.Nonce(),
			To:        nil,
			Receipt:   receipt,
		}
	}
	return &TransactionModel{
		Hash:      txn.Hash().Hex(),
		Size:      txn.Size().String(),
		Time:      time,
		From:      txnMsg.From().Bytes(),
		Type:      txn.Type(),
		ChainID:   txn.ChainId().Uint64(),
		Data:      txn.Data(),
		Gas:       txn.Gas(),
		GasPrice:  txn.GasPrice().String(),
		GasTipCap: txn.GasTipCap().String(),
		GasFeeCap: txn.GasFeeCap().String(),
		Value:     txn.Value().String(),
		Nonce:     txn.Nonce(),
		To:        txn.To().Bytes(),
		Receipt:   receipt,
	}
}
func (m *TransactionModel) GetId() int64 {
	return m.ID
}

func (m *TransactionModel) GetTable() string {
	return "block_models"
}

func (m *TransactionModel) GetRedisDb() int64 {
	return 1
}

func (m *TransactionModel) GetRedisKey() string {
	key := fmt.Sprintf("%s:%s", m.GetTable(), strconv.FormatInt(10, 10))
	return key
}
