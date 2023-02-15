// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// See the file LICENSE for licensing terms.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package server

import (
	"context"
	"time"

	"github.com/berachain/stargazer/services/base/cosmos"
	config "github.com/berachain/stargazer/services/base/server/config"
	server "github.com/berachain/stargazer/services/base/server/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server interface {
	Start()
	Shutdown() error
	Notify() <-chan error
}

type Service struct {
	// `cosmosClient` provides the gRPC connection to the Cosmos node.
	cosmosClient *cosmos.Client
	// `engine` is the gin engine responsible for handling the JSON-RPC requests.
	engine *gin.Engine
	// `logger` is the logger for the service.
	logger *zap.Logger
	// `notify` is the channel that is used to notify the service has stopped.
	notify chan error
	// `shutdownTimeout` is the delay between the service being stopped and the HTTP server being shutdown.
	shutdownTimeout time.Duration
	// config is the configuration for the service.
	config server.ServerConfig
}

// `New` returns a new `Service` object.
func NewService(ctx context.Context, logger *zap.Logger, client *cosmos.Client, config server.ServerConfig) *Service {
	// Configure the JSON-RPC API.
	return &Service{
		cosmosClient: client,
		config:       config,
		logger:       logger,
		notify:       make(chan error, 1),
		engine:       gin.Default(),
	}
}

func (s *Service) GetEngine() *gin.Engine {
	return s.engine
}

func (s *Service) GetCosmosClient() *cosmos.Client {
	return s.cosmosClient
}

func (s *Service) GetConfig() config.ServerConfig {
	return s.config
}

func (s *Service) GetLogger() *zap.Logger {
	return s.logger
}

func (s *Service) GetNotify() chan error {
	return s.notify
}

func (s *Service) GetShutdownTimeout() time.Duration {
	return s.shutdownTimeout
}
