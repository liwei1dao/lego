package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

/*
CBC对称加密 注意 key 的长度必须固定 16, 24, 32
*/

// =================== CBC ======================
func AesEncryptCBC(_origData string, key string) (ciphertext string) {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher([]byte(key))
	blockSize := block.BlockSize()                         // 获取秘钥块的长度
	origData := pkcs5Padding([]byte(_origData), blockSize) // 补全码
	iv := []byte("0000000000000000")
	blockMode := cipher.NewCBCEncrypter(block, iv) // 加密模式
	encrypted := make([]byte, len(origData))       // 创建数组
	blockMode.CryptBlocks(encrypted, origData)     // 加密
	ciphertext = base64.StdEncoding.EncodeToString(encrypted)
	return ciphertext
}
func AesDecryptCBC(ciphertext string, key string) (origData string) {
	encrypted, _ := base64.StdEncoding.DecodeString(ciphertext)
	block, _ := aes.NewCipher([]byte(key)) // 分组秘钥
	iv := []byte("0000000000000000")
	blockMode := cipher.NewCBCDecrypter(block, iv) // 加密模式
	decrypted := make([]byte, len(encrypted))      // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)    // 解密
	decrypted = pkcs5UnPadding(decrypted)          // 去除补全码
	origData = string(decrypted)
	return origData
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
