package zookeeper

type (
	ISys interface {
		Start() error
		RegisterFunction(serviceName, fname string, fn interface{}, meta string) error
		Stop() error
	}
)

func NewSys(option ...Option) (sys ISys, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}
