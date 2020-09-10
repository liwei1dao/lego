package ethpay

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"testing"

	"github.com/liwei1dao/lego/sys/sdks/ethpay/pay"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

//创建钱包
func Test_CreateNewWallet(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // 0xfad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:]) // 0x049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address) // 0x96216849c49358B10257cb55b28eA603c874b05E

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:])) // 0x96216849c49358b10257cb55b28ea603c874b05e

	/*
		first
		6e92292424b5ef148eeeba7759cfdf9def532b2f2e8027b4c0cdddafac6a8d8e
		295fc4ab99f3427f74891c6c1d98050dd03e6f6c425dfea434a36908055b570ccad0ba57da16be678fd01a361d21b3a1f16a92b642b8ca36aeba704f9bce2911
		0xd8600d3C91c583A05047E761E4a86d224a5AC7ca
		0xd8600d3c91c583a05047e761e4a86d224a5ac7ca

		second
		d60e7a84978c08d2fa8f9ff53adb57898fc641e5c4186b2c807bc38d5340afb6
		aa2ed40791dde905685223d67a628a0242830571770ac45a6921ea892d28d85c65516f69962690227cfff25607fd547e2309f6fe69caa3316b213657204e9836
		0x5e62646910b49b4035a5FabE0E422c1A15aC0776
		0x5e62646910b49b4035a5fabe0e422c1a15ac0776
	*/
}

//查看钱包
func Test_LookWallet(t *testing.T) {
	client, err := ethclient.Dial("https://rinkeby.infura.io/v3/e48d4006b90e4b22b2ea6d4385e52ca9")
	if err != nil {
		log.Fatal(err)
	}

	account := common.HexToAddress("0xd84551C36a4e2e0F951a633852D557Be8CAaAf6a")
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance) // 25893180161173005034
}

//转账
func Test_Transfer(t *testing.T) {
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

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000)                // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0x5e62646910b49b4035a5FabE0E422c1A15aC0776")
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}

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
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	address, tx, instance, err := pay.DeployWallet(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(address.Hex())   // 0x147B8eb97fD247D06C4006D269c90C1908Fb5D54
	fmt.Println(tx.Hash().Hex()) // 0xdae8ba5444eefdc99f4d45cd0c4f24056cba6a02cefbf78066ef9f4188ff7dc0

	_ = instance

	/*  store
	0x70C217Eb968dC35025468c080E394142f79012E2
	0x8d08d1ad286a2977ed01d1a920bece7d46803c229187ab09f6c8200b10540cb1
	*/

	/* pay
	0x1620CaefD3cc4d0a33D1CAF2294e0EaC10D1Be81
	0x44624f593b5bbf05ae920710cb91dc27a8dd25646d506699f990a644dd8a5672
	*/

	/* account
	0xBaB9203562CB500484fb230303AcF73C99a7Dc8b
	0xa827c869882d0350f80cb6d826e860c3b9c2c348c09f04d3a6aaa33ab6678be4
	*/
}

//加载 智能合约
func Test_LoadSmartContract(t *testing.T) {
	client, err := ethclient.Dial("https://rinkeby.infura.io/v3/e48d4006b90e4b22b2ea6d4385e52ca9")
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0x70C217Eb968dC35025468c080E394142f79012E2")
	instance, err := store.NewStore(address, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("contract is loaded")
	_ = instance
}

//加载 智能合约 字节码 byte
func Test_LoadByteSmartContract(t *testing.T) {
	client, err := ethclient.Dial("https://rinkeby.infura.io/v3/e48d4006b90e4b22b2ea6d4385e52ca9")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0x70C217Eb968dC35025468c080E394142f79012E2")
	bytecode, err := client.CodeAt(context.Background(), contractAddress, nil) // nil is latest block
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hex.EncodeToString(bytecode)) // 60806...10029
}

//检查地址知否是合约
func Test_CheckAddrIsContract(t *testing.T) {
	client, err := ethclient.Dial("https://rinkeby.infura.io/v3/e48d4006b90e4b22b2ea6d4385e52ca9")
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0x7dfc481efab957b7abb7289692127ac737b10146")
	bytecode, err := client.CodeAt(context.Background(), address, nil) // nil is latest block
	if err != nil {
		log.Fatal(err)
	}

	isContract := len(bytecode) > 0

	fmt.Printf("is contract: %v\n", isContract) // is contract: true
}

//执行 智能合约 get
func Test_CallGet_SmartContract(t *testing.T) {
	client, err := ethclient.Dial("https://rinkeby.infura.io/v3/e48d4006b90e4b22b2ea6d4385e52ca9")
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0x70C217Eb968dC35025468c080E394142f79012E2")
	instance, err := store.NewStore(address, client)
	if err != nil {
		log.Fatal(err)
	}

	version, err := instance.Version(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(version) // "1.0"
}

//执行 智能合约 set
func Test_CallSet_SmartContract(t *testing.T) {
	client, err := ethclient.Dial("https://rinkeby.infura.io/v3/e48d4006b90e4b22b2ea6d4385e52ca9")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("d60e7a84978c08d2fa8f9ff53adb57898fc641e5c4186b2c807bc38d5340afb6")
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
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	address := common.HexToAddress("0x70C217Eb968dC35025468c080E394142f79012E2")
	instance, err := store.NewStore(address, client)
	if err != nil {
		log.Fatal(err)
	}

	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], []byte("foo"))
	copy(value[:], []byte("bar"))

	tx, err := instance.SetItem(auth, key, value)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s \n", tx.Hash().Hex()) // tx sent: 0x8d490e535678e9a24360e955d75b27ad307bdfb97a1dca51d0f3035dcee3e870

	result, err := instance.Items(nil, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(result[:])) // "bar"
}
