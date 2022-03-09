package md5

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

/*
MD5加密
*/

//MD5加密 大写
func MD5EncToUpper(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

//MD5加密 小写
func MD5EncToLower(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return strings.ToLower(hex.EncodeToString(h.Sum(nil)))
}

//参数签名
func ParamSign(param map[string]interface{}, key string) (orsign, sign string) {
	a := sort.StringSlice{}
	for k, _ := range param {
		a = append(a, k)
	}
	sort.Sort(a)
	orsign = ""
	for _, k := range a {
		switch param[k].(type) { //只签名基础数据
		case bool, byte, int8, int16, uint16, int32, uint32, int, int64, uint64, float64, float32, string:
			orsign = orsign + fmt.Sprintf("%s=%v&", k, param[k])
			break
		}
	}
	orsign = orsign + fmt.Sprintf("key=%s", key)
	sign = MD5EncToLower(orsign)
	return
}
