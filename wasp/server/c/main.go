package main

import (
	"context"
	"time"

	"github.com/berachain/stargazer/wasp/server"
)

func main() {
	a, b := server.NewContainer(context.Background(), server.DefaultContainerConfig())
	if b != nil {
		panic(b)
	}
	x := time.Second
	err := a.Stop(context.Background(), &x)
	if err != nil {
		panic(err)
	}
}
