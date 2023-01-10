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

package types

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	proto "github.com/cosmos/gogoproto/proto"
)

// https://github.com/ethereum/go-ethereum/issues/23154#issuecomment-876254171
// We expose TxData for only one reason: making creation of types.Transaction
// from Go code more convenient. The methods of TxData are an internal implementation detail
// and will never be stable API. It's not a good idea to send these fields over protocol buffers.
//  If your app needs to send/receive Ethereum transactions over protocol buffers,
// it should accept them as opaque binary data, which can be parsed by Transaction.UnmarshalBinary.

// NOTE: All non-protected transactions (i.e non EIP155 signed) will fail if the
// AllowUnprotectedTxs parameter is disabled.
func NewTxDataFromTx(tx *Transaction) (TxData, error) {
	var err error
	var inner TxData
	switch tx.Type() {
	case uint8(DynamicFeeTxType):
		// txData, err = newDynamicFeeTx(tx)
	case uint8(AccessListTxType):
		// txData, err = newAccessListTx(tx)
	default:
		// txData, err = newLegacyTx(tx)
	}
	// if err != nil {
	// 	return nil, err
	// }

	return inner, err
}

func BytesToEthereumTransaction(b []byte) (*Transaction, error) {
	tx := new(Transaction)
	err := tx.UnmarshalBinary(b)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// ==============================================================================
// TxData <-> Proto.Any
// ==============================================================================

// PackTxData constructs a new Any packed with the given tx data value. It returns
// an error if the client state can't be casted to a protobuf message or if the concrete
// implementation is not registered to the protobuf codec.
func PackTxData(txData TxData) (*codectypes.Any, error) {
	msg, ok := txData.(proto.Message)
	if !ok {
		return nil, errorsmod.Wrapf(errortypes.ErrPackAny, "cannot proto marshal %T", txData)
	}

	anyTxData, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, errorsmod.Wrap(errortypes.ErrPackAny, err.Error())
	}

	return anyTxData, nil
}

// UnpackTxData unpacks an Any into a TxData. It returns an error if the
// client state can't be unpacked into a TxData.
func UnpackTxData(data *codectypes.Any) (TxData, error) {
	if data == nil {
		return nil, errorsmod.Wrap(errortypes.ErrUnpackAny, "protobuf Any message cannot be nil")
	}

	txData, ok := data.GetCachedValue().(TxData)
	if !ok {
		return nil, errorsmod.Wrapf(errortypes.ErrUnpackAny, "cannot unpack Any into TxData %T", data)
	}

	return txData, nil
}

// deriveChainId derives the chain id from the given v parameter.
func DeriveChainID(v *big.Int) *big.Int {
	if v.BitLen() <= 64 { //nolint: gomnd // from geth
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((v - 35) / 2) //nolint: gomnd // from geth
	}
	v = new(big.Int).Sub(v, big.NewInt(35)) //nolint: gomnd // from geth
	return v.Div(v, big.NewInt(2))          //nolint: gomnd // from geth
}
