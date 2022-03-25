package influxdb

import (
	"time"

	iclient "github.com/influxdata/influxdb1-client/v2"
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
	this.client, err = iclient.NewHTTPClient(iclient.HTTPConfig{
		Addr:     this.options.Addr,
		Username: this.options.Username,
		Password: this.options.Password,
	})
	return
}

func (this *Influxdb) Ping(timeout time.Duration) (time.Duration, string, error) {
	return this.client.Ping(timeout)
}

func (this *Influxdb) Write(bp iclient.BatchPoints) error {
	return this.client.Write(bp)
}

func (this *Influxdb) Query(q iclient.Query) (*iclient.Response, error) {
	return this.client.Query(q)
}

func (this *Influxdb) QueryAsChunk(q iclient.Query) (*iclient.ChunkedResponse, error) {
	return this.client.QueryAsChunk(q)
}

func (this *Influxdb) Close() error {
	return this.client.Close()
}

func (this *Influxdb) NewPoint(name string,
	tags map[string]string,
	fields map[string]interface{},
	t ...time.Time) (*iclient.Point, error) {
	return iclient.NewPoint(name, tags, fields, t...)
}

func (this *Influxdb) NewBatchPoints(conf iclient.BatchPointsConfig) (iclient.BatchPoints, error) {
	return iclient.NewBatchPoints(conf)
}
