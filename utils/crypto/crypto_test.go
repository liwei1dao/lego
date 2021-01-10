package crypto

import (
	"fmt"
	"testing"
)

func Test_CBC(t *testing.T) {
	encrypted := AesEncryptCBC([]byte("asdjoiqwjeio"), "123456781234567812345678")
	fmt.Printf("encrypted:%s", string(encrypted))
}
