package oss

import "testing"

func Test_OSSUploadFile(t *testing.T) {
	if err := OnInit(nil,
		SetEndpoint("http://oss-cn-hongkong.aliyuncs.com"),
		SetAccessKeyId("xxxxxxxxxxxxx"),                 //账号AccessKeyId
		SetAccessKeySecret("xxxxxxxxxxxxxxxxxxxxxxxxx"), //账号AccessKeySecret
		SetBucketName("xxxxxxx"),                        //存储空间
	); err != nil {
		t.Logf("初始化OSS 系统失败 err:%s", err.Error())
	} else {
		t.Logf("初始化OSS 成功")
	}
	// if err := CreateBucket("hitoolchat"); err != nil {
	// 	t.Logf("创建 CreateBucket  err:%s", err.Error())
	// } else {
	// 	t.Logf("创建 CreateBucket 成功")
	// }
	// if err := UploadFile("test/liwei1dao.jpg", "F:/liwei1dao.jpg"); err != nil {
	// 	t.Logf("上传OSS 系统失败 err:%s", err.Error())
	// } else {
	// 	t.Logf("上传OSS 成功")
	// }
	if err := DownloadFile("test/liwei1dao.jpg", "F:/liwei2dao.jpg"); err != nil {
		t.Logf("下载OSS 系统失败 err:%s", err.Error())
	} else {
		t.Logf("下载OSS 成功")
	}
}
