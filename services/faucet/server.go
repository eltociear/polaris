package faucet

import (
	"context"
	"net/http"
	"time"

	"github.com/berachain/stargazer/services/base/config"
	"github.com/berachain/stargazer/services/base/cosmos"
	server "github.com/berachain/stargazer/services/base/server"
	"github.com/berachain/stargazer/services/faucet/api"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var _ server.Server = &FaucetServer{}

type FaucetServer struct {
	service *server.Service
}

func NewFaucetServer() *FaucetServer {
	ctx := context.Background()
	logger, _ := zap.NewProduction()
	defer logger.Sync() //nolint: errcheck // ignore error
	// Create a new Cosmos client to connect to the node.
	cfg := config.DefaultSigningConfig()
	// Create new cosmos client.
	client := cosmos.New(ctx, cfg.Client, logger)

	// Create new service
	service := server.NewService(ctx, logger, client, cfg.Server)
	// Add logging middleware.
	service.GetEngine().Use(ginzap.Ginzap(logger, time.RFC3339, true))
	// Add tracing middleware.
	service.GetEngine().Use(gin.Recovery())
	// Configure the JSON-RPC API.
	faucetServer := &FaucetServer{
		service: service,
	}

	// Register routes.
	faucetServer.RegisterAPI()
	return faucetServer
}

func (s *FaucetServer) Start() {
	go func() {
		s.service.GetLogger().Info("Starting Faucer server at:", zap.String("address", s.service.GetConfig().Address))
		s.service.GetNotify() <- s.service.GetEngine().Run(s.service.GetConfig().Address)
		close(s.service.GetNotify())
	}()
}

func (s *FaucetServer) Shutdown() error {
	// Set a timeout for the shutdown.
	_, cancel := context.WithTimeout(
		context.Background(),
		s.service.GetShutdownTimeout(),
	)
	defer cancel()
	return nil
}

func (s *FaucetServer) Notify() <-chan error {
	return s.service.GetNotify()
}

func (s *FaucetServer) RegisterAPI() {
	api := api.NewFaucetApi(s.service.GetCosmosClient())
	s.service.GetEngine().GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	s.service.GetEngine().GET("/ding", func(c *gin.Context) {
		api.GetBalance()
		c.String(http.StatusOK, "pong")
	})
}
