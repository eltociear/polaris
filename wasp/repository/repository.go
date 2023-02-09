package repository

import (
	"github.com/berachain/stargazer/wasp/database"
	"github.com/berachain/stargazer/wasp/queryClient"
)

type Repositories struct {
	BlockRepo       *BlockRepo
	TransactionRepo *TransactionRepo
	AccountRepo     *AccountRepo
}

func InitRepositories(db *database.Database, qc *queryClient.QueryClient) *Repositories {
	return &Repositories{
		BlockRepo:       NewBlockRepo(db, qc),
		TransactionRepo: NewTransactionRepo(db, qc),
		AccountRepo:     NewAccountRepo(db, qc),
	}
}
