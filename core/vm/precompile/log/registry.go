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

package log

import (
	"github.com/berachain/stargazer/lib/common"
	"github.com/berachain/stargazer/lib/errors"
	"github.com/berachain/stargazer/types/abi"
)

// `Registry` stores a mapping of `EventType`s to `*log.PrecompileLog`s.
type Registry struct {
	// `eventTypesToLogs` is a map of `EventType`s to `*log.PrecompileLog` for all supported
	// events.
	eventTypesToLogs map[string]*PrecompileLog

	// `factory` is the `LogFactory` used to create `sdk.Event`s. ( //
	Translator Translator // todo: generalize
}

// `NewRegistry` creates and returns a new, empty `Registry`.
func NewRegistry(
	translator Translator,
) *Registry {
	return &Registry{
		eventTypesToLogs: make(map[string]*PrecompileLog),
		Translator:       translator,
	}
}

// `RegisterEvent` registers an Ethereum event as a precompile log.
func (lr *Registry) RegisterEvent(
	moduleEthAddress common.Address,
	abiEvent abi.Event,
) error {
	if _, found := lr.eventTypesToLogs[abiEvent.Name]; found {
		return errors.Wrap(ErrEthEventAlreadyRegistered, abiEvent.Name)
	}

	var err error
	lr.eventTypesToLogs[abiEvent.Name], err = NewPrecompileLog(
		moduleEthAddress,
		abiEvent,
	)

	// todo register translator events here?

	return err
}

// `GetPrecompileLog` returns the precompile log corresponding to the given event.
func (lr *Registry) GetPrecompileLog(eventType string) *PrecompileLog {
	return lr.eventTypesToLogs[abi.ToCamelCase(eventType)]
}
