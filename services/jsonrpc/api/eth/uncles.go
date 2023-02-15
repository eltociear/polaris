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

package eth

import (
	"github.com/berachain/stargazer/lib/common"
	"github.com/berachain/stargazer/lib/common/hexutil"
)

// `GetUncleByBlockHashAndIndex` returns nil since there are no uncles in Tendermint.
func (*api) GetUncleByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) map[string]interface{} {
	return nil
}

// `GetUncleByBlockNumberAndIndex` returns nil since there are no uncles in Tendermint.
func (*api) GetUncleByBlockNumberAndIndex(number, idx hexutil.Uint) map[string]interface{} {
	return nil
}

// `GetUncleCountByBlockHash` returns 0 since there are no uncles in Tendermint.
func (*api) GetUncleCountByBlockHash(hash common.Hash) hexutil.Uint {
	return 0
}

// `GetUncleCountByBlockNumber` returns 0 since there are no uncles in Tendermint.
func (*api) GetUncleCountByBlockNumber(blockNum int64) hexutil.Uint {
	return 0
}
