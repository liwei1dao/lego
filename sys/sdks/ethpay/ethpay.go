package ethpay

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/liwei1dao/lego/sys/sdks/ethpay/pay"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func newEthPay(opt ...Option) (IETHPay, error) {
	opts, err := newOptions(opt...)
	if err != nil {
		return nil, err
	}
	ethpay := &EthPay{
		opt: opts,
	}
	err = ethpay.Init()
	return ethpay, err
}

type EthPay struct {
	opt            *Options
	client         *ethclient.Client
	privateKey     *ecdsa.PrivateKey
	ControllerAddr common.Address //控制账号钱包地址
	WalletAdrr     common.Address //钱包合约地址
	RecieverAddr   common.Address //回收账号钱包地址
}

func (this *EthPay) Init() (err error) {
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
	this.ControllerAddr = crypto.PubkeyToAddress(*publicKeyECDSA)
	this.WalletAdrr = common.HexToAddress(this.opt.WalletAdrr)
	this.RecieverAddr = common.HexToAddress(this.opt.FundRecoveryAddr)
	return
}

//查看钱包余额
func (this *EthPay) LookBalance(addr string) uint64 {
	account := common.HexToAddress(addr)
	balance, err := this.client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return 0
	}
	return balance.Uint64()
}

//获取用户支付合约地址
func (this *EthPay) GetUserPayAddr(uhash string) (addr string, err error) {
	parsed, err := abi.JSON(strings.NewReader(pay.AccountABI))
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
	inithash := keccak256(mustHexDecode(pay.AccountBin), param)

	// Wallet 合约地址
	address := mustHexDecode(this.opt.WalletAdrr)
	salt := mustHexDecode(uhash)
	addr = "0x" + hex.EncodeToString(Keccak256([]byte{0xff}, address, salt, inithash)[12:])
	return
}

//部署支付账号合约
func (this *EthPay)DeployAccountContract(uhash string)(trans string,err error){
	nonce, err := this.client.PendingNonceAt(context.Background(), this.ControllerAddr)
	if err != nil {
		return "",err
	}
	gasPrice, err := this.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "",err
	}
	auth := bind.NewKeyedTransactor(this.privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	instance, err := pay.NewWallet(this.WalletAdrr, this.client)
	if err != nil {
		return "",err
	}
	
	salt := [32]byte{}
	copy(salt[:], MustHexDecode(uhash))

	tx, err := instance.Create(auth, this.RecieverAddr, salt)
	if err != nil {
		return "",err
	}
	
	trans = tx.Hash().Hex()
	return trans,nil
}

//回收用户支付合约下的金额
func (this *EthPay) RecycleUserMoney(uaddr string)(trans string,err error) {
	nonce, err := this.client.PendingNonceAt(context.Background(), this.ControllerAddr)
	if err != nil {
		return  "",err
	}s
	gasPrice, err := this.client.SuggestGasPrice(context.Background())
	if err != nil {
		return  "",err
	}

	auth := bind.NewKeyedTransactor(this.privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	address := common.HexToAddress(uaddr)
	instance, err := pay.NewAccount(address, this.client)
	if err != nil {
		return  "",err
	}

	tx, err := instance.Flush(auth)
	if err != nil {
		return  "",err
	}
	trans = tx.Hash().Hex()
	return trans,nil
}

//监听用户支付行为
func (this *EthPay)MonitorUserPay(uaddr string){

}