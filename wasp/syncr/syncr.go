package syncr

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/berachain/stargazer/wasp/repository"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type SyncrClient struct {
	client *ethclient.Client
	repos  *repository.Repositories
}

func NewSyncrClient(client *ethclient.Client, repos *repository.Repositories) *SyncrClient {
	return &SyncrClient{
		client: client,
		repos:  repos,
	}
}

func (s *SyncrClient) Run() {
	ctx := context.Background()

	headers := make(chan *types.Header)
	sub, _ := s.client.SubscribeNewHead(ctx, headers)

	fmt.Print("Listening...")
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			// latest block - 1
			block, err := s.client.BlockByNumber(context.Background(), big.NewInt(0).Sub(header.Number, big.NewInt(1)))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("BLOCK NUMBER: %s\n", block.Number().String())
			s.parseBlock(ctx, block)
		}

	}
}

func (s *SyncrClient) parseBlock(ctx context.Context, block *types.Block) {
	parsedTxns := s.repos.TransactionRepo.BuildTransactionList(block)
	parsedBlock := s.repos.BlockRepo.BuildBlock(block, *parsedTxns)
	s.repos.BlockRepo.CreateBlock(ctx, parsedBlock)
}
