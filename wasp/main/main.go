package main

import (
	"fmt"
	"log"

	"github.com/berachain/stargazer/wasp/database"
	"github.com/berachain/stargazer/wasp/queryClient"
	"github.com/berachain/stargazer/wasp/repository"
	"github.com/berachain/stargazer/wasp/syncr"
	"github.com/ethereum/go-ethereum/ethclient"
)

var exit = make(chan bool)

func main() {
	db, err := database.NewDatabase()
	if err != nil {
		panic(err)
	}

	// client, err := ethclient.Dial(os.Getenv("RPC_ENDPOINT"))
	client, err := ethclient.Dial("wss://eth-goerli.g.alchemy.com/v2/2Vd54oL5HObq1Yl_aZfLZzBz37_FCNdP")

	if err != nil {
		log.Fatal(err)
	}
	queryClient := queryClient.NewQueryClient(client, db.Gorm)
	repos := repository.InitRepositories(db, queryClient)

	syncrClient := syncr.NewSyncrClient(client, repos)
	go syncrClient.Run()
	<-exit // This blocks until the exit channel receives some input
	fmt.Println("Done.")
}
