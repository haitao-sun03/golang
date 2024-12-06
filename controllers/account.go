package controllers

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"io/ioutil"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/haitao-sun03/go/config"
	"github.com/haitao-sun03/go/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
)

type AccountController struct{}

type CreateParam struct {
	Password string `json:"password"`
}

func (a AccountController) CreateAccount(ctx *gin.Context) {
	param := CreateParam{}
	ctx.Bind(&param)

	log.Info("password in :", param.Password)
	fmt.Println("keystore path", config.Config.Geth.KeystorePath)
	ks := keystore.NewKeyStore(config.Config.Geth.KeystorePath, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(param.Password)

	if err != nil {
		log.WithError(err).Error("create account error")
	}
	log.Info("new account :", account.Address.Hex())

	Success(ctx, http.StatusOK, "success", account.Address.Hex(), 0)
}

func (a AccountController) ImportAccount(ctx *gin.Context) {
	file := config.Config.Geth.KeystorePath + "/UTC--2024-11-28T11-35-15.079340300Z--6bd04fc15e031eead0e61cc78b6f7be6c179ffc0"
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	// password := "secret"
	accountImport, err := ks.Import(jsonBytes, "666666", "666666")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("import account :", accountImport.Address.Hex())
	if err := os.Remove(file); err != nil {
		log.Fatal(err)
	}
	Success(ctx, http.StatusOK, "success", accountImport.Address.Hex(), 0)
}

func (a AccountController) Foo(ctx *gin.Context) {

	account := common.HexToAddress("0x6c0db8c49190b517b949429b9dea1c2b32143bd2")
	blockNumber := big.NewInt(0)
	fmt.Println("===", blockNumber)
	balance, err := config.GethClient.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.WithError(err).Error("BalanceAt error")
	}
	fmt.Println(balance)

	pendingBalance, err := config.GethClient.PendingBalanceAt(context.Background(), account)
	if err != nil {
		log.WithError(err).Error("PendingBalanceAt error")
	}
	fmt.Println(pendingBalance) // 25729324269165216042
	Success(ctx, http.StatusOK, "success", "foo", 0)
}

func (a AccountController) Wallet(ctx *gin.Context) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.WithError(err).Error("GenerateKey error")
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes))
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // 0xfad9c8855b740a0b7e
	publicKey := privateKey.Public()
	fmt.Println("publicKey ", publicKey)
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fmt.Println("publicKeyECDSA ", publicKeyECDSA)

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("===", hexutil.Encode(publicKeyBytes))
	fmt.Println("===", hexutil.Encode(publicKeyBytes)[4:]) // 0x049a7df67f79246283f
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address) // 0x96216849c49358B10257cb55b28eA603c874b05E
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:]))

}

