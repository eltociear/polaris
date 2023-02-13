package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/berachain/stargazer/wasp"
	"github.com/ethereum/go-ethereum/core/types"
)

var _ wasp.BaseModel = (*EthBlockModel)(nil)

type EthBlockModel struct {
	wasp.BasePersistenceModal `gorm:"type:int;auto_increment;not null;"`
	ParentHash                []byte             `gorm:"type:bytea;not null;" json:"parentHash"`
	UncleHash                 []byte             `gorm:"type:bytea;not null;" json:"sha3Uncles"`
	Coinbase                  []byte             `gorm:"type:bytea;not null;" json:"miner"`
	Root                      []byte             `gorm:"type:bytea;not null;" json:"stateRoot"`
	TxHash                    []byte             `gorm:"type:bytea;not null;" json:"transactionHash"`
	ReceiptHash               []byte             `gorm:"type:bytea;not null;" json:"receiptsRoot"`
	Bloom                     []byte             `gorm:"type:bytea;not null;" json:"logsBloom"`
	Difficulty                string             `gorm:"type:varchar(64);not null;" json:"difficulty"`
	BlockNumber               string             `gorm:"type:varchar(64);not null;column:block_number;unique" json:"number"`
	GasLimit                  string             `gorm:"type:varchar(64);not null;" json:"gasLimit"`
	GasUsed                   string             `gorm:"type:varchar(64);not null;" json:"gasUsed"`
	Time                      string             `gorm:"type:varchar(64);not null;" json:"timestamp"`
	Extra                     []byte             `gorm:"type:bytea;not null;" json:"extraData"`
	MixDigest                 []byte             `gorm:"type:bytea;not null;" json:"mixHash"`
	Nonce                     uint64             `gorm:"type:int;not null;" json:"nonce"`
	BaseFee                   []byte             `gorm:"type:bytea;not null;" json:"baseFeePerGas"`
	Txs                       []TransactionModel `gorm:"foreignKey:Number;references:block_number"`
	Hash                      []byte             `gorm:"type:bytea;not null;unique" json:"hash"`
	Size                      string             `gorm:"type:varchar(64);not null;" json:"size"`
	ReceivedAt                time.Time          `gorm:"type:timestamp;not null;" json:"receivedAt"`
}

func GethToBlockModel(gethBlock *types.Block, txns []TransactionModel) *EthBlockModel {

	return &EthBlockModel{
		ParentHash:  gethBlock.ParentHash().Bytes(),
		UncleHash:   gethBlock.UncleHash().Bytes(),
		Coinbase:    gethBlock.Coinbase().Bytes(),
		Root:        gethBlock.Root().Bytes(),
		TxHash:      gethBlock.TxHash().Bytes(),
		ReceiptHash: gethBlock.ReceiptHash().Bytes(),
		Bloom:       gethBlock.Bloom().Big().Bytes(),
		Difficulty:  gethBlock.Difficulty().String(),
		BlockNumber: gethBlock.Number().String(),
		GasLimit:    strconv.FormatUint(gethBlock.GasLimit(), 10),
		GasUsed:     strconv.FormatUint(gethBlock.GasUsed(), 10),
		Time:        strconv.FormatUint(gethBlock.Time(), 10),
		Extra:       gethBlock.Extra(),
		MixDigest:   gethBlock.MixDigest().Bytes(),
		Nonce:       gethBlock.Nonce(),
		BaseFee:     gethBlock.BaseFee().Bytes(),
		Hash:        gethBlock.Hash().Bytes(),
		Size:        gethBlock.Size().String(),
		ReceivedAt:  gethBlock.ReceivedAt,
		// ReceivedFrom: gethBlock.ReceivedFrom,
		Txs: txns,
	}
}

func (m *EthBlockModel) GetId() int64 {
	return m.ID
}

func (m *EthBlockModel) GetTable() string {
	return "block_models"
}

func (m *EthBlockModel) GetRedisDb() int64 {
	return 1
}

func (m *EthBlockModel) GetRedisKey() string {
	key := fmt.Sprintf("%s:%s", m.GetTable(), m.BlockNumber)
	return key
}
