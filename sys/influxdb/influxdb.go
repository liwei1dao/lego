package influxdb

import (
	"context"
	"time"

	iclient "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

func newSys(options Options) (sys *Influxdb, err error) {
	sys = &Influxdb{options: options}
	err = sys.init()
	return
}

type Influxdb struct {
	options Options
	client  iclient.Client
}

func (this *Influxdb) init() (err error) {
	this.client = iclient.NewClient(this.options.Addr, this.options.Token)
	_, err = this.Ping(this.getContext())
	return
}

func (this *Influxdb) getContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), time.Duration(this.options.TimeOut)*time.Second)
	return
}
func (this *Influxdb) Setup(username, password, org, bucket string, timeout int) (*domain.OnboardingResponse, error) {
	return this.client.Setup(context.Background(), username, password, org, bucket, timeout)
}

func (this *Influxdb) Ping(ctx context.Context) (bool, error) {
	return this.client.Ping(ctx)
}

func (this *Influxdb) WriteAPI(org, bucket string) api.WriteAPI {
	return this.client.WriteAPI(org, bucket)
}
func (this *Influxdb) WriteAPIBlocking(org, bucket string) api.WriteAPIBlocking {
	return this.client.WriteAPIBlocking(org, bucket)
}
func (this *Influxdb) QueryAPI(org string) api.QueryAPI {
	return this.client.QueryAPI(org)
}

func (this *Influxdb) Close() {
	this.client.Close()
}
