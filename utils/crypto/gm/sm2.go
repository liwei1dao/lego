package crypto

/*
国密算法库  github.com/tjfoc/gmsm/ 封装
*/

import (
	"bytes"
	"crypto/rand"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm3"
	"github.com/tjfoc/gmsm/sm4"
	"github.com/tjfoc/gmsm/x509"
)

///生成密钥对
func GenerateKey(key string) (*sm2.PrivateKey, error) {
	return sm2.GenerateKey(bytes.NewBufferString(key))
}

///读取私钥
func ReadPrivateKeyFromPem(privateKeyPem []byte, pwd []byte) (*sm2.PrivateKey, error) {
	return x509.ReadPrivateKeyFromPem(privateKeyPem, pwd)
}

///读取私钥
func ReadPublicKeyFromPem(privateKeyPem []byte) (*sm2.PublicKey, error) {
	return x509.ReadPublicKeyFromPem(privateKeyPem)
}

///国密—SM2-加密
func GM_SM2_Encry(origData []byte, pub *sm2.PublicKey) (ciphertext []byte, err error) {
	ciphertext, err = pub.EncryptAsn1(origData, rand.Reader)
	return
}

///国密—SM2-解密
func GM_SM2_Decry(ciphertext []byte, priv *sm2.PrivateKey) (origData []byte, err error) {
	origData, err = priv.DecryptAsn1(ciphertext)
	return
}

///国密—SM2-签名
func GM_SM2_Sign(origData, key string) (sign string, err error) {
	var (
		priv     *sm2.PrivateKey
		msg      []byte
		signdata []byte
	)
	if priv, err = sm2.GenerateKey(bytes.NewBufferString(key)); err != nil {
		return
	}
	msg = []byte(origData)
	if signdata, err = priv.Sign(rand.Reader, msg, nil); err != nil {
		return
	}
	sign = string(signdata)
	return
}

///国密—SM2-验签
func GM_SM2_Verify(origData, sign, key string) (isok bool, err error) {
	var (
		priv *sm2.PrivateKey
		pub  *sm2.PublicKey
	)
	if priv, err = sm2.GenerateKey(bytes.NewBufferString(key)); err != nil {
		return
	}
	pub = &priv.PublicKey
	isok = pub.Verify([]byte(origData), []byte(sign))
	return
}

///国密—SM3-哈希
func GM_SM3_Hash(origData string) (hash string) {
	h := sm3.New()
	h.Write([]byte(origData))
	hash = string(h.Sum(nil))
	return
}

///国密-SM4-sm4Ecb模式pksc7填充加密
func GM_SM4_Ecb(origData string, key string) (ciphertext string, err error) {
	var (
		iv      []byte
		ecbdata []byte
	)
	iv = []byte("0000000000000000")
	if err = sm4.SetIV(iv); err != nil { //设置SM4算法实现的IV值,不设置则使用默认值
		return
	}
	if ecbdata, err = sm4.Sm4Ecb([]byte(key), []byte(origData), true); err != nil { //sm4Ecb模式pksc7填充加密
		return
	}
	ciphertext = string(ecbdata)
	return
}

///国密-SM4-sm4Ecb模式pksc7填充加密
func GM_SM4_Dec(ciphertext string, key string) (origData string, err error) {
	var (
		iv      []byte
		ecbdata []byte
	)
	iv = []byte("0000000000000000")
	if err = sm4.SetIV(iv); err != nil { //设置SM4算法实现的IV值,不设置则使用默认值
		return
	}
	if ecbdata, err = sm4.Sm4Ecb([]byte(key), []byte(ciphertext), false); err != nil { //sm4Ecb模式pksc7填充加密
		return
	}
	origData = string(ecbdata)
	return
}
