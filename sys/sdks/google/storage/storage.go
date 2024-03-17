package storage

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func newSys(options Options) (sys *Storage, err error) {
	sys = &Storage{options: options}
	err = sys.init()
	return
}

type Storage struct {
	options Options
	client  *storage.Client
	bucket  *storage.BucketHandle
}

func (this *Storage) init() (err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(this.options.TimeOut))
	jsonKey, err := ioutil.ReadFile(this.options.ServiceAccountPath)
	if err != nil {
		return err
	}
	conf, err := google.JWTConfigFromJSON(jsonKey, storage.ScopeReadWrite)
	if err != nil {
		return err
	}
	if this.client, err = storage.NewClient(ctx, option.WithTokenSource(conf.TokenSource(ctx))); err != nil {
		return
	}
	this.bucket = this.client.Bucket(this.options.BucketName)
	return
}

///上传文件
func (this *Storage) UploadFile(r io.Reader, object string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := this.bucket.Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, r); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}

//下载文件
func (this *Storage) DownloadFile(w io.Writer, object string) (err error) {
	var (
		rc *storage.Reader
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	rc, err = this.bucket.Object(object).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	if _, err = io.Copy(w, rc); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	return nil

}

func (this *Storage) Close() (err error) {
	err = this.client.Close()
	return
}
