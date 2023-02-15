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

package cosmos

import (
	"context"
	"os"
	"time"

	"github.com/berachain/stargazer/services/base/config"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	tmjsonclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"go.uber.org/zap"
)

// `Client` is a wrapper around the Cosmos SDK `client.Context` that implements querying and
// transaction capabilities for the Cosmos SDK.
type Client struct {
	// `ctx` is the context instance.
	ctx context.Context
	// `clientCtx` is the Cosmos SDK `client.Context` instance.
	clientCtx client.Context
	// `logger` is the logger instance.
	logger *zap.Logger
}

// `New` creates a new `CosmosClient`.
func New(
	ctx context.Context,
	cfg config.CosmosConnection,
	logger *zap.Logger,
) *Client {
	clientCtx, err := CreateClientContext(cfg)
	if err != nil {
		panic(err)
	}
	return &Client{
		ctx:       ctx,
		clientCtx: clientCtx,
		logger:    logger,
	}
}

// `CreateClientContext` creates a new `client.Context` instance.
func CreateClientContext(config config.CosmosConnection) (client.Context, error) {
	httpClient, err := tmjsonclient.DefaultHTTPClient(config.CMRPCEndpoint)
	if err != nil {
		return client.Context{}, err
	}

	httpClient.Timeout, err = time.ParseDuration(config.RPCTimeout)
	if err != nil {
		return client.Context{}, err
	}
	tmRPC, err := rpchttp.NewWithClient(config.CMRPCEndpoint, "/websocket", httpClient)
	if err != nil {
		return client.Context{}, err
	}

	clientCtx := client.Context{
		ChainID: config.ChainID,
		// InterfaceRegistry: oc.Encoding.InterfaceRegistry,
		Output:        os.Stderr,
		BroadcastMode: flags.BroadcastSync,
		// TxConfig:          oc.Encoding.TxConfig,
		// AccountRetriever:  authtypes.AccountRetriever{},
		// Codec:       oc.Encoding.Codec,
		// LegacyAmino: oc.Encoding.Amino,
		// Input:       os.Stdin,
		NodeURI: config.CMRPCEndpoint,
		Client:  tmRPC,
		// Keyring:      kr,
		// FromAddress:  oc.OracleAddr,
		// FromName:     keyInfo.Name,
		// From:         keyInfo.Name,
		OutputFormat: "json",
		UseLedger:    false,
		Simulate:     false,
		GenerateOnly: false,
		Offline:      false,
		SkipConfirm:  true,
	}

	return clientCtx, nil
}