func (a AccountController) BlockHeaderAndBody(ctx *gin.Context) {

	header, err := config.GethClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(header.Number.String())
	fmt.Println(header.Coinbase.Hex())
	fmt.Println(header.Difficulty.Uint64())
	fmt.Println(header.Nonce)
	fmt.Println(header.Hash())
	fmt.Println(header.TxHash)
	fmt.Println(header.Root)
	fmt.Println(header.ReceiptHash)

	fmt.Println("---------------------")

	// blockNumber := big.NewInt(1)
	block, err := config.GethClient.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := config.GethClient.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("chainID : ", chainID)
	fmt.Println(block.Number().Uint64())     // 5671744
	fmt.Println(block.Time())                // 1527211625
	fmt.Println(block.Difficulty().Uint64()) // 3217000136609065
	fmt.Println(block.Hash().Hex())          // 0x9e8751ebb5069389b855bba72d949

	fmt.Println("通过block.Transactions()遍历txs---------------------")
	for _, tx := range block.Transactions() {
		fmt.Println("hash : ", tx.Hash())
		fmt.Println("value : ", tx.Value())    // 10000000000000000
		fmt.Println("tx nonce : ", tx.Nonce()) // 110644
		fmt.Println("data : ", tx.Data())
		fmt.Println("chainId : ", tx.ChainId())

		from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
		if err != nil {
			fmt.Printf("Error retrieving sender address: %v\n", err)
			continue
		}
		fmt.Println("from : ", from.Hex())         // 0x55fE59D8Ad77035154dDd0AD0388D09D
		fmt.Println("to : ", tx.To().Hex())        // 0x55fE59D8Ad77035154dDd0AD0388D09D
		fmt.Println("gas limit : ", tx.Gas())      // 105000
		fmt.Println("gas price : ", tx.GasPrice()) // 102000000000

		// block.Transactions() :have been packaged to block,all confirmed
		// but the method (TransactionByHash) can get tx through tx hash,the tx maybe pending or confirmed
		// _, isPending, err := config.GethClient.TransactionByHash(context.Background(), tx.Hash())

		// if err != nil {
		// 	log.Fatal(err)
		// }
		// if isPending {
		// 	log.Println("Transaction is still pending")
		// 	return
		// }

		// 获取交易收据
		receipt, err := config.GethClient.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("gas used : ", receipt.GasUsed)

		transactionFee := new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(receipt.GasUsed)))
		fmt.Println("gas fee : ", transactionFee)

		fmt.Println("---------------------")
	}

	fmt.Println("通过block hash以及tx数量遍历txs---------------------")

	// 获取区块中的交易数量
	count, err := config.GethClient.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}
	// 可以通过block.Hash()以及block中tx数量遍历到其中每个tx
	for idx := uint(0); idx < count; idx++ {
		tx, err := config.GethClient.TransactionInBlock(context.Background(), block.Hash(), idx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("hash : ", tx.Hash())
		fmt.Println("value : ", tx.Value())    // 10000000000000000
		fmt.Println("tx nonce : ", tx.Nonce()) // 110644
		fmt.Println("data : ", tx.Data())
		fmt.Println("chainId : ", tx.ChainId())

		from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
		if err != nil {
			fmt.Printf("Error retrieving sender address: %v\n", err)
			continue
		}
		fmt.Println("from : ", from.Hex())         // 0x55fE59D8Ad77035154dDd0AD0388D09D
		fmt.Println("to : ", tx.To().Hex())        // 0x55fE59D8Ad77035154dDd0AD0388D09D
		fmt.Println("gas limit : ", tx.Gas())      // 105000
		fmt.Println("gas price : ", tx.GasPrice()) // 102000000000

		// 获取交易收据
		receipt, err := config.GethClient.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("gas used : ", receipt.GasUsed)

		transactionFee := new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(receipt.GasUsed)))
		fmt.Println("gas fee : ", transactionFee)

		fmt.Println("---------------------")
	}
	// fmt.Println("---------------------")

	// count, err := config.GethClient.TransactionCount(context.Background(), block.Hash())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(count)
}

func suggestBaseFee() (*big.Int, error) {

	// 获取最新的区块头
	header, err := config.GethClient.HeaderByNumber(context.Background(), nil) // nil 表示最新区块
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve latest block header: %v", err)
	}

	// 从区块头中提取baseFee
	baseFee := header.BaseFee

	return baseFee, nil
}

