package cosmos

import (
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

// `LatestBlockNumber` returns the the latest block number as reported at the application layer.
func (c *Client) SignerBalance() (int64, error) {
	addr, err := c.clientCtx.Keyring.List()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(addr)
	request := &types.QueryBalanceRequest{
		Address: "addr",
		Denom:   "abera",
	}

	byteData, err := request.Marshal()
	if err != nil {
		return 0, err
	}
	res, code, err := c.clientCtx.QueryWithData("cosmos.bank.v1beta1.Query/AllBalances", byteData)
	if err != nil {
		return 0, err
	}
	if code == 0 {
		return 0, nil
	}

	fmt.Println(res)
	return 0, nil
}
