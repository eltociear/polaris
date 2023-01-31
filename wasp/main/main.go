package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/berachain/stargazer/wasp/database"
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
	chainID, err := client.NetworkID(context.Background())
	signerType := types.NewEIP155Signer(chainID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Listening...")
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			block, err := client.BlockByNumber(context.Background(), big.NewInt(0).Sub(header.Number, big.NewInt(1)))
			if err != nil {
				log.Fatal(err)
			}
			code := repos.BlockRepo.CreateBlock(ctx, block, signerType)
			fmt.Println(code)
		}

	}
}
