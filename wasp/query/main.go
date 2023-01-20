package main

import (
	"fmt"
	"os"

	"github.com/berachain/stargazer/wasp/backend"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

const numRootArgs = 0

var rootCmd = &cobra.Command{
	Use:   "iquery validators",
	Args:  cobra.MatchAll(cobra.ExactArgs(numRootArgs), cobra.OnlyValidArgs),
	Short: "Foundry contract generator",
	Run: func(cmd *cobra.Command, args []string) {
		grpcConn, err := grpc.Dial(
			"127.0.0.1:9090",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())),
		)
		if err != nil {
			fmt.Print("unable to connect")
			panic(err)
		}

		backend := backend.NewBackend(grpcConn, true)
		res, err := backend.GetValidators()
		if err != nil {
			fmt.Print("unable to query")
			panic(err)
		}
		fmt.Printf("%v\n", res)
	},
}
