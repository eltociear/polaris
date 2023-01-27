package main

import (
	"github.com/berachain/stargazer/wasp/database"
	"github.com/berachain/stargazer/wasp/repository"
)

func main() {
	db, err := database.NewDatabase()
	if err != nil {
		panic(err)
	}
	repos := repository.InitRepositories(db)

	repos.BlockRepo.GetBlock()
}
