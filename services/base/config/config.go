package config

import (
	"time"

	config "github.com/berachain/stargazer/services/base/server/config"
)

type (
	// `Server` defines the configuration for the JSON-RPC server.
	ServerConfig struct {
		// `API` defines a list of JSON-RPC namespaces to be enabled.
		EnableAPIs []string `mapstructure:"api"`

		// `Address` defines the HTTP server to listen on.
		Address string `mapstructure:"address"`

		// `WsAddress` defines the WebSocket server to listen on.
		WSAddress string `mapstructure:"ws-address"`

		// `MetricsAddress` defines the metrics server to listen on.
		MetricsAddress string `mapstructure:"metrics-address"`

		// `HTTPReadHeaderTimeout` is the read timeout of http json-rpc server.
		HTTPReadHeaderTimeout time.Duration `mapstructure:"http-read-header-timeout"`

		// `HTTPReadTimeout` is the read timeout of http json-rpc server.
		HTTPReadTimeout time.Duration `mapstructure:"http-read-timeout"`

		// `HTTPWriteTimeout` is the write timeout of http json-rpc server.
		HTTPWriteTimeout time.Duration `mapstructure:"http-write-timeout"`

		// HTTPIdleTimeout is the idle timeout of http json-rpc server.
		HTTPIdleTimeout time.Duration `mapstructure:"http-idle-timeout"`

		// `HTTPBaseRoute` defines the base path for the JSON-RPC server.
		BaseRoute string `mapstructure:"base-path"`

		// `TLSConfig` defines the TLS configuration for the JSON-RPC server.
		TLSConfig *TLSConfig `mapstructure:"tls-config"`
	}

	// `TLSConfig` defines a certificate and matching private key for the server.
	TLSConfig struct {
		// `CertPath` the file path for the certificate .pem file
		CertPath string `mapstructure:"cert-path"`

		// KeyPath the file path for the key .pem file
		KeyPath string `toml:"key-path"`
	}
)

type (
	// RPC defines RPC configuration of both the gRPC and CometBFT nodes.
	CosmosConnection struct {
		CMRPCEndpoint string `mapstructure:"cmrpc-endpoint" validate:"required"`
		GRPCEndpoint  string `mapstructure:"grpc-endpoint" validate:"required"`
		RPCTimeout    string `mapstructure:"rpc-timeout" validate:"required"`
		ChainID       string `mapstructure:"chain-id" validate:"required"`
	}
)

// `Config` is the configuration for the JSON-RPC service.
type Config struct {
	// `Server` is the configuration for the JSON-RPC server.
	Server config.ServerConfig
	// `Client` is the configuration for the Cosmos gRPC client.
	Client CosmosConnection
}

// `DefaultConfig` returns a default configuration for the JSON-RPC service.
func DefaultConfig() *Config {
	return &Config{
		Server: *DefaultServer(),
		Client: *DefaultCosmosConnection(),
	}
}
