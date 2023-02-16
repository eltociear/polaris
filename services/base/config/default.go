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

package config

import (
	"bufio"
	"os"
	"time"

	"github.com/evmos/ethermint/crypto/hd"
	"github.com/evmos/ethermint/encoding"

	config "github.com/berachain/stargazer/services/base/server/config"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

var (
	// `DefaultAPINamespaces` is the default namespaces the JSON-RPC server exposes.
	DefaultAPINamespaces = []string{"eth", "node"}
)

const (
	// `DefaultGRPCAddress` is the default address the gRPC server binds to.
	DefaultGRPCAddress = "0.0.0.0:9900"

	// `DefaultJSONRPCAddress` is the default address the JSON-RPC server binds to.
	DefaultJSONRPCAddress = "127.0.0.1:8545"

	// `DefaultJSONRPCWSAddress` is the default address the JSON-RPC WebSocket server binds to.
	DefaultJSONRPCWSAddress = "127.0.0.1:8546"

	// `DefaultJSOPNRPCMetricsAddress` is the default address the JSON-RPC Metrics server binds to.
	DefaultJSONRPCMetricsAddress = "127.0.0.1:6065"

	// `DefaultHTTPReadHeaderTimeout` is the default read timeout of http json-rpc server.
	DefaultHTTPReadHeaderTimeout = 5 * time.Second

	// `DefaultHTTPReadTimeout` is the default read timeout of http json-rpc server.
	DefaultHTTPReadTimeout = 10 * time.Second

	// `DefaultHTTPWriteTimeout` is the default write timeout of http json-rpc server.
	DefaultHTTPWriteTimeout = 10 * time.Second

	// `DefaultHTTPIdleTimeout` is the default idle timeout of http json-rpc server.
	DefaultHTTPIdleTimeout = 120 * time.Second

	// `DefaultBaseRoute` is the default base path for the JSON-RPC server.
	DefaultJSONRPCBaseRoute = "/"
)

// `DefaultServer` returns the default TLS configuration.
func DefaultServer() *config.ServerConfig {
	return &config.ServerConfig{
		EnableAPIs:            DefaultAPINamespaces,
		Address:               DefaultJSONRPCAddress,
		WSAddress:             DefaultJSONRPCWSAddress,
		MetricsAddress:        DefaultJSONRPCMetricsAddress,
		BaseRoute:             DefaultJSONRPCBaseRoute,
		HTTPReadHeaderTimeout: DefaultHTTPReadHeaderTimeout,
		HTTPReadTimeout:       DefaultHTTPReadTimeout,
		HTTPWriteTimeout:      DefaultHTTPWriteTimeout,
		HTTPIdleTimeout:       DefaultHTTPIdleTimeout,
		TLSConfig:             DefaultTLSConfig(),
	}
}

// DefaultConfig returns the default TLS configuration.
func DefaultTLSConfig() *config.TLSConfig {
	return &config.TLSConfig{
		CertPath: "",
		KeyPath:  "",
	}
}

const (
	// `DefaultCMRPCEndpoint` is the default address of the Comet RPC server.
	DefaultCMRPCEndpoint = "http://0.0.0.0:26657"

	// `DefaultRPCTimeout` is the default timeout for the RPC server.
	DefaultRPCTimeout = "10s"

	// `DefaultChainID` is the default chain ID.
	DefaultChainID = "berachain_420-1"

	DefaultKeyringFile        = ".keyring"
	DefaultKeyringDir         = "./"
	DefaultKeyringServiceName = "test"
)

// DefaultRPC returns the default RPC configuration.
func DefaultCosmosConnection() *CosmosConnection {
	return &CosmosConnection{
		CMRPCEndpoint: DefaultCMRPCEndpoint,
		GRPCEndpoint:  DefaultGRPCAddress,
		RPCTimeout:    DefaultRPCTimeout,
		ChainID:       DefaultChainID,
		Keyring:       nil,
	}
}

// DefaultRPC returns the default RPC configuration.
func DefaultSigningCosmosConnection() *CosmosConnection {
	DefaultEncodingConfig := encoding.MakeConfig(beramodules.GetModuleManager())
	KeyringOptions := []keyring.Option{hd.EthSecp256k1Option()}
	buf := bufio.NewReader(os.Stdin)

	DefaultKeyring, _ := keyring.New(DefaultKeyringServiceName, DefaultKeyringFile, DefaultKeyringDir, buf, DefaultEncodingConfig.Codec, KeyringOptions...)

	return &CosmosConnection{
		CMRPCEndpoint: DefaultCMRPCEndpoint,
		GRPCEndpoint:  DefaultGRPCAddress,
		RPCTimeout:    DefaultRPCTimeout,
		ChainID:       DefaultChainID,
		Keyring:       DefaultKeyring,
	}
}
