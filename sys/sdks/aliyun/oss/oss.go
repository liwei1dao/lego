package oss

import (
	"bytes"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func newSys(options Options) (sys *OSS, err error) {
	sys = &OSS{options: options}
	err = sys.Init()
	return
}

type OSS struct {
	options Options
	client  *oss.Client
	bucket  *oss.Bucket
}

func (this *OSS) Init() (err error) {
	if this.options.SecurityToken != "" {
		this.client, err = oss.New(this.options.Endpoint, this.options.AccessKeyId, this.options.AccessKeySecret, oss.SecurityToken(this.options.SecurityToken))
	} else {
		this.client, err = oss.New(this.options.Endpoint, this.options.AccessKeyId, this.options.AccessKeySecret)
	}
	if err != nil {
		return err
	}
	if ok, err := this.client.IsBucketExist(this.options.BucketName); !ok || err != nil {
		if err = this.CreateBucket(this.options.BucketName); err == nil {
			this.bucket, err = this.client.Bucket(this.options.BucketName)
			return err
		} else {
			return err
		}
	} else {
		this.bucket, err = this.client.Bucket(this.options.BucketName)
		return err
	}
}

//创建存储空间。
func (this *OSS) CreateBucket(bucketName string) (err error) {
	err = this.client.CreateBucket(bucketName)
	return err
}

//上传文件
//<objectName>上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
//<localFileName>由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt。
func (this *OSS) UploadFile(objectName string, localFileName string) (err error) {
	err = this.bucket.PutObjectFromFile(objectName, localFileName)
	return err
}

//上传对象
// <objectName>上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
// <localFileName>由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt。
// 上传文件。
func (this *OSS) UploadObject(objectKey string, reader io.Reader, options ...oss.Option) (err error) {
	err = this.bucket.PutObject(objectKey, reader, options...)
	return err
}

//下载文件。
//<objectName>从OSS下载文件时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
func (this *OSS) DownloadFile(objectName string, downloadedFileName string) (err error) {
	err = this.bucket.GetObjectToFile(objectName, downloadedFileName)
	return err
}

//下载文件到缓存
func (this *OSS) GetObject(objectName string, options ...oss.Option) ([]byte, error) {
	if file, err := this.bucket.GetObject(objectName, options...); err != nil {
		return nil, err
	} else {
		defer file.Close()
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(file); err != nil {
			return nil, err
		} else {
			return buf.Bytes(), nil
		}
	}
}

//删除文件
func (this *OSS) DeleteFile(objectName string) (err error) {
	err = this.bucket.DeleteObject(objectName)
	return
}

//获取临时访问url
func (this *OSS) GetURL(objectName string, expired int64, options ...oss.Option) (url string, err error) {
	url, err = this.bucket.SignURL(objectName, oss.HTTPGet, expired, options...)
	return
}
