package jsonrpc

import (
	"context"
	"time"

	ethlog "github.com/berachain/stargazer/eth/log"
	"github.com/berachain/stargazer/services/base/config"
	"github.com/berachain/stargazer/services/base/cosmos"
	server "github.com/berachain/stargazer/services/base/server"
	"github.com/berachain/stargazer/services/jsonrpc/api"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var _ server.Server = &JsonRpcServer{}

type JsonRpcServer struct {
	service *server.Service
	// `rpcserver` is the externally facing JSON-RPC Server.
	rpcserver *ethrpc.Server
}

func NewJsonRpcServer() *JsonRpcServer {

	ctx := context.Background()
	rpcserver := ethrpc.NewServer()
	// Create a new logger instance.
	logger, _ := zap.NewProduction()
	defer logger.Sync() //nolint: errcheck // ignore error

	// Create a new Cosmos client to connect to the node.
	cfg := config.DefaultConfig()
	// Create new cosmos client.
	client := cosmos.New(ctx, cfg.Client, logger)
	// Create new service
	service := server.NewService(ctx, logger, client, cfg.Server)
	// Configure the JSON-RPC API.
	service.GetEngine().Use(ginzap.Ginzap(logger, time.RFC3339, true))
	// Configure the JSON-RPC API.
	ethlog.Root().SetHandler(ethlog.FuncHandler(func(r *ethlog.Record) error {
		sugared := logger.Sugar()
		switch r.Lvl { //nolint:nolintlint,exhaustive // linter is bugged.
		case ethlog.LvlTrace, ethlog.LvlDebug:
			sugared.Debug(r.Msg, r.Ctx)
		case ethlog.LvlInfo, ethlog.LvlWarn:
			sugared.Info(r.Msg, r.Ctx)
		case ethlog.LvlError, ethlog.LvlCrit:
			sugared.Error(r.Msg, r.Ctx)
		}
		return nil
	}))
	// Configure the JSON-RPC API.
	service.GetEngine().Use(gin.Recovery())
	// Configure the JSON-RPC API.
	service.GetEngine().Any(service.GetConfig().BaseRoute, gin.WrapH(rpcserver))

	// Create the JSON-RPC server.
	jsonRpcServer := &JsonRpcServer{
		service:   service,
		rpcserver: rpcserver,
	}

	// Register the JSON-RPC API.
	for _, namespace := range cfg.Server.EnableAPIs {
		if err := jsonRpcServer.RegisterAPI(api.Build(namespace, service.GetCosmosClient(), service.GetLogger())); err != nil {
			panic(err)
		}
	}

	return jsonRpcServer
}

func (s *JsonRpcServer) Start() {
	// Start the JSON-RPC server.
	go func() {
		s.service.GetLogger().Info("Starting JSON-RPC server at:", zap.String("address", s.service.GetConfig().Address))
		s.service.GetNotify() <- s.service.GetEngine().Run(s.service.GetConfig().Address)
		close(s.service.GetNotify())
	}()
}

func (s *JsonRpcServer) Shutdown() error {
	// Set a timeout for the shutdown.
	_, cancel := context.WithTimeout(
		context.Background(),
		s.service.GetShutdownTimeout(),
	)
	defer cancel()
	// Stop the RPC Server
	s.rpcserver.Stop()
	// TODO: stop the gin server
	return nil
}

func (s *JsonRpcServer) Notify() <-chan error {
	return s.service.GetNotify()
}

func (s *JsonRpcServer) RegisterAPI(service api.Service) error {
	return s.rpcserver.RegisterName(service.Namespace(), service)
}
