package repository

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/berachain/stargazer/lib/common"
	"github.com/berachain/stargazer/lib/crypto"
	"github.com/berachain/stargazer/wasp/abi"
	"github.com/berachain/stargazer/wasp/database"
	"github.com/berachain/stargazer/wasp/models"
	"github.com/berachain/stargazer/wasp/multicall"
	"github.com/berachain/stargazer/wasp/query"
	"github.com/berachain/stargazer/wasp/queryClient"
	gethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"golang.org/x/crypto/sha3"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type TransactionRepo struct {
	db *database.Database
	qc *queryClient.QueryClient
}

func NewTransactionRepo(db *database.Database, qc *queryClient.QueryClient) *TransactionRepo {
	return &TransactionRepo{
		db: db,
		qc: qc,
	}
}

func (r *TransactionRepo) BuildTransactionList(block *types.Block) *[]models.TransactionModel {
	txns := []models.TransactionModel{}
	chainID, err := r.qc.NetworkID(context.Background())
	if err != nil {
		panic("unable to retrieve chainId")
	}
	signerType := types.LatestSignerForChainID(chainID)
	for _, t := range block.Transactions() {
		txn := *r.BuildTransaction(
			t,
			block.Number().String(),
			block.Time(),
			block.BaseFee(),
			signerType)
		txns = append(txns, txn)
	}
	return &txns
}
func (r *TransactionRepo) BuildTransaction(
	txn *types.Transaction,
	blockNumber string,
	time uint64,
	baseFee *big.Int,
	signerType types.Signer) *models.TransactionModel {

	receipt := r.BuildTransactionReceipt(txn, blockNumber, time, baseFee)
	txnModel := models.GethToTransactionModel(
		txn,
		blockNumber,
		time,
		baseFee,
		signerType,
		receipt)

	if txnModel.To == nil {
		//contract creation
		fmt.Println("CONTRACT CREATION")
		fmt.Println(hex.EncodeToString(txnModel.Hash))

		c := r.BuildContract(txnModel)
		res := r.CreateContract(c)
		if res == 1 {
			panic("failed to create contract")
		}
		fmt.Println(c)
		return txnModel
	}

	// detect transfer even
	return txnModel
}

func (r *TransactionRepo) BuildTransactionReceipt(
	txn *types.Transaction,
	blockNumber string,
	time uint64,
	baseFee *big.Int) models.EthTxnReceipt {
	gethReciept, err := r.qc.GetTransactionReceiptByHash(txn.Hash())
	if err != nil {
		panic(err)
	}
	ethTxnLogs := r.BuildTransactionLogs(gethReciept)
	txnReceiptModel := models.GethToReceiptModel(gethReciept, ethTxnLogs)
	return *txnReceiptModel
}

func (r *TransactionRepo) BuildTransactionLogs(receipt *types.Receipt) []models.EthLog {
	logs := []models.EthLog{}
	for _, log := range receipt.Logs {
		logModel := *models.GethToEthLogModel(log)
		r.ParsedLog(log)
		logs = append(logs, logModel)
	}
	return logs
}

func (r *TransactionRepo) TestDefaultABIs(txnData []byte, contractAddress []byte) (*multicall.ErcInfo, int64) {
	multicall := multicall.New(r.qc.GetEthClient())
	info, isErc20 := multicall.IsContractErc20(contractAddress)
	if isErc20 {
		return info, 1
	}
	info, isErc721 := multicall.IsContractERC721(contractAddress)
	if isErc721 {
		return info, 2
	}
	fmt.Println(info)

	return nil, 0
}

func (r *TransactionRepo) BuildContract(txnModel *models.TransactionModel) *models.Contract {

	calculatedContractAddress := r.calculateContractAddress(txnModel.From, txnModel.Nonce)
	_, abiId := r.TestDefaultABIs(txnModel.Data, calculatedContractAddress)
	fmt.Println(abiId)
	//find ABIs
	return &models.Contract{
		Address:       calculatedContractAddress,
		DeployTxnHash: txnModel.Hash,
		Creator:       txnModel.From,
		AbiId:         abiId,
		// Name:          "",
		// Symbol:        "",
		// Decimal:       0,
		// TotalSupply:   "",
	}
}

func (r *TransactionRepo) CreateContract(contract *models.Contract) int {
	data, err := json.Marshal(contract)
	if err != nil {
		panic(err)
	}

	req := &database.SetRequest{
		RedisDb: contract.GetRedisDb(),
		Key:     contract.GetRedisKey(),
		Value:   data,
	}

	err = r.db.Set(req, func() error {
		res := r.db.Gorm.Create(&contract)
		return res.Error
	})

	r.CreateContractBalance(contract)

	if err != nil {
		return 1
	}

	return 0
}

