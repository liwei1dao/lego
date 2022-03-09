package aes

import (
	"fmt"
	"testing"
)

func Test_CBC(t *testing.T) {
	token := AesEncryptCBC("asdjoiqwjeio", "123456781234567812345678")
	fmt.Printf("encrypted:%s", token)

	origData := AesDecryptCBC(token, "123456781234567812345678")
	fmt.Printf("encrypted:%s", string(origData))
}
