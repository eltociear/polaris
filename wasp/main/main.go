package main

import (
	"fmt"
	"log"
	"os"

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
	client, err := ethclient.Dial(os.Getenv("RPC_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}
	queryClient := queryClient.NewQueryClient(client)
	repos := repository.InitRepositories(db, queryClient)

	syncrClient := syncr.NewSyncrClient(client, repos)
	go syncrClient.Run()
	<-exit // This blocks until the exit channel receives some input
	fmt.Println("Done.")
}
