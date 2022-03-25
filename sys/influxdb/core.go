package influxdb

type (
	IInfluxdb interface {
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
