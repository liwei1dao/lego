package registry

import (
	"fmt"
	"testing"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func Test_Sys_Nacose(t *testing.T) {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId: "cb351549-86d2-4b65-9416-9dbe24856bdb",
	}
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "172.20.27.145",
			ContextPath: "/nacos",
			Port:        10005,
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

	if succ, err := client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "127.0.0.2",
		Port:        8848,
		Weight:      1,
		GroupName:   "demo",
		ServiceName: "test_1",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata: map[string]string{
			"type":         "test",
			"category":     "test",
			"version":      fmt.Sprintf("%e", 1.0),
			"rpcid":        "rwercsfsdwer",
			"rpcsubscribe": "{}",
		},
	}); err != nil {
		fmt.Printf("RegisterInstance err:%v\n", err)
	} else {
		fmt.Printf("RegisterInstance succ:%v\n", succ)
	}

	if slist, err := client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		NameSpace: "cb351549-86d2-4b65-9416-9dbe24856bdb",
		GroupName: "demo",
		PageNo:    1,
		PageSize:  10,
	}); err != nil {
		fmt.Printf("GetAllServicesInfo err:%v\n", err)
	} else {
		fmt.Printf("GetAllServicesInfo :%+v\n", slist)
		for _, v := range slist.Doms {
			if services, err := client.SelectInstances(vo.SelectInstancesParam{
				GroupName:   "demo",
				ServiceName: v,
				HealthyOnly: true,
			}); err != nil {
				fmt.Printf("SelectInstances err:%v\n", err)
			} else {
				fmt.Printf("SelectInstances :%+v\n", services)
			}
			if services, err := client.SelectAllInstances(vo.SelectAllInstancesParam{
				GroupName:   "demo",
				ServiceName: v,
			}); err != nil {
				fmt.Printf("SelectAllInstances err:%v\n", err)
			} else {
				fmt.Printf("SelectAllInstances :%+v\n", services)
			}
		}
	}
}
