package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	gstorage "cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	fstorage "firebase.google.com/go/v4/storage"
	"google.golang.org/api/option"
)

func newSys(options Options) (sys *Storage, err error) {
	sys = &Storage{options: options}
	err = sys.init()
	return
}

type Storage struct {
	options Options
	app     *firebase.App
	client  *fstorage.Client
	bucket  *gstorage.BucketHandle
}

func (this *Storage) init() (err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(this.options.TimeOut))
	opt := option.WithCredentialsFile(this.options.ServiceAccountPath)
	config := &firebase.Config{
		ProjectID:     "resstation",
		StorageBucket: this.options.BucketName,
	}
	if this.app, err = firebase.NewApp(context.Background(), config, opt); err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}
	if this.client, err = this.app.Storage(ctx); err != nil {
		return fmt.Errorf("error Storage app: %v", err)
	}
	if this.bucket, err = this.client.DefaultBucket(); err != nil {
		return fmt.Errorf("error Bucket app: %v", err)
	}
	return
}

///上传文件
func (this *Storage) UploadFile(r io.Reader, object string) (err error) {
	var (
		wc  *gstorage.Writer
		ctx context.Context
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()
	wc = this.bucket.Object(object).NewWriter(ctx)
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
		rc  *gstorage.Reader
		ctx context.Context
	)
	ctx, _ = context.WithTimeout(context.Background(), time.Second*50)
	rc, err = this.bucket.Object(object).NewReader(ctx)
	if err != nil {
		err = fmt.Errorf("Object(%q).NewReader: %v", object, err)
		return
	}
	defer rc.Close()
	if _, err := io.Copy(w, rc); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	return nil

}
