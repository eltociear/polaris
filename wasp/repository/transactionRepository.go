package repository

import (
	"context"
	"math/big"

	"github.com/berachain/stargazer/wasp/database"
	"github.com/berachain/stargazer/wasp/models"
	"github.com/berachain/stargazer/wasp/queryClient"

	"github.com/ethereum/go-ethereum/core/types"
)

type TransactionRepo struct {
	db *database.Database
	qc *queryClient.QueryClient
}

func NewTransactionRepo(db *database.Database, qc *queryClient.QueryClient) *TransactionRepo {
	return &TransactionRepo{
		db: db,
		qc: qc,
	}
}

func (r *TransactionRepo) BuildTransactionList(block *types.Block) *[]models.TransactionModel {
	txns := []models.TransactionModel{}
	chainID, err := r.qc.NetworkID(context.Background())
	if err != nil {
		panic("unable to retrieve chainId")
	}
	signerType := types.NewEIP155Signer(chainID)

	for _, t := range block.Transactions() {
		txn := *r.BuildTransaction(
			t,
			block.Number().String(),
			block.Time(),
			block.BaseFee(),
			signerType)
		txns = append(txns, txn)
	}
	return &txns
}
func (r *TransactionRepo) BuildTransaction(
	txn *types.Transaction,
	blockNumber string,
	time uint64,
	baseFee *big.Int,
	signerType types.Signer) *models.TransactionModel {

	receipt := r.BuildTransactionReceipt(txn, blockNumber, time, baseFee)
	txnModel := models.GethToTransactionModel(
		txn,
		blockNumber,
		time,
		baseFee,
		signerType,
		receipt)

	return txnModel
}

func (r *TransactionRepo) BuildTransactionReceipt(
	txn *types.Transaction,
	blockNumber string,
	time uint64,
	baseFee *big.Int) models.EthTxnReceipt {
	gethReciept, err := r.qc.GetTransactionReceiptByHash(txn.Hash())
	if err != nil {
		panic(err)
	}
	ethTxnLogs := r.BuildTransactionLogs(gethReciept)
	txnReceiptModel := models.GethToReceiptModel(gethReciept, ethTxnLogs)
	return *txnReceiptModel
}

func (r *TransactionRepo) BuildTransactionLogs(receipt *types.Receipt) []models.EthLog {
	logs := []models.EthLog{}
	for _, log := range receipt.Logs {
		logModel := *models.GethToEthLogModel(log)
		logs = append(logs, logModel)
	}
	return logs
}
