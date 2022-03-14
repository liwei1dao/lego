package lib

import (
	"github.com/liwei1dao/lego/core"
)

const (
	SC_ServiceGateRouteComp core.S_Comps = "SC_GateRouteComp" //s_comps.ISC_GateRouteComp
)

const (
	SM_GateModule    core.M_Modules = "SM_GateModule"    //gate模块 	网关服务模块
	SM_HttpModule    core.M_Modules = "SM_HttpModule"    //hhtp模块		web服务模块
	SM_MonitorModule core.M_Modules = "SM_MonitorModule" //monitor模块	监控模块
	SM_ConsoleModule core.M_Modules = "SM_ConsoleModule" //cnsole模块	控制台模块
	SM_LiveModule    core.M_Modules = "SM_LiveModule"    //live模块		流媒体模块
	SM_CollyModule   core.M_Modules = "SM_CollyModule"   //colly模块	爬虫模块
)
