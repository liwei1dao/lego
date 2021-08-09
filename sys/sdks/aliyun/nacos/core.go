package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type (
	ISys interface {
		RegisterInstance(service vo.RegisterInstanceParam) (success bool, err error)
		DeregisterInstance(service vo.DeregisterInstanceParam) (success bool, err error)
		GetService(param vo.GetServiceParam) (services model.Service, err error)
		SelectAllInstances(param vo.SelectAllInstancesParam) (instances []model.Instance, err error)
		SelectInstances(param vo.SelectInstancesParam) (instances []model.Instance, err error)
		SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (instance *model.Instance, err error)
		Subscribe(param *vo.SubscribeParam) (err error)
		Unsubscribe(param *vo.SubscribeParam) (err error)
		GetAllServicesInfo(param vo.GetAllServiceInfoParam) (serviceInfos model.ServiceList, err error)
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

///注册服务
func RegisterInstance(service vo.RegisterInstanceParam) (success bool, err error) {
	return defsys.RegisterInstance(service)
}

///注销服务
func DeregisterInstance(service vo.DeregisterInstanceParam) (success bool, err error) {
	return defsys.DeregisterInstance(service)
}

///获取服务
func GetService(param vo.GetServiceParam) (services model.Service, err error) {
	return defsys.GetService(param)
}

///获取全部服务实例
func SelectAllInstances(param vo.SelectAllInstancesParam) (instances []model.Instance, err error) {
	return defsys.SelectAllInstances(param)
}
func SelectInstances(param vo.SelectInstancesParam) (instances []model.Instance, err error) {
	return defsys.SelectInstances(param)
}
func SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (instance *model.Instance, err error) {
	return defsys.SelectOneHealthyInstance(param)
}
func Subscribe(param *vo.SubscribeParam) (err error) {
	return defsys.Subscribe(param)
}
func Unsubscribe(param *vo.SubscribeParam) (err error) {
	return defsys.Unsubscribe(param)
}
func GetAllServicesInfo(param vo.GetAllServiceInfoParam) (serviceInfos model.ServiceList, err error) {
	return defsys.GetAllServicesInfo(param)
}
