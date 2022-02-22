package crypto

import (
	"fmt"
	"testing"
)

func Test_GM_SM2_Encry(t *testing.T) {
	origData, err := GM_SM2_Encry("token", "123456781234567812345678")
	fmt.Printf("origData:%s err:%v", origData, err)
}

func Test_GM_SM2_Decry(t *testing.T) {
	origData, err := GM_SM2_Decry("token", "123456781234567812345678")
	fmt.Printf("origData:%s err:%v", origData, err)
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
