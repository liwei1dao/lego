package crypto

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

func Test_GM_SM2_Encry(t *testing.T) {
	priv, err := sm2.GenerateKey(nil) // 生成密钥对
	if err != nil {
		t.Fatal(err)
	}
	privPem, err := x509.WritePrivateKeyToPem(priv, nil) // 生成密钥文件
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("privPem:%s\n", string(privPem))
	}
	pubKey, _ := priv.Public().(*sm2.PublicKey)
	pubkeyPem, err := x509.WritePublicKeyToPem(pubKey) // 生成公钥文件
	_, err = x509.ReadPrivateKeyFromPem(privPem, nil)  // 读取密钥
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("pubkeyPem:%s\n", string(pubkeyPem))
	}
	pubKey, err = x509.ReadPublicKeyFromPem(pubkeyPem) // 读取公钥
	if err != nil {
		t.Fatal(err)
	}
	cipher, err := pubKey.EncryptAsn1([]byte(
		`liwei1dao
		liwei2dao
		liwei3dao
		liwei4dao`), rand.Reader)
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("cipher:%s\n", string(cipher))
	}

	origData, err := priv.DecryptAsn1(cipher)
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("origData:%s\n", string(origData))
	}
}

func Test_GM_SM2_Decry(t *testing.T) {

}

func Test_GM_SM2_Sign(t *testing.T) {
	origData, err := GM_SM2_Sign("token", "123456781234567812345678")
	fmt.Printf("origData:%s err:%v", origData, err)
}

func Test_GM_SM2_Verify(t *testing.T) {
	isok, err := GM_SM2_Verify("token", "123456781234567812345678", "")
	fmt.Printf("isok:%v rr:%v", isok, err)
}

func Test_GM_SM3_Hash(t *testing.T) {
	hash := GM_SM3_Hash("123456781234567812345678")
	fmt.Printf("hash:%v", hash)
}

func Test_GM_SM4_Ecb(t *testing.T) {
	ciphertext, err := GM_SM4_Ecb("123456781234567812345678", "")
	fmt.Printf("ciphertext:%v err:%v", ciphertext, err)
}

func Test_GM_SM4_Dec(t *testing.T) {
	ciphertext, err := GM_SM4_Dec("123456781234567812345678", "")
	fmt.Printf("ciphertext:%v err:%v", ciphertext, err)
}

func Test_GM_SM2_Hex(t *testing.T) {
	// if priv, err := x509.ReadPrivateKeyFromHex(""); err != nil {

	// }
}
