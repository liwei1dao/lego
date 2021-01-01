package token

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/liwei1dao/lego/sys/sdks/eth/token/solidity"
)

//部署 智能合约
func Test_DeploySmartContract(t *testing.T) {
	client, err := ethclient.Dial("https://rinkeby.infura.io/v3/e48d4006b90e4b22b2ea6d4385e52ca9")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("6e92292424b5ef148eeeba7759cfdf9def532b2f2e8027b4c0cdddafac6a8d8e")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(4000000) // in units
	auth.GasPrice = gasPrice

	//HiToolCoin:代币名称 HTC:代币符号 18:代币小数点 big.NewInt(1):代币汇率 big.NewInt(1000000):代币发行总量 fromAddress:代币接受eth地址
	address, tx, instance, err := solidity.DeployHiToolCoin(auth, client, "HiToolCoin", "HTC", 18, big.NewInt(1), big.NewInt(1000000), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(address.Hex())   // 0x147B8eb97fD247D06C4006D269c90C1908Fb5D54
	fmt.Println(tx.Hash().Hex()) // 0xdae8ba5444eefdc99f4d45cd0c4f24056cba6a02cefbf78066ef9f4188ff7dc0

	_ = instance

	/*
		THC
		0x262C59905D6865FDbafA925C0Fae598Ab27327A7
		0x249c5e38d10455edaeba3b2b4ea38f4e2245dd4d6cc6f3fd9d7a09e74a8afe51
	*/
}
