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
	"github.com/berachain/stargazer/eth/common"
	"github.com/ethereum/go-ethereum/rlp"
)

// `StargazerHeader` represents a wrapped Ethereum header that allows for specifying a custom
// blockhash to make it compatible with a non-ethereum chain.
//
//go:generate rlpgen -type StargazerHeader -out header.rlpgen.go -decoder
type StargazerHeader struct {
	// `Header` is an embedded ethereum header.
	*Header
	// `hostHash` is the block hash on the host chain.
	hostHash common.Hash
}

// `NewEmptyStargazerHeader` returns an empty `StargazerHeader`.
func NewEmptyStargazerHeader() *StargazerHeader {
	return &StargazerHeader{Header: &Header{}}
}

// `NewStargazerHeader` returns a `StargazerHeader` with the given `header` and `hash`.
func NewStargazerHeader(header *Header, hash common.Hash) *StargazerHeader {
	return &StargazerHeader{Header: header, hostHash: hash}
}

// `Author` returns the address of the original block producer.
func (h *StargazerHeader) Author() common.Address {
	return h.Coinbase
}

// `UnmarshalBinary` decodes a block from the Ethereum RLP format.
func (h *StargazerHeader) UnmarshalBinary(data []byte) error {
	return rlp.DecodeBytes(data, h)
}

// `MarshalBinary` encodes the block into the Ethereum RLP format.
func (h *StargazerHeader) MarshalBinary() ([]byte, error) {
	bz, err := rlp.EncodeToBytes(h)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

// `Hash` returns the block hash of the header, we override the geth implementation
// to use the hash of the host chain, as the implementing chain might want to use it's
// real block hash opposed to hashing the "fake" header.
func (h *StargazerHeader) Hash() common.Hash {
	if h.hostHash == (common.Hash{}) {
		h.hostHash = h.Header.Hash()
	}
	return h.hostHash
}

// `SetHash` sets the hash of the header.
func (h *StargazerHeader) SetHash(hash common.Hash) {
	h.hostHash = hash
}
