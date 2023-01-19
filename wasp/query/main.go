package main

import (
	"flag"
	"fmt"

	"github.com/berachain/stargazer/wasp/backend"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = flag.String("addr", "localhost:26657", "the address to connect to")

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Print("unable to connect")
		panic(err)
	}

	backend := backend.NewBackend(conn, true)
	res, err := backend.GetValidators()
	if err != nil {
		fmt.Print("unable to read \n")
		panic(err)
	}

	fmt.Println(res)
}
