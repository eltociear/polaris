package multicall

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/berachain/stargazer/wasp/abi"

	"github.com/ethereum/go-ethereum"
	gethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Call struct {
	Name     string         `json:"name"`
	Target   common.Address `json:"target"`
	CallData []byte         `json:"call_data"`
}

type CallResponse struct {
	Success    bool   `json:"success"`
	ReturnData []byte `json:"returnData"`
}

func (call Call) GetMultiCall() abi.Multicall2Call {
	return abi.Multicall2Call{Target: call.Target, CallData: call.CallData}
}

func randomSigner() *bind.TransactOpts {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	signer, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1))
	if err != nil {
		panic(err)
	}

	signer.NoSend = true
	signer.Context = context.Background()
	signer.GasPrice = big.NewInt(0)

	return signer
}

type EthMultiCaller struct {
	Signer          *bind.TransactOpts
	Client          *ethclient.Client
	Abi             gethabi.ABI
	ContractAddress common.Address
}

func New(client *ethclient.Client) EthMultiCaller {
	// Load Multicall abi for later use
	mcAbi, err := gethabi.JSON(strings.NewReader(abi.MultiCall2ABI))
	if err != nil {
		panic(err)
	}

	contractAddress := common.HexToAddress(os.Getenv("MULTICALL_ADDRESS"))

	return EthMultiCaller{
		Signer:          randomSigner(),
		Client:          client,
		Abi:             mcAbi,
		ContractAddress: contractAddress,
	}
}

