// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package ante

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// `SetAnteHandler` sets the required ante handler for a Stargazer Cosmos SDK Chain.
func SetAnteHandler(
	ak ante.AccountKeeper,
	bk authtypes.BankKeeper,
	fgk ante.FeegrantKeeper,
	txCfg client.TxConfig,
) func(bApp *baseapp.BaseApp) {
	return func(bApp *baseapp.BaseApp) {
		fmt.Println("SETTING TXCFG", txCfg.SignModeHandler())
		opt := ante.HandlerOptions{
			AccountKeeper:          ak,
			BankKeeper:             bk,
			ExtensionOptionChecker: extOptCheckerfunc,
			SignModeHandler:        txCfg.SignModeHandler(),
			FeegrantKeeper:         fgk,
			SigGasConsumer:         SigVerificationGasConsumer,
		}
		ch, _ := ante.NewAnteHandler(
			opt,
		)
		bApp.SetAnteHandler(
			ch,
		)
	}
}

func extOptCheckerfunc(a *codectypes.Any) bool {
	return a.TypeUrl == "/stargazer.evm.v1alpha1.ExtensionOptionsEthTransaction"
}