func (a AccountController) TransferEther(ctx *gin.Context) {
	accountPath := fmt.Sprint(config.Config.Geth.KeystorePath, "/UTC--2024-11-28T03-26-22.871269100Z--6c0db8c49190b517b949429b9dea1c2b32143bd2")
	password := ""
	privateKey, err := utils.GetKeystorePK(accountPath, password)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Error("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	// 获取公钥地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("from address is :", fromAddress)
	// 获取该公钥交易的nonce
	nonce, err := config.GethClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.WithError(err).Error("get from account nonce fail")
	}
	fmt.Println("from account nonce is :", nonce)
	// 转账金额
	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	// 设置gas limit和gas price
	gasLimit := uint64(21000) // in units
	// 建议的矿工小费：用户可以选择支付一个额外的小费给矿工，以鼓励矿工更快地处理他们的交易。小费可以是0，但为了确保交易被优先处理，通常会设置一个较小的小费。
	gasTipCap, err := config.GethClient.SuggestGasTipCap(context.Background())
	if err != nil {
		log.WithError(err).Error("Failed to suggest gas tip cap")
	}
	// 建议的BaseFee：每个区块有一个动态调整的基础费用，根据网络的拥堵情况自动调整。基础费用被销毁，而不是支付给矿工。
	baseFee, err := suggestBaseFee()
	if err != nil {
		log.WithError(err).Error("Failed to get base fee")
	}

	gasFeeCap := new(big.Int).Add(gasTipCap, baseFee)
	// 接收地址对象
	toAddress := common.HexToAddress("0xea2194aeffcb3e2c192d5b8d7522b13bd6c7bac1")
	var data []byte
	chainID, err := config.GethClient.NetworkID(context.Background())
	if err != nil {
		log.WithError(err).Error("获取chainID失败")
	}
	fmt.Println("chainID is :", chainID)
	// 创建交易对象
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap, // maxPriorityFeePerGas：矿工小费上限
		GasFeeCap: gasFeeCap, // maxFeePerGas：交易费用上限
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      data,
	})
	// tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	// 使用from account 即发起交易者签名交易
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		log.WithError(err).Error("签名tx失败")
	}
	fmt.Println("signedTx is :", signedTx.Hash())

	// 发送交易
	err = config.GethClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.WithError(err).Error("发送交易失败")
	}

	// 查询交易状态
	_, isPending, err := config.GethClient.TransactionByHash(context.Background(), signedTx.Hash())
	if err != nil {
		log.WithError(err).Error("查询交易状态失败")
	}
	fmt.Println(isPending)
}

func (AccountController) TransferToken(ctx *gin.Context) {

	accountPath := fmt.Sprint(config.Config.Geth.KeystorePath, "/UTC--2024-11-28T03-26-22.871269100Z--6c0db8c49190b517b949429b9dea1c2b32143bd2")
	password := ""
	privateKey, err := utils.GetKeystorePK(accountPath, password)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Error("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	// 获取公钥地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("from address is :", fromAddress)
	// 获取该公钥交易的nonce
	nonce, err := config.GethClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.WithError(err).Error("get from account nonce fail")
	}
	fmt.Println("from account nonce is :", nonce)

	value := big.NewInt(0)
	toAddress := common.HexToAddress("0xcE8De9742BBA3a5D039Cf7516fCC6a9eC0839B6A")
	tokenAddress := common.HexToAddress("0xA03384C52E88b60cecc00EEF44f827caECe072b5")
	transferFnSignature := []byte("transfer(address,uint256)")

	// 使用crypto.Keccak256Hash来生成Keccak-256哈希
	hash := crypto.Keccak256Hash(transferFnSignature)
	// 提取前4个字节作为方法ID
	methodID := hash.Bytes()[:4]
	fmt.Println(hexutil.Encode(methodID)) //

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAddress))

	amount := new(big.Int)
	amount.SetString("100", 10)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	fmt.Printf("Padded amount: %x\n", paddedAmount)
	fmt.Println("hexutil.Encode: ", hexutil.Encode(paddedAmount))

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := config.GethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})

	if err != nil {
		log.WithError(err).Error("EstimateGas error")
	}
	gasLimit *= 2
	fmt.Println("gasLimit:", gasLimit)

	gasPrice, err := config.GethClient.SuggestGasPrice(context.Background())
	gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(120))
	gasPrice = new(big.Int).Div(gasPrice, big.NewInt(100))
	if err != nil {
		log.WithError(err).Error("SuggestGasPrice error")
	}
	fmt.Println("gasPrice:", gasPrice)

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	chainID, err := config.GethClient.NetworkID(context.Background())
	if err != nil {
		log.WithError(err).Error("chainID error")
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.WithError(err).Error("SignTx error")
	}

	fmt.Println("signedTx before is :", signedTx)
	fmt.Println("signedTx hash before is :", signedTx.Hash().Hex())
	// test *types.Transaction to hex string
	hexString := utils.Transection2HexStr(signedTx)
	fmt.Println("hex tx String is : ", hexString)

	// test hex string to *types.Transaction
	signedTxBack := utils.HexStr2Transection(hexString)
	fmt.Println("signedTx back is : ", signedTxBack)
	fmt.Println("signedTx hash back is :", signedTxBack.Hash().Hex())
	fmt.Println(signedTxBack == signedTx) //打印false

	err = config.GethClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.WithError(err).Error("SendTransaction error")
	}
	log.Printf("tx sent: %s", signedTx.Hash().Hex())

	// 查询交易状态
	_, isPending, err := config.GethClient.TransactionByHash(context.Background(), signedTx.Hash())
	if err != nil {
		log.WithError(err).Error("查询交易状态失败")
	}
	fmt.Println(isPending)
}

