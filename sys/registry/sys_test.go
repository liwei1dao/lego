package registry_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func Test_Sys_Nacose(t *testing.T) {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId: "ac1b23d5-1c14-4485-9e08-f3cde8c83163",
	}
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "172.20.27.126",
			ContextPath: "/nacos",
			Port:        8888,
			Scheme:      "http",
		},
	}

	// 创建服务发现客户端的另一种方式 (推荐)
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		fmt.Printf("初始化系统失败\n")
	}

	// if succ, err := client.RegisterInstance(vo.RegisterInstanceParam{
	// 	Ip:          "127.0.0.2",
	// 	Port:        8848,
	// 	Weight:      1,
	// 	ServiceName: "test_1",
	// 	Enable:      true,
	// 	Healthy:     true,
	// 	Ephemeral:   true,
	// 	Metadata: map[string]string{
	// 		"type":         "test",
	// 		"category":     "test",
	// 		"version":      fmt.Sprintf("%e", 1.0),
	// 		"rpcid":        "rwercsfsdwer",
	// 		"rpcsubscribe": "{}",
	// 	},
	// }); err != nil {
	// 	fmt.Printf("RegisterInstance err:%v\n", err)
	// } else {
	// 	fmt.Printf("RegisterInstance succ:%v\n", succ)
	// }

	if slist, err := client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		NameSpace: "ac1b23d5-1c14-4485-9e08-f3cde8c83163",
		GroupName: "datacollector",
		PageNo:    1,
		PageSize:  20,
	}); err != nil {
		fmt.Printf("GetAllServicesInfo err:%v\n", err)
	} else {
		fmt.Printf("GetAllServicesInfo :%+v\n", slist)
		for _, v := range slist.Doms {
			if instances, err := client.SelectInstances(vo.SelectInstancesParam{
				ServiceName: v,
				GroupName:   "datacollector",
				HealthyOnly: true,
			}); err == nil {
				fmt.Printf("instances :%+v\n", instances)
			} else {
				fmt.Printf("instances err:%v\n", err)
			}
			// if services, err := client.SelectInstances(vo.SelectInstancesParam{
			// 	ServiceName: v,
			// 	HealthyOnly: true,
			// }); err != nil {
			// 	fmt.Printf("SelectInstances err:%v\n", err)
			// } else {
			// 	fmt.Printf("SelectInstances :%+v\n", services)
			// }
			// if services, err := client.SelectAllInstances(vo.SelectAllInstancesParam{
			// 	ServiceName: v,
			// }); err != nil {
			// 	fmt.Printf("SelectAllInstances err:%v\n", err)
			// } else {
			// 	fmt.Printf("SelectAllInstances :%+v\n", services)
			// }
		}
	}
}

func Test_Sys_Consul(t *testing.T) {
	config := api.DefaultConfig()
	config.Address = "172.20.27.145:10003"
	if client, err := api.NewClient(config); err != nil {
		fmt.Printf("NewClient Err:%v\n", err)
		return
	} else {
		if err := client.Agent().ServiceRegister(&api.AgentServiceRegistration{
			ID:   "test",
			Name: "test",
			Tags: []string{"test"},
			Meta: map[string]string{},
		}); err != nil {
			fmt.Printf("ServiceRegister Err:%v\n", err)
			return
		} else {
			fmt.Printf("ServiceRegister Succ\n")
		}
	}
}
