package monitor

import (
	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib"
)

func NewModule() core.IServiceMonitor {
	m := new(Monitor)
	return m
}

type Monitor struct {
	cbase.ModuleBase
	service            base.IClusterService
	ServiceMonitor     *core.ServiceMonitor
	ServiceSetting     map[string]func(newvalue string) (err error)
	ModulesSetting     map[core.M_Modules]map[string]func(newvalue string) (err error)
	ServiceMonitorComp *ServiceMonitorComp
}

func (this *Monitor) GetType() core.M_Modules {
	return lib.SM_MonitorModule
}

func (this *Monitor) Init(service core.IService, module core.IModule, setting map[string]interface{}) (err error) {
	this.service = service.(base.IClusterService)
	this.ServiceMonitor = &core.ServiceMonitor{
		ServiceId:       this.service.GetId(),
		ServiceType:     this.service.GetType(),
		ServiceCategory: this.service.GetCategory(),
		ServiceVersion:  this.service.GetVersion(),
		ServiceTag:      this.service.GetTag(),
		Setting:         make(map[string]*core.SettingItem),
		ModuleMonitor:   make(map[core.M_Modules]*core.ModuleMonitor),
	}
	this.ServiceSetting = make(map[string]func(newvalue string) (err error))
	this.ModulesSetting = make(map[core.M_Modules]map[string]func(newvalue string) (err error))
	err = this.ModuleBase.Init(service, module, setting)
	for k, v := range this.service.GetSettings().Settings {
		this.ServiceMonitor.Setting[k] = &core.SettingItem{
			ItemName: k,
			IsWrite:  false,
			Data:     v,
		}
	}
	for k, v := range this.service.GetSettings().Modules {
		this.ServiceMonitor.ModuleMonitor[core.M_Modules(k)] = &core.ModuleMonitor{
			ModuleName: core.M_Modules(k),
			Setting:    make(map[string]*core.SettingItem),
		}
		this.ModulesSetting[core.M_Modules(k)] = make(map[string]func(newvalue string) (err error))
		for k1, v1 := range v {
			this.ServiceMonitor.ModuleMonitor[core.M_Modules(k)].Setting[k1] = &core.SettingItem{
				ItemName: k1,
				IsWrite:  false,
				Data:     v1,
			}
		}
	}

	return
}

func (this *Monitor) Start() (err error) {
	err = this.ModuleBase.Start()
	this.service.RegisterGO(Rpc_GetServiceMonitorInfo, this.Rpc_GetServiceMonitorInfo)
	return
}

func (this *Monitor) OnInstallComp() {
	this.ModuleBase.OnInstallComp()
	this.ServiceMonitorComp = this.RegisterComp(new(ServiceMonitorComp)).(*ServiceMonitorComp)
}

//注册服务配置信息
func (this *Monitor) RegisterServiceSettingItem(name string, iswrite bool, value interface{}, f func(newvalue string) (err error)) {
	this.ServiceMonitor.Setting[name] = &core.SettingItem{
		ItemName: name,
		IsWrite:  iswrite,
		Data:     value,
	}
	this.ServiceSetting[name] = f
}

//注册模块配置信息
func (this *Monitor) RegisterModuleSettingItem(module core.M_Modules, name string, iswrite bool, value interface{}, f func(newvalue string) (err error)) {
	this.ServiceMonitor.ModuleMonitor[module].Setting[name] = &core.SettingItem{
		ItemName: name,
		IsWrite:  iswrite,
		Data:     value,
	}
	this.ModulesSetting[module][name] = f
}

//RPC------------------------------------------------------------------------------------------------------------------------------------------
//读取服务监控信息
func (this *Monitor) Rpc_GetServiceMonitorInfo() (result *core.ServiceMonitor, err string) {
	return this.ServiceMonitor, ""
}
