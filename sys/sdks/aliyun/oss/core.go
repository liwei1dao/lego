package oss

import "github.com/liwei1dao/lego/core"

var (
	service core.IService
	defoss  IOSS
)

type (
	IOSS interface {
		CreateBucket(bucketName string) (err error)
		UploadFile(objectName string, localFileName string) (err error)
		DownloadFile(objectName string, downloadedFileName string) (err error)
		DeleteFile(objectName string) (err error)
	}
)

func OnInit(s core.IService, opt ...Option) (err error) {
	defoss, err = newOSS(opt...)
	return
}

func CreateBucket(bucketName string) (err error) {
	return defoss.CreateBucket(bucketName)
}

func UploadFile(localFileName string, objectName string) (err error) {
	return defoss.UploadFile(localFileName, objectName)
}

func DownloadFile(objectName string, downloadedFileName string) (err error) {
	return defoss.DownloadFile(objectName, downloadedFileName)
}

func DeleteFile(objectName string) (err error) {
	return defoss.DeleteFile(objectName)
}
