package lgid

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

/*
系统描述:ID系统 提供多进程级别的唯一id生成器，支持绑定tag标签，可以用于id反推数据
*/
type (
	ISys interface {
		//唯一的数值id
		NewNumberID() int64
		//唯一的字符串id
		NewStringID() string
		NewStringIDForNumber() string
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

// 唯一的数值id
func NewNumberID() int64 {
	return defsys.NewNumberID()
}

// 唯一的字符串id
func NewStringID() string {
	return defsys.NewStringID()
}

// 获取字符串的数值型id
func NewStringIDForNumber() string {
	return defsys.NewStringIDForNumber()
}

// tada 会话id
// 生成指定长度的随机唯一十六进制字符串
func NewTadaSessionId() string {
	length := 38
	// 获取当前时间戳（纳秒）
	timestamp := time.Now().UnixNano()
	// 剩余需要的随机字节数
	randomLength := (length - 16) / 2 // 16 个字符用于时间戳（8 字节）
	// 随机字节
	randomBytes := make([]byte, randomLength)
	if _, err := rand.Read(randomBytes); err != nil {
		fmt.Println("Failed to generate random bytes:", err)
		return ""
	}
	// 时间戳部分（转为 16 进制）
	timestampHex := fmt.Sprintf("%016x", timestamp)
	// 随机部分
	randomHex := hex.EncodeToString(randomBytes)
	// 拼接时间戳和随机数，确保唯一性
	return timestampHex + randomHex
}
