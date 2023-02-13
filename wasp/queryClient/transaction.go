package queryClient

import (
	"context"

	"github.com/berachain/stargazer/wasp"
	"github.com/berachain/stargazer/wasp/models"
)

func (c *QueryClient) GetTransactionByHash(ctx context.Context, hash string) (*models.TransactionModel, error) {
	txnClient := c.db.TransactionModel
	tx, err := txnClient.WithContext(ctx).Preload(txnClient.Receipt).Where(txnClient.Hash.Eq(hash)).First()
	return tx, err
}
func (c *QueryClient) GetLatestTransactions(ctx context.Context, pagination wasp.Pagination) ([]*models.TransactionModel, error) {
	txnClient := c.db.TransactionModel
	offset := (pagination.Page - 1) * pagination.Limit
	tx, err := txnClient.WithContext(ctx).Preload(txnClient.Receipt).Offset(offset).Limit(pagination.Limit).Find()
	return tx, err
}
func (c *QueryClient) GetTransactionsByBlockNumber(ctx context.Context, blockNumber string, pagination wasp.Pagination) ([]*models.TransactionModel, error) {
	txnClient := c.db.TransactionModel
	offset := (pagination.Page - 1) * pagination.Limit
	tx, err := txnClient.WithContext(ctx).Preload(txnClient.Receipt).Where(txnClient.Number.Eq(blockNumber)).Offset(offset).Limit(pagination.Limit).Find()
	return tx, err
}
func (c *QueryClient) GetTransactionsByBlockHash(ctx context.Context, hash string, pagination wasp.Pagination) ([]*models.TransactionModel, error) {
	txnClient := c.db.TransactionModel
	offset := (pagination.Page - 1) * pagination.Limit
	tx, err := txnClient.WithContext(ctx).Preload(txnClient.Receipt).Where(txnClient.Hash.Eq(hash)).Offset(offset).Limit(pagination.Limit).Find()
	return tx, err
}
func (c *QueryClient) GetTransactionCount(ctx context.Context) (int64, error) {
	txnClient := c.db.TransactionModel
	count, err := txnClient.WithContext(ctx).Select(txnClient.ALL).Count()
	return count, err
}
func (c *QueryClient) GetTransactionReceipt(ctx context.Context, hash string) (*models.EthTxnReceipt, error) {
	receiptClient := c.db.EthTxnReceipt
	receipt, err := receiptClient.WithContext(ctx).Where(receiptClient.TxHash.Eq(hash)).First()
	return receipt, err
}
