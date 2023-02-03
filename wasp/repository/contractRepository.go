package repository

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/berachain/stargazer/wasp/abi"
	"github.com/berachain/stargazer/wasp/database"
	"github.com/berachain/stargazer/wasp/models"
	"github.com/berachain/stargazer/wasp/queryClient"
	decoder "github.com/mingjingc/abi-decoder"
	"golang.org/x/crypto/sha3"
)

type ContractRepo struct {
	db *database.Database
	qc *queryClient.QueryClient
}

func NewContractRepo(db *database.Database, qc *queryClient.QueryClient) *ContractRepo {
	return &ContractRepo{
		db: db,
		qc: qc,
	}
}

func BuildContract(txnModel *models.TransactionModel) *models.Contract {

	calculatedContractAddress := calculateContractAddress(txnModel.From, txnModel.Nonce)
	fmt.Printf("from: %s", hex.EncodeToString(txnModel.From))
	calculatedAddressAsHex := hex.EncodeToString(calculatedContractAddress)
	fmt.Println("CALCULATED ADDRESS")
	fmt.Println(calculatedAddressAsHex)
	//find ABIs
	return &models.Contract{
		Address:       calculatedContractAddress,
		DeployTxnHash: txnModel.Hash,
		Creator:       txnModel.From,
	}
}

func calculateContractAddress(address []byte, nonce uint64) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(address)
	hash.Write([]byte{byte(nonce)})
	calculatedAddressAsBytes := hash.Sum(nil)
	return calculatedAddressAsBytes[12:]
}

func TestDefaultABIs(txnData []byte) {
	txDataHex := hex.EncodeToString(txnData)

	txDataDecoder := decoder.NewABIDecoder()
	txDataDecoder.SetABI(abi.ERC20ABI)
	method, err := txDataDecoder.DecodeMethod(txDataHex)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(method.Name)
}
