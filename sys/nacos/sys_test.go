package nacos_test

import (
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/sys/nacos"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func Test_SysInit(t *testing.T) {
	sys, err := nacos.NewSys(
		nacos.SetNacosAddr("127.0.0.1"),
		nacos.SetPort(8848),
		nacos.SetNamespaceId("17d02ef9-afaa-4878-ad6c-9fd697f3b628"),
	)
	if err != nil {
		fmt.Printf("启动系统错误err:%v", err)
	}
	succ, err := sys.Naming_RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        8848,
		Weight:      0.69,
		GroupName:   "demo",
		ClusterName: "demo",
		ServiceName: "demo1",
		Metadata: map[string]string{
			"tag": "demo",
		},
	})
	fmt.Printf("sys RegisterInstance succ:%v err:%v", succ, err)
	// time.Sleep(time.Second * 3)
	serviceInfos, err := sys.Naming_GetAllServicesInfo(vo.GetAllServiceInfoParam{
		NameSpace: "17d02ef9-afaa-4878-ad6c-9fd697f3b628",
		GroupName: "demo",
	})

	fmt.Printf("sys GetAllServicesInfo serviceInfos:%+v err:%v", serviceInfos, err)

	// succ, err = sys.DeregisterInstance()
	// fmt.Printf("sys DeregisterInstance succ:%v err:%v", succ, err)
}
