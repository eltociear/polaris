package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/berachain/stargazer/wasp/database"
	"github.com/berachain/stargazer/wasp/proto"
	"github.com/berachain/stargazer/wasp/repository"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var exit = make(chan bool)

func main() {
	db, err := database.NewDatabase()
	if err != nil {
		panic(err)
	}

	repos := repository.InitRepositories(db)

	go syncr(repos)
	<-exit // This blocks until the exit channel receives some input
	fmt.Println("Done.")
}

func syncr(repos *repository.Repositories) {
	ctx := context.Background()
	client, err := ethclient.Dial(os.Getenv("RPC_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(ctx, headers)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Listening...")
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}
			msg := &proto.CreateBlockRequest{
				Block: &proto.Block{
					Number: block.Number().Int64(),
				},
			}
			res := repos.BlockRepo.CreateBlock(ctx, msg)
			fmt.Println(block.Number().Int64()) // 3477413
			fmt.Println(res.Code)               // 3477413

		}
	}
}
