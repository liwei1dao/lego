package oss_test

import (
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/sys/sdks/aliyun/oss"
)

func Test_OSSUploadFile(t *testing.T) {
	sys, err := oss.NewSys(
		oss.SetEndpoint("http://oss-cn-shenzhen.aliyuncs.com"),
		oss.SetAccessKeyId("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
		oss.SetAccessKeySecret("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
		oss.SetSecurityToken("xxxxxxxxxxxxxxxxxxx"),
		oss.SetBucketName("liwei1dao"),
	)
	if err != nil {
		fmt.Printf("初始化OSS 系统失败 err:%v", err)
		return
	} else {
		fmt.Printf("初始化OSS 系统成功")
	}
	// if err := CreateBucket("hitoolchat"); err != nil {
	// 	t.Logf("创建 CreateBucket  err:%s", err.Error())
	// } else {
	// 	t.Logf("创建 CreateBucket 成功")
	// }
	if err := sys.UploadFile("test/tuoluo_icon.png", "./tuoluo_icon.png"); err != nil {
		fmt.Printf("上传OSS 系统失败 err:%s", err.Error())
	} else {
		fmt.Printf("上传OSS 成功")
	}
	// if file, err := GetObject("test/liwei1dao.jpg"); err != nil {
	// 	t.Logf("下载OSS 系统失败 err:%s", err.Error())
	// } else {
	// 	t.Logf("下载OSS 成功 len:%d", len(file))
	// }
}
