package main

import (
	"context"
	"fmt"

	"github.com/berachain/stargazer/wasp/database"
	"github.com/berachain/stargazer/wasp/proto"
	"github.com/berachain/stargazer/wasp/repository"
)

func main() {
	ctx := context.Background()
	db, err := database.NewDatabase()
	if err != nil {
		panic(err)
	}
	repos := repository.InitRepositories(db)

	msg := &proto.ReadBlockRequest{
		Number: 1,
	}
	res := repos.BlockRepo.GetBlock(ctx, msg)
	fmt.Print(res.GetBlock().Number)

	// msg := &proto.CreateBlockRequest{
	// 	Block: &proto.Block{
	// 		Number: 1,
	// 	},
	// }
	// res := repos.BlockRepo.CreateBlock(ctx, msg)
	// fmt.Printf("\n%d", res.Code)
}
