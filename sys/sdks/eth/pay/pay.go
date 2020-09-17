package pay

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/liwei1dao/lego/sys/sdks/eth/pay/solidity"
)

func newPay(opt ...Option) (IPay, error) {
	opts, err := newOptions(opt...)
	if err != nil {
		return nil, err
	}
	Pay := &Pay{
		opt: opts,
	}
	err = Pay.Init()
	return Pay, err
}

type Pay struct {
	opt            *Options
	client         *ethclient.Client
	privateKey     *ecdsa.PrivateKey
	WalletAdrr     common.Address //钱包合约地址
	walletInstance *solidity.Wallet
	ControllerAddr common.Address //控制账号钱包地址
	RecieverAddr   common.Address //回收账号钱包地址
}

func (this *Pay) Init() (err error) {
	this.client, err = ethclient.Dial(this.opt.EthPoolAdrr)
	if err != nil {
		return
	}

	this.privateKey, err = crypto.HexToECDSA(this.opt.ControllerPrivateKey)
	if err != nil {
		return
	}

	publicKey := this.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	this.WalletAdrr = common.HexToAddress(this.opt.WalletAdrr)
	this.walletInstance, err = solidity.NewWallet(this.WalletAdrr, this.client)
	if err != nil {
		return
	}

	this.ControllerAddr = crypto.PubkeyToAddress(*publicKeyECDSA)
	this.RecieverAddr = common.HexToAddress(this.opt.FundRecoveryAddr)
	return
}

//查看钱包余额
func (this *Pay) LookBalance(addr string) uint64 {
	account := common.HexToAddress(addr)
	balance, err := this.client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return 0
	}
	return balance.Uint64()
}

//获取用户支付合约地址
func (this *Pay) GetUserPayAddr(uhash string) (addr string, err error) {
	parsed, err := abi.JSON(strings.NewReader(solidity.AccountABI))
	if err != nil {
		return "", err
	}
	// Account 合约构造函数设置了 reciever 参数
	// 为了以后能生成这个地址，这个需要持久化保存
	param, err := parsed.Pack("", this.RecieverAddr)
	if err != nil {
		return "", err
	}
	// 计算 Account 合约初始化哈希
	// 360c3c0304ab4f09eee311be7433387a83c3d62c7150e7654dfa339f5294eb45
	inithash := keccak256(mustHexDecode(solidity.AccountBin), param)
	// Wallet 合约地址
	address := mustHexDecode(this.opt.WalletAdrr)
	salt := mustHexDecode(uhash)
	addr = "0x" + hex.EncodeToString(keccak256([]byte{0xff}, address, salt, inithash)[12:])
	return
}

//部署支付账号合约
func (this *Pay) DeployAccountContract(uhash string) (trans string, err error) {
	nonce, err := this.client.PendingNonceAt(context.Background(), this.ControllerAddr)
	if err != nil {
		return "", err
	}
	gasPrice, err := this.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	auth := bind.NewKeyedTransactor(this.privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	salt := [32]byte{}
	copy(salt[:], mustHexDecode(uhash))

	tx, err := this.walletInstance.Create(auth, this.RecieverAddr, salt)
	if err != nil {
		return "", err
	}

	trans = tx.Hash().Hex()
	return trans, nil
}

//回收用户支付合约下的金额
func (this *Pay) RecycleUserMoney(uaddr string) (trans string, err error) {
	nonce, err := this.client.PendingNonceAt(context.Background(), this.ControllerAddr)
	if err != nil {
		return "", err
	}
	gasPrice, err := this.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	auth := bind.NewKeyedTransactor(this.privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	address := common.HexToAddress(uaddr)
	instance, err := solidity.NewAccount(address, this.client)
	if err != nil {
		return "", err
	}

	tx, err := instance.Flush(auth)
	if err != nil {
		return "", err
	}
	trans = tx.Hash().Hex()
	return trans, nil
}

//监听用户支付行为
func (this *Pay) MonitorUserPay(uaddr string, timeout time.Time) (value uint64, err error) {
	contractAbi, err := abi.JSON(strings.NewReader(string(solidity.AccountABI)))
	if err != nil {
		return 0, err
	}

	contractAddress := common.HexToAddress(uaddr)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := this.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return 0, err
	}

	ctx, cancel := context.WithDeadline(context.Background(), timeout)

	for {
		select {
		case err := <-sub.Err():
			return 0, err
		case <-ctx.Done():
			cancel()
			return 0, fmt.Errorf("Time Out")
		case vLog := <-logs:
			// fmt.Println(vLog.BlockHash.Hex()) // 0x3404b8c050aa0aacd0223e91b5c32fee6400f357764771d0684fa7b3f448f1a8
			// fmt.Println(vLog.BlockNumber)     // 2394201
			// fmt.Println(vLog.TxHash.Hex())    // 0x280201eda63c9ff6f305fcee51d5eb86167fab40ca3108ec784e8652a0e2b1a6
			event := struct {
				Vaule *big.Int
			}{}
			err := contractAbi.Unpack(&event, "Recharge", vLog.Data)
			if err != nil {
				return 0, err
			}
			// fmt.Println(event.Vaule)
			return event.Vaule.Uint64(), nil
		}
	}
}
