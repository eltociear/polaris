package api

import (
	"fmt"

	"github.com/berachain/stargazer/services/base/cosmos"
)

type FaucetApi struct {
	client *cosmos.Client
}

func NewFaucetApi(client *cosmos.Client) *FaucetApi {
	return &FaucetApi{
		client: client,
	}
}

func (api *FaucetApi) GetBalance() {
	res, err := api.client.SignerBalance()
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
