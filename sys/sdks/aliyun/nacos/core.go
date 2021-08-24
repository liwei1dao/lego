package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type (
	ISys interface {
		Naming_RegisterInstance(service vo.RegisterInstanceParam) (success bool, err error)
		Naming_DeregisterInstance(service vo.DeregisterInstanceParam) (success bool, err error)
		Naming_GetService(param vo.GetServiceParam) (services model.Service, err error)
		Naming_SelectAllInstances(param vo.SelectAllInstancesParam) (instances []model.Instance, err error)
		Naming_SelectInstances(param vo.SelectInstancesParam) (instances []model.Instance, err error)
		Naming_SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (instance *model.Instance, err error)
		Naming_Subscribe(param *vo.SubscribeParam) (err error)
		Naming_Unsubscribe(param *vo.SubscribeParam) (err error)
		Naming_GetAllServicesInfo(param vo.GetAllServiceInfoParam) (serviceInfos model.ServiceList, err error)
		Config_GetConfig(param vo.ConfigParam) (string, error)
		Config_PublishConfig(param vo.ConfigParam) (bool, error)
		Config_DeleteConfig(param vo.ConfigParam) (bool, error)
		Config_ListenConfig(params vo.ConfigParam) (err error)
		Config_CancelListenConfig(params vo.ConfigParam) (err error)
		Config_SearchConfig(param vo.SearchConfigParam) (*model.ConfigPage, error)
		Config_PublishAggr(param vo.ConfigParam) (published bool, err error)
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
func Naming_RegisterInstance(service vo.RegisterInstanceParam) (success bool, err error) {
	return defsys.Naming_RegisterInstance(service)
}

///注销服务
func Naming_DeregisterInstance(service vo.DeregisterInstanceParam) (success bool, err error) {
	return defsys.Naming_DeregisterInstance(service)
}

///获取服务
func Naming_GetService(param vo.GetServiceParam) (services model.Service, err error) {
	return defsys.Naming_GetService(param)
}

///获取全部服务实例
func Naming_SelectAllInstances(param vo.SelectAllInstancesParam) (instances []model.Instance, err error) {
	return defsys.Naming_SelectAllInstances(param)
}
func Naming_SelectInstances(param vo.SelectInstancesParam) (instances []model.Instance, err error) {
	return defsys.Naming_SelectInstances(param)
}
func Naming_SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (instance *model.Instance, err error) {
	return defsys.Naming_SelectOneHealthyInstance(param)
}
func Naming_Subscribe(param *vo.SubscribeParam) (err error) {
	return defsys.Naming_Subscribe(param)
}
func Naming_Unsubscribe(param *vo.SubscribeParam) (err error) {
	return defsys.Naming_Unsubscribe(param)
}
func Naming_GetAllServicesInfo(param vo.GetAllServiceInfoParam) (serviceInfos model.ServiceList, err error) {
	return defsys.Naming_GetAllServicesInfo(param)
}

func Config_GetConfig(param vo.ConfigParam) (string, error) {
	return defsys.Config_GetConfig(param)
}
func Config_PublishConfig(param vo.ConfigParam) (bool, error) {
	return defsys.Config_PublishConfig(param)
}
func Config_DeleteConfig(param vo.ConfigParam) (bool, error) {
	return defsys.Config_DeleteConfig(param)
}
func Config_ListenConfig(params vo.ConfigParam) (err error) {
	return defsys.Config_ListenConfig(params)
}
func Config_CancelListenConfig(params vo.ConfigParam) (err error) {
	return defsys.Config_ListenConfig(params)
}
func Config_SearchConfig(param vo.SearchConfigParam) (*model.ConfigPage, error) {
	return defsys.Config_SearchConfig(param)
}
func Config_PublishAggr(param vo.ConfigParam) (published bool, err error) {
	return defsys.Config_PublishAggr(param)
}