var (
	AccountPath string
	password    = ""
)

type MintIn struct {
	Amount int64 `json:"amount"`
}

func (AccountController) Mint(ctx *gin.Context) {
	mintIn := MintIn{}
	ctx.Bind(&mintIn)

	privateKey, err := utils.GetKeystorePK(AccountPath, password)
	if err != nil {
		log.WithError(err).Error("GetKeystorePK error")
		Fail(ctx, http.StatusInternalServerError, err)
		return
	}
	opts, shouldReturn := ConstructTransactionOpts(config.GethClient, privateKey, ctx, 300000, 1*1e9, 2*1e9)
	if shouldReturn {
		return
	}

	tx, err := config.ERC20Contract.Mint(opts, big.NewInt(mintIn.Amount))
	if err != nil {
		log.WithError(err).Error("Mint error")
		Fail(ctx, http.StatusInternalServerError, err)
		return
	}
	hashStr := tx.Hash().Hex()
	Success(ctx, http.StatusOK, "success", hashStr, 0)
}

type TransferToken struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

func (AccountController) TransferTokenWithABI(ctx *gin.Context) {
	transferToken := TransferToken{}
	ctx.Bind(&transferToken)
	privateKey, err := utils.GetKeystorePK(AccountPath, password)
	if err != nil {
		log.WithError(err).Error("GetKeystorePK error")
		Fail(ctx, http.StatusInternalServerError, err)
		return
	}
	// 设置EIP-1559相关的费用
	opts, shouldReturn := ConstructTransactionOpts(config.GethClient, privateKey, ctx, 3000000, 10*1e9, 20*1e9)
	if shouldReturn {
		return
	}
	tx, err := config.ERC20Contract.Transfer(opts, common.HexToAddress(transferToken.To), big.NewInt(int64(transferToken.Amount)))
	if err != nil {
		log.WithError(err).Error("Transfer error")
		Fail(ctx, http.StatusInternalServerError, err)
		return
	}
	hashStr := tx.Hash().Hex()
	Success(ctx, http.StatusOK, "success", hashStr, 0)
}

type BalanceOfIn struct {
	Account string `json:"account" binding:"required"`
}

func (AccountController) BalanceOf(ctx *gin.Context) {
	balanceIn := BalanceOfIn{}
	ctx.Bind(&balanceIn)
	fmt.Println("account is :", balanceIn.Account)
	balance, err := config.ERC20Contract.BalanceOf(nil, common.HexToAddress(balanceIn.Account))
	if err != nil {
		log.WithError(err).Error("BalanceOf error")
		Fail(ctx, http.StatusInternalServerError, err)
		return
	}

	Success(ctx, http.StatusOK, "success", balance, 0)
}
