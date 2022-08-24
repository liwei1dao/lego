package jwt

import (
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt"
)

func Test_ParamSign(t *testing.T) {
	if token, err := CreateToken("mysces", "liwei1dao"); err != nil {
		fmt.Printf("err:%v", err)
		return
	} else {
		tokenObj, _ := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("mysces"), nil
		})
		if key, ok := tokenObj.Claims.(*jwt.StandardClaims); ok && tokenObj.Valid {
			fmt.Printf("key:%v", key)
		} else {
			fmt.Printf("void")
		}
	}
}
