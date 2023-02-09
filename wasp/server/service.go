package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/berachain/stargazer/wasp/database"
	"github.com/berachain/stargazer/wasp/handler"
	"github.com/berachain/stargazer/wasp/queryClient"
	"github.com/berachain/stargazer/wasp/server/config"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Service struct {
	queryClient     *queryClient.QueryClient
	engine          *gin.Engine
	logger          *zap.Logger
	config          *config.Server
	shutdownTimeout time.Duration
	server          *http.Server
	notify          chan error
}

// `New` returns a new `Service` object.
func New() *Service {
	// Create a new logger instance.
	logger, _ := zap.NewProduction()
	defer logger.Sync() //nolint: errcheck /

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
	handler := handler.NewHandler(queryClient)

	engine := InitRouter(logger, handler)
	srv := &http.Server{
		Addr:    config.DefaultServer().Address,
		Handler: engine,
	}

	s := &Service{
		queryClient: queryClient,
		engine:      engine,
		config:      config.DefaultServer(),
		logger:      logger,
		server:      srv,
	}

	return s
}

// `Shutdown` stops the service.
func (s *Service) Shutdown() error {
	_, cancel := context.WithTimeout(
		context.Background(),
		s.shutdownTimeout,
	)
	defer cancel()
	// TODO: stop the gin server
	return nil
}

// `Start` starts the service.
func (s *Service) Start() error {
	ctx := context.Background()
	// listen and serve
	s.engine.Run(":10081")
	fmt.Printf("Server is listening on %s", s.server.Addr)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal or an error to gracefully shutdown the server.
	var err error
	select {
	case sig := <-interrupt:
		s.logger.Info(sig.String())
	case err = <-s.Notify():
		s.logger.Error(err.Error())
	}

	// Ensure that if the switch statement outputs an error, we return it to the CLI.
	if err != nil {
		return err
	}

	// Shutdown the server.
	if sErr := s.server.Shutdown(ctx); sErr != nil {
		s.logger.Error(sErr.Error())
		return sErr
	}

	return nil
}

func (s *Service) Notify() <-chan error {
	return s.notify
}
