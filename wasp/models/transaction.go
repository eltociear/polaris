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
	Number                    string `gorm:"type:varchar(64);not null;column:number;"`
	Hash                      []byte `gorm:"type:bytea;not null;"`
	Size                      string `gorm:"type:varchar(64);not null"`
	Time                      uint64 `gorm:"type:int;not null;"`
	From                      []byte `gorm:"type:bytea;not null;"`
	Type                      uint8  `gorm:"type:smallint;not null;"`
	ChainID                   uint64 `gorm:"type:int;not null;"`
	Data                      []byte `gorm:"type:bytea;not null;"`
	Gas                       uint64 `gorm:"type:int;not null;"`
	GasPrice                  string `gorm:"type:varchar(64);not null;"`
	GasTipCap                 string `gorm:"type:varchar(64);not null;"`
	GasFeeCap                 string `gorm:"type:varchar(64);not null;"`
	Value                     string `gorm:"type:varchar(64);not null;"`
	Nonce                     uint64 `gorm:"type:int;not null;"`
	To                        []byte `gorm:"type:bytea;"`
}

func GethToTransactionModel(txn *types.Transaction, blockNumber string, time uint64, baseFee *big.Int, signer types.Signer) *TransactionModel {

	txnMsg, _ := txn.AsMessage(signer, baseFee)

	if txn.To() == nil {
		fmt.Println(txn.Hash().Hex())
		fmt.Println("CONTRACT CREATION")
		return &TransactionModel{
			Hash:      txn.Hash().Bytes(),
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
			To:        nil,
		}
	}
	return &TransactionModel{
		Hash:      txn.Hash().Bytes(),
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
