package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/berachain/stargazer/wasp/handler"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRouter(logger *zap.Logger, handler *handler.Handler) *gin.Engine {
	fmt.Println("Initializing router...")
	engine := gin.Default()
	engine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	engine.Use(gin.Recovery())
	engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	engine.GET("/block", handler.GetLatestBlock)
	engine.GET("/block/:height", handler.GetBlock)
	engine.GET("/blocks", handler.GetBlocks)

	engine.GET("transaction/hash/:hash", handler.GetTransactionByHash)
	engine.GET("transaction/block/hash/:hash", handler.GetTransactionsByBlockHash)
	engine.GET("transaction/block/number/:number", handler.GetTransactionsByBlockNumber)
	engine.GET("transaction/receipt/:hash", handler.GetTransactionReceipt)

	engine.GET("transactions", handler.GetLatestTransactions)
	engine.GET("transactions/count", handler.GetTransactionCount)

	return engine
}