func (r *TransactionRepo) CreateContractBalance(contract *models.Contract) {
	ctx := context.Background()
	balance, err := r.qc.GetErc20Balance(ctx, contract.Address, contract.Address)
	if err != nil {
		fmt.Println(err)
	}
	err = r.HandleNewErc20Balance(ctx, contract.Address, contract.Address, balance)
	if err != nil {
		panic(err)
	}
}
func (r *TransactionRepo) calculateContractAddress(address []byte, nonce uint64) []byte {
	hash := sha3.NewLegacyKeccak256()

	inputRLP, err := rlp.EncodeToBytes([]interface{}{address, nonce})
	if err != nil {
		fmt.Println(err)
	}

	hash.Write(inputRLP)
	calculatedAddressAsBytes := hash.Sum(nil)
	return calculatedAddressAsBytes[12:]
}

var erc20Abi, _ = gethabi.JSON(strings.NewReader(abi.ERC20ABI))
var dead = common.HexToAddress("0x0000000000000000000000000000000000000000")

var test = common.HexToAddress("0xb089e5fd920e0cae7fa25ef3da177920da820655077399023ee4c48d951e0c1d")

func (r *TransactionRepo) ParsedLog(ethlog *types.Log) (Transfer, error) {
	fmt.Println(ethlog.TxHash.Hex())
	if ethlog.TxHash.Hex() == test.Hex() {
		fmt.Print("REEE")
	}
	logTransferSig := []byte("Transfer(address,address,uint256)")
	LogApprovalSig := []byte("Approval(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	logApprovalSigHash := crypto.Keccak256Hash(LogApprovalSig)
	switch ethlog.Topics[0].Hex() {
	case logTransferSigHash.Hex():
		var transferEvent Transfer
		if len(ethlog.Topics) == 4 {
			transferEvent.Type = 0
			transferEvent.From = common.HexToAddress(ethlog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(ethlog.Topics[2].Hex())
			transferEvent.Contract = ethlog.Address
			transferEvent.Tokens = new(big.Int).SetBytes(ethlog.Topics[3].Bytes())
			r.UpdateBalance(&transferEvent, ethlog.BlockNumber)
			return transferEvent, nil
		}
		val, _ := erc20Abi.Unpack("Transfer", ethlog.Data)

		value, ok := val[0].(*big.Int)
		if !ok {
			log.Fatal("cannot unmarshal token value")
		}

		transferEvent.Type = 1
		transferEvent.From = common.HexToAddress(ethlog.Topics[1].Hex())
		transferEvent.To = common.HexToAddress(ethlog.Topics[2].Hex())
		transferEvent.Contract = ethlog.Address
		transferEvent.Tokens = value
		r.UpdateBalance(&transferEvent, ethlog.BlockNumber)
		return transferEvent, nil

	case logApprovalSigHash.Hex():

		var approvalEvent LogApproval

		// _, err := erc20Abi.Unpack("Approval", ethlog.Data)

		approvalEvent.TokenOwner = common.HexToAddress(ethlog.Topics[1].Hex())
		approvalEvent.Spender = common.HexToAddress(ethlog.Topics[2].Hex())
	}
	return Transfer{}, nil
}

func (r *TransactionRepo) UpdateBalance(transfer *Transfer, blockNumber uint64) {
	ctx := context.Background()
	// NFT CASE
	if transfer.Type == 0 {
		// MINT CASE
		if transfer.From.Hex() == dead.Hex() {
			err := r.HandleNewAccount(ctx, transfer.To.Bytes(), blockNumber)
			if err != nil {
				fmt.Println(err)
			}
			receiverErc721BalanceModel, err := r.HandleNewErc721Balance(ctx, transfer.To.Bytes(), transfer.Contract.Bytes())
			if err != nil {
				log.Fatal(err)
			}
			err = r.CreateErc721Token(ctx, receiverErc721BalanceModel.ID, transfer.Tokens.Int64())
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// TRANSFER CASE
			err := r.HandleNewAccount(ctx, transfer.From.Bytes(), blockNumber)
			if err != nil {
				fmt.Println(err)
			}
			err = r.HandleNewAccount(ctx, transfer.To.Bytes(), blockNumber)
			if err != nil {
				fmt.Println(err)
			}

			// DELETE NFT TOKEN ID FROM SENDER
			senderErc721BalanceModel, err := r.HandleNewErc721Balance(ctx, transfer.From.Bytes(), transfer.Contract.Bytes())
			if err != nil {
				log.Fatal(err)
			}
			err = r.DeleteErc721Token(ctx, senderErc721BalanceModel.ID, transfer.Tokens.Int64())
			if err != nil {
				log.Fatal(err)
			}

			receiverErc721BalanceModel, err := r.HandleNewErc721Balance(ctx, transfer.To.Bytes(), transfer.Contract.Bytes())
			if err != nil {
				log.Fatal(err)
			}
			err = r.CreateErc721Token(ctx, receiverErc721BalanceModel.ID, transfer.Tokens.Int64())
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		// ERC20 TRANSFER CASE
		err := r.HandleNewAccount(ctx, transfer.From.Bytes(), blockNumber)
		if err != nil {
			fmt.Println(err)
		}
		err = r.HandleNewAccount(ctx, transfer.To.Bytes(), blockNumber)
		if err != nil {
			fmt.Println(err)
		}
		// subtract balance from sender
		senderUpdatedAmount := big.NewInt(0).Sub(big.NewInt(0), transfer.Tokens)
		err = r.HandleNewErc20Balance(ctx, transfer.From.Bytes(), transfer.Contract.Bytes(), senderUpdatedAmount)
		if err != nil {
			log.Fatal(err)
		}
		// add balance to receiver
		receiverUpdatedAmount := big.NewInt(0).Add(big.NewInt(0), transfer.Tokens)
		err = r.HandleNewErc20Balance(ctx, transfer.To.Bytes(), transfer.Contract.Bytes(), receiverUpdatedAmount)
		if err != nil {
			log.Fatal(err)
		}

	}
}

func (r *TransactionRepo) HandleNewErc20Balance(ctx context.Context, address []byte, contractAddress []byte, amount *big.Int) error {
	db := query.Use(r.db.Gorm)
	erc20client := db.Erc20Balance
	queryDao := erc20client.WithContext(ctx).Where(erc20client.Address.Eq(address), erc20client.ContractAddress.Eq(contractAddress))
	erc20BalanceModel, err := queryDao.First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		erc20BalanceModel := &models.Erc20Balance{
			Address:         address,
			ContractAddress: contractAddress,
			Amount:          amount.String(),
		}
		res := r.db.Gorm.Create(&erc20BalanceModel)
		if res.Error != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}
	// if exists, update record with new balance
	originalAmount, ok := new(big.Int).SetString(erc20BalanceModel.Amount, 10)
	if !ok {
		fmt.Println("SetString: error")
	}
	updatedBalance := big.NewInt(0).Add(originalAmount, amount)
	_, err = queryDao.Update(erc20client.Amount, updatedBalance.String())
	if err != nil {
		return err
	}
	return nil
}

func (r *TransactionRepo) HandleNewErc721Balance(ctx context.Context, address []byte, contractAddress []byte) (*models.Erc721Balance, error) {
	db := query.Use(r.db.Gorm)
	erc721client := db.Erc721Balance
	erc721BalanceModel, err := erc721client.WithContext(ctx).Where(erc721client.Address.Eq(address), erc721client.ContractAddress.Eq(contractAddress)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		erc721BalanceModel := &models.Erc721Balance{
			Address:         address,
			ContractAddress: contractAddress,
		}
		res := r.db.Gorm.Create(&erc721BalanceModel)
		if res.Error != nil {
			return &models.Erc721Balance{}, err
		}
		return erc721BalanceModel, nil
	} else if err != nil {
		return &models.Erc721Balance{}, err
	}
	return erc721BalanceModel, nil
}

func (r *TransactionRepo) CreateErc721Token(ctx context.Context, id int64, tokenId int64) error {
	db := query.Use(r.db.Gorm)
	erc721Tokenclient := db.Erc721Tokens
	erc721Token := models.Erc721Tokens{
		BalanceId: id,
		TokenId:   tokenId,
	}
	err := erc721Tokenclient.WithContext(ctx).Create(&erc721Token)
	if err != nil {
		return err
	}
	return nil
}

func (r *TransactionRepo) DeleteErc721Token(ctx context.Context, id int64, tokenId int64) error {
	db := query.Use(r.db.Gorm)
	erc721Tokenclient := db.Erc721Tokens
	_, err := erc721Tokenclient.WithContext(ctx).Where(erc721Tokenclient.BalanceId.Eq(id), erc721Tokenclient.TokenId.Eq(tokenId)).Delete()
	if err != nil {
		return err
	}
	return nil
}

func (r *TransactionRepo) HandleNewAccount(ctx context.Context, address []byte, blockNumber uint64) error {
	_, error := r.GetEthAccount(ctx, address)
	if errors.Is(error, gorm.ErrRecordNotFound) {
		b, err := r.qc.GetEthBalance(ctx, address, blockNumber)
		if err != nil {
			return err
		}
		ethAccount := &models.EthAccount{
			Address:    address,
			Alias:      "",
			Balance:    b.String(),
			IsContract: false,
		}
		r.db.Gorm.Create(&ethAccount)
	}
	return nil
}

func (r *TransactionRepo) GetEthAccount(ctx context.Context, address []byte) (*models.EthAccount, error) {
	db := query.Use(r.db.Gorm)
	ethAccount := db.EthAccount
	ethAccountModel, err := ethAccount.WithContext(ctx).Where(ethAccount.Address.Eq(address)).First()
	if err != nil {
		return &models.EthAccount{}, err
	}
	return ethAccountModel, nil
}

type LogTransfer struct {
	From     common.Address
	To       common.Address
	Contract string

	Tokens *big.Int
}

type Transfer struct {
	Type     int
	From     common.Address
	To       common.Address
	Contract common.Address
	Tokens   *big.Int
}
type LogApproval struct {
	TokenOwner common.Address
	Spender    common.Address
	Tokens     *big.Int
}