func (caller *EthMultiCaller) Execute(calls []Call) map[string]CallResponse {
	var responses []CallResponse

	var multiCalls = make([]abi.Multicall2Call, 0, len(calls))

	// Add calls to multicall structure for the contract
	for _, call := range calls {
		multiCalls = append(multiCalls, call.GetMultiCall())
	}

	// Prepare calldata for multicall
	callData, err := caller.Abi.Pack("tryAggregate", false, multiCalls)
	if err != nil {
		panic(err)
	}

	// Perform multicall
	resp, err := caller.Client.CallContract(context.Background(), ethereum.CallMsg{To: &caller.ContractAddress, Data: callData}, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Unpack results
	unpackedResp, _ := caller.Abi.Unpack("tryAggregate", resp)

	a, err := json.Marshal(unpackedResp[0])
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(a, &responses)
	if err != nil {
		panic(err)
	}

	// Create mapping for results. Be aware that we sometimes get two empty results initially, unsure why
	results := make(map[string]CallResponse)
	for i, response := range responses {
		results[calls[i].Name] = response
	}

	return results
}

var erc20Abi, _ = gethabi.JSON(strings.NewReader(abi.ERC20ABI))
var erc721Abi, _ = gethabi.JSON(strings.NewReader(abi.ERC721ABI))

func GetBalanceCall(name string, tokenAddress common.Address) Call {
	parsedData, err := erc20Abi.Pack("balanceOf", tokenAddress)
	if err != nil {
		panic(err)
	}
	return Call{
		Target:   tokenAddress,
		CallData: parsedData,
		Name:     name,
	}
}

func GetNameCall(name string, tokenAddress common.Address) Call {
	parsedData, err := erc20Abi.Pack("name")
	if err != nil {
		panic(err)
	}
	return Call{
		Target:   tokenAddress,
		CallData: parsedData,
		Name:     name,
	}
}

func GetSymbolCall(name string, tokenAddress common.Address) Call {
	parsedData, err := erc20Abi.Pack("symbol")
	if err != nil {
		panic(err)
	}
	return Call{
		Target:   tokenAddress,
		CallData: parsedData,
		Name:     name,
	}
}

func GetDecimalCall(name string, tokenAddress common.Address) Call {
	parsedData, err := erc20Abi.Pack("decimals")
	if err != nil {
		panic(err)
	}
	return Call{
		Target:   tokenAddress,
		CallData: parsedData,
		Name:     name,
	}
}

func GetTotalSupplyCall(name string, tokenAddress common.Address) Call {
	parsedData, err := erc20Abi.Pack("totalSupply")
	if err != nil {
		panic(err)
	}
	return Call{
		Target:   tokenAddress,
		CallData: parsedData,
		Name:     name,
	}
}

func GetAllowanceCall(name string, tokenAddress common.Address) Call {
	parsedData, err := erc20Abi.Pack("allowance", tokenAddress, tokenAddress)
	if err != nil {
		panic(err)
	}
	return Call{
		Target:   tokenAddress,
		CallData: parsedData,
		Name:     name,
	}
}

func (m *EthMultiCaller) IsContractErc20(contractAddress []byte) (*ErcInfo, bool) {
	contract := common.BytesToAddress(contractAddress)
	var userCalls = make([]Call, 0)
	userCalls = append(userCalls, GetBalanceCall("balance", contract))
	userCalls = append(userCalls, GetNameCall("name", contract))
	userCalls = append(userCalls, GetSymbolCall("symbol", contract))
	userCalls = append(userCalls, GetDecimalCall("decimal", contract))
	userCalls = append(userCalls, GetTotalSupplyCall("totalsupply", contract))
	userCalls = append(userCalls, GetAllowanceCall("allownace", contract))

	response := m.Execute(userCalls)
	for _, value := range response {
		if !value.Success {
			return &ErcInfo{}, false
		}
	}
	nameString := hex.EncodeToString(response["name"].ReturnData)

	// for _, c := range response["name"].ReturnData {
	// 	if c != 0 {
	// 		nameString += hex.EncodeToString(c[:])
	// 	}
	// }

	var symbolString string
	for _, c := range response["symbol"].ReturnData {
		if c != 0 {
			symbolString += string(c)
		}
	}

	ercInfo := &ErcInfo{
		Name:        nameString,
		Symbol:      symbolString,
		Decimal:     0,
		TotalSupply: "",
	}
	return ercInfo, true
}

func GetOwnerCall(name string, tokenAddress common.Address) Call {
	parsedData, err := erc721Abi.Pack("ownerOf", big.NewInt(0))
	if err != nil {
		panic(err)
	}
	return Call{
		Target:   tokenAddress,
		CallData: parsedData,
		Name:     name,
	}
}

func GetTokenUriCall(name string, tokenAddress common.Address) Call {
	parsedData, err := erc721Abi.Pack("tokenURI", big.NewInt(0))
	if err != nil {
		panic(err)
	}
	return Call{
		Target:   tokenAddress,
		CallData: parsedData,
		Name:     name,
	}
}

func GetApproved(name string, tokenAddress common.Address) Call {
	parsedData, err := erc721Abi.Pack("getApproved", big.NewInt(0))
	if err != nil {
		panic(err)
	}
	return Call{
		Target:   tokenAddress,
		CallData: parsedData,
		Name:     name,
	}
}

func GetSupportsInterface(name string, tokenAddress common.Address) Call {
	parsedData, err := erc721Abi.Pack("supportsInterface", [4]byte{0, 0, 0, 0})
	if err != nil {
		panic(err)
	}
	return Call{
		Target:   tokenAddress,
		CallData: parsedData,
		Name:     name,
	}
}

func (m *EthMultiCaller) IsContractERC721(contractAddress []byte) (*ErcInfo, bool) {
	contract := common.BytesToAddress(contractAddress)
	var userCalls = make([]Call, 0)
	userCalls = append(userCalls, GetBalanceCall("balance", contract))
	userCalls = append(userCalls, GetNameCall("name", contract))
	userCalls = append(userCalls, GetSymbolCall("symbol", contract))
	userCalls = append(userCalls, GetOwnerCall("owner", contract))
	userCalls = append(userCalls, GetTokenUriCall("owner", contract))
	userCalls = append(userCalls, GetSupportsInterface("supportsInterface", contract))

	response := m.Execute(userCalls)
	for _, value := range response {
		if !value.Success {
			return &ErcInfo{}, false
		}
	}

	var nameString string
	for _, c := range response["name"].ReturnData {
		if c != 0 {
			nameString += string(c)
		}
	}

	var symbolString string
	for _, c := range response["symbol"].ReturnData {
		if c != 0 {
			symbolString += string(c)
		}
	}
	ercInfo := &ErcInfo{
		Name:        nameString,
		Symbol:      symbolString,
		Decimal:     0,
		TotalSupply: "",
	}
	return ercInfo, true
}

type ErcInfo struct {
	Name        string
	Symbol      string
	Decimal     int64
	TotalSupply string
}
