package repository

import (
	"github.com/berachain/stargazer/wasp/database"
)

type Repositories struct {
	BlockRepo *BlockRepo
}

func InitRepositories(db *database.Database) *Repositories {
	return &Repositories{
		BlockRepo: NewBlockRepo(db),
	}
}
