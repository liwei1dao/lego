package base64

import "encoding/base64"

///编码
func EncodeToString(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

///解码
func DecodeString(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
