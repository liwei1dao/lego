package sts_test

import (
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/sys/sdks/aliyun/sts"
)

func Test_STS(t *testing.T) {
	sys, err := sts.NewSys(
		sts.SetRegionId("cn-shenzhen"),
		sts.SetAccessKeyId("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
		sts.SetAccessKeySecret("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
	)
	if err != nil {
		fmt.Printf("初始化OSS 系统失败 err:%v", err)
		return
	} else {
		fmt.Printf("初始化OSS 系统成功")
		auth, err := sys.AssumeRole("xxxxxxxxxxxxxxxxxxxxxxxxxxx", "SessionTest")
		fmt.Printf("初始化OSS AssumeRole auth:%+v err:%v", auth, err)
	}
}
