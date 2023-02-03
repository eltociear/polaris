package repository

import (
	"github.com/berachain/stargazer/wasp/database"
	"github.com/berachain/stargazer/wasp/queryClient"
)

type AccountRepo struct {
	db *database.Database
	qc *queryClient.QueryClient
}

func NewAccountRepo(db *database.Database, qc *queryClient.QueryClient) *AccountRepo {
	return &AccountRepo{
		db: db,
		qc: qc,
	}
}

func (r *AccountRepo) UpdateAccounts() {

}
func (r *AccountRepo) ParseTransfers() {

}
func (r *AccountRepo) UpdateBalances() {

}
