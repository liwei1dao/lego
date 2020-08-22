package monitor

import (
	"github.com/liwei1dao/lego/core"
)

const ( //Rpc
	Rpc_GetServiceMonitorInfo core.Rpc_Key = "GetServiceMonitorInfo" //获取服务监控信息

	Rpc_SetMonitorServiceSetting core.Rpc_Key = "SetMonitorServiceSetting" //设置服务配置信息
	Rpc_SetMonitorModuleSetting  core.Rpc_Key = "SetMonitorModuleSetting"  //设置模块配置信息
)
