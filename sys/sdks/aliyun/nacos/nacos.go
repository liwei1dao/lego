package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func newSys(options Options) (sys *Nacos, err error) {
	sys = &Nacos{
		options: options,
	}
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         options.NamespaceId,
		TimeoutMs:           options.TimeoutMs,
		NotLoadCacheAtStart: true,
		LogDir:              "./nacos/log",
		CacheDir:            "./nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      options.NacosAddr,
			ContextPath: "/nacos",
			Port:        options.Port,
			Scheme:      "http",
		},
	}
	if sys.options.NacosClientType == NamingClient || sys.options.NacosClientType == All {
		// 创建服务发现客户端的另一种方式 (推荐)
		sys.namingClient, err = clients.NewNamingClient(
			vo.NacosClientParam{
				ClientConfig:  &clientConfig,
				ServerConfigs: serverConfigs,
			},
		)
	}
	if sys.options.NacosClientType == ConfigClient || sys.options.NacosClientType == All {
		// 创建服务发现客户端的另一种方式 (推荐)
		sys.configClient, err = clients.NewConfigClient(
			vo.NacosClientParam{
				ClientConfig:  &clientConfig,
				ServerConfigs: serverConfigs,
			},
		)
	}
	return
}

type Nacos struct {
	options      Options
	namingClient naming_client.INamingClient
	configClient config_client.IConfigClient
}

///注册服务
func (this *Nacos) Naming_RegisterInstance(service vo.RegisterInstanceParam) (success bool, err error) {
	success, err = this.namingClient.RegisterInstance(service)
	return
}

///注销服务
func (this *Nacos) Naming_DeregisterInstance(service vo.DeregisterInstanceParam) (success bool, err error) {
	success, err = this.namingClient.DeregisterInstance(service)
	return
}

///获取服务
func (this *Nacos) Naming_GetService(param vo.GetServiceParam) (services model.Service, err error) {
	services, err = this.namingClient.GetService(param)
	return
}

/// SelectAllInstance可以返回全部实例列表,包括healthy=false,enable=false,weight<=0
func (this *Nacos) Naming_SelectAllInstances(param vo.SelectAllInstancesParam) (instances []model.Instance, err error) {
	instances, err = this.namingClient.SelectAllInstances(param)
	return
}

///选中实例列表
func (this *Nacos) Naming_SelectInstances(param vo.SelectInstancesParam) (instances []model.Instance, err error) {
	instances, err = this.namingClient.SelectInstances(param)
	return
}

/// SelectOneHealthyInstance将会按加权随机轮询的负载均衡策略返回一个健康的实例 实例必须满足的条件：health=true,enable=true and weight>0
func (this *Nacos) Naming_SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (instance *model.Instance, err error) {
	instance, err = this.namingClient.SelectOneHealthyInstance(param)
	return
}

/// Subscribe key=serviceName+groupName+cluster 注意:我们可以在相同的key添加多个SubscribeCallback.
func (this *Nacos) Naming_Subscribe(param *vo.SubscribeParam) (err error) {
	err = this.namingClient.Subscribe(param)
	return
}

func (this *Nacos) Naming_Unsubscribe(param *vo.SubscribeParam) (err error) {
	err = this.namingClient.Unsubscribe(param)
	return
}

func (this *Nacos) Naming_GetAllServicesInfo(param vo.GetAllServiceInfoParam) (serviceInfos model.ServiceList, err error) {
	serviceInfos, err = this.namingClient.GetAllServicesInfo(param)
	return
}

func (this *Nacos) Config_GetConfig(param vo.ConfigParam) (string, error) {
	return this.configClient.GetConfig(param)
}

func (this *Nacos) Config_PublishConfig(param vo.ConfigParam) (bool, error) {
	return this.configClient.PublishConfig(param)
}

func (this *Nacos) Config_DeleteConfig(param vo.ConfigParam) (bool, error) {
	return this.configClient.PublishConfig(param)
}

func (this *Nacos) Config_ListenConfig(params vo.ConfigParam) (err error) {
	return this.configClient.ListenConfig(params)
}

func (this *Nacos) Config_CancelListenConfig(params vo.ConfigParam) (err error) {
	return this.configClient.CancelListenConfig(params)
}

func (this *Nacos) Config_SearchConfig(param vo.SearchConfigParam) (*model.ConfigPage, error) {
	return this.configClient.SearchConfig(param)
}

func (this *Nacos) Config_PublishAggr(param vo.ConfigParam) (published bool, err error) {
	return this.configClient.PublishAggr(param)
}
