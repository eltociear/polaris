package main

import (
	"github.com/berachain/stargazer/wasp/database"
)

func main() {
	db, err := database.NewDatabase()
	if err != nil {
		panic(err)
	}
	// repos := repository.InitRepositories(db)

}
