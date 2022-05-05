package jwt

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/liwei1dao/lego/sys/gin/engine"
)

func NewJWT(key string) *JWT {
	return &JWT{jwtkey: key}
}

type JWT struct {
	jwtkey string
}

// CreateToken 生成token
func (this *JWT) CreateToken(Id string) (string, error) {
	expireTime := time.Now().Add(2 * time.Hour) //过期时间
	nowTime := time.Now()                       //当前时间
	claims := jwt.StandardClaims{
		Id:        Id,                //用户Id
		ExpiresAt: expireTime.Unix(), //过期时间戳
		IssuedAt:  nowTime.Unix(),    //当前时间戳
		Issuer:    "blogLeo",         //颁发者签名
		Subject:   "userToken",       //签名主题

	}
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenStruct.SignedString([]byte(this.jwtkey))
}

// CheckToken 验证token
func (this *JWT) CheckToken(token string) (*jwt.StandardClaims, bool) {
	tokenObj, _ := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return this.jwtkey, nil
	})
	if key, _ := tokenObj.Claims.(*jwt.StandardClaims); tokenObj.Valid {
		return key, true
	} else {
		return nil, false
	}
}

// JwtMiddleware jwt中间件
func (this *JWT) JwtMiddleware() engine.HandlerFunc {
	return func(c *engine.Context) {
		//从请求头中获取token
		tokenStr := c.Request.Header.Get("Authorization")
		//用户不存在
		if tokenStr == "" {
			c.JSON(http.StatusOK, engine.H{"code": 0, "msg": "用户不存在"})
			c.Abort() //阻止执行
			return
		}
		//token格式错误
		tokenSlice := strings.SplitN(tokenStr, " ", 2)
		if len(tokenSlice) != 2 && tokenSlice[0] != "Bearer" {
			c.JSON(http.StatusOK, engine.H{"code": 0, "msg": "token格式错误"})
			c.Abort() //阻止执行
			return
		}
		//验证token
		tokenStruck, ok := this.CheckToken(tokenSlice[1])
		if !ok {
			c.JSON(http.StatusOK, engine.H{"code": 0, "msg": "token不正确"})
			c.Abort() //阻止执行
			return
		}
		//token超时
		if time.Now().Unix() > tokenStruck.ExpiresAt {
			c.JSON(http.StatusOK, engine.H{"code": 0, "msg": "token过期"})
			c.Abort() //阻止执行
			return
		}
		c.Set("UserId", tokenStruck.Id)
		c.Next()
	}
}
