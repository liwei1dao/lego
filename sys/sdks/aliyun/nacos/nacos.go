package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
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
	// 创建服务发现客户端的另一种方式 (推荐)
	sys.client, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	return
}

type Nacos struct {
	options Options
	client  naming_client.INamingClient
}

///注册服务
func (this *Nacos) RegisterInstance(service vo.RegisterInstanceParam) (success bool, err error) {
	success, err = this.client.RegisterInstance(service)
	return
}

///注销服务
func (this *Nacos) DeregisterInstance(service vo.DeregisterInstanceParam) (success bool, err error) {
	success, err = this.client.DeregisterInstance(service)
	return
}

///获取服务
func (this *Nacos) GetService(param vo.GetServiceParam) (services model.Service, err error) {
	services, err = this.client.GetService(param)
	return
}

/// SelectAllInstance可以返回全部实例列表,包括healthy=false,enable=false,weight<=0
func (this *Nacos) SelectAllInstances(param vo.SelectAllInstancesParam) (instances []model.Instance, err error) {
	instances, err = this.client.SelectAllInstances(param)
	return
}

///选中实例列表
func (this *Nacos) SelectInstances(param vo.SelectInstancesParam) (instances []model.Instance, err error) {
	instances, err = this.client.SelectInstances(param)
	return
}

/// SelectOneHealthyInstance将会按加权随机轮询的负载均衡策略返回一个健康的实例 实例必须满足的条件：health=true,enable=true and weight>0
func (this *Nacos) SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (instance *model.Instance, err error) {
	instance, err = this.client.SelectOneHealthyInstance(param)
	return
}

/// Subscribe key=serviceName+groupName+cluster 注意:我们可以在相同的key添加多个SubscribeCallback.
func (this *Nacos) Subscribe(param *vo.SubscribeParam) (err error) {
	err = this.client.Subscribe(param)
	return
}

func (this *Nacos) Unsubscribe(param *vo.SubscribeParam) (err error) {
	err = this.client.Unsubscribe(param)
	return
}

func (this *Nacos) GetAllServicesInfo(param vo.GetAllServiceInfoParam) (serviceInfos model.ServiceList, err error) {
	serviceInfos, err = this.client.GetAllServicesInfo(param)
	return
}
