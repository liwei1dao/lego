package influxdb

import (
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

type (
	IInfluxdb interface {
		Setup(username, password, org, bucket string, timeout int) (*domain.OnboardingResponse, error)
		QueryAPI(org string) api.QueryAPI
		WriteAPI(org, bucket string) api.WriteAPI
		WriteAPIBlocking(org, bucket string) api.WriteAPIBlocking
		Close()
	}
)

var (
	defsys IInfluxdb
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys IInfluxdb, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}
func Setup(username, password, org, bucket string, timeout int) (*domain.OnboardingResponse, error) {
	return defsys.Setup(username, password, org, bucket, timeout)
}
func QueryAPI(org string) api.QueryAPI {
	return defsys.QueryAPI(org)
}
func WriteAPI(org, bucket string) api.WriteAPI {
	return defsys.WriteAPI(org, bucket)
}
func WriteAPIBlocking(org, bucket string) api.WriteAPIBlocking {
	return defsys.WriteAPIBlocking(org, bucket)
}
func Close() {
	defsys.Close()
}
