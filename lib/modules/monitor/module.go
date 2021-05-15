package monitor

import (
	"fmt"
	"os"
	"runtime"

	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib"
	"github.com/shirou/gopsutil/process"
)

func NewModule() core.IServiceMonitor {
	m := new(Monitor)
	return m
}

type Monitor struct {
	cbase.ModuleBase
	service        base.IClusterService
	ServiceMonitor *core.ServiceMonitor
	ServiceSetting map[string]func(newvalue string) (err error)
	ModulesSetting map[core.M_Modules]map[string]func(newvalue string) (err error)
	Process        *process.Process
}

func (this *Monitor) GetType() core.M_Modules {
	return lib.SM_MonitorModule
}

func (this *Monitor) NewOptions() (options core.IModuleOptions) {
	return new(Options)
}

func (this *Monitor) Init(service core.IService, module core.IModule, options core.IModuleOptions) (err error) {
	this.service = service.(base.IClusterService)
	this.Process, err = process.NewProcess(int32(os.Getpid()))
	this.ServiceMonitor = &core.ServiceMonitor{
		ServiceId:       this.service.GetId(),
		ServiceType:     this.service.GetType(),
		ServiceCategory: this.service.GetCategory(),
		ServiceVersion:  this.service.GetVersion(),
		ServiceTag:      this.service.GetTag(),
		Setting:         make(map[string]*core.SettingItem),
		SysSetting:      make(map[string]*core.SysSetting),
		ModuleMonitor:   make(map[core.M_Modules]*core.ModuleMonitor),
	}
	this.ServiceSetting = make(map[string]func(newvalue string) (err error))
	this.ModulesSetting = make(map[core.M_Modules]map[string]func(newvalue string) (err error))
	err = this.ModuleBase.Init(service, module, options)
	for k, v := range this.service.GetSettings().Settings {
		this.ServiceMonitor.Setting[k] = &core.SettingItem{
			ItemName: k,
			IsWrite:  false,
			Data:     v,
		}
	}
	for k, v := range this.service.GetSettings().Sys {
		this.ServiceMonitor.SysSetting[k] = &core.SysSetting{
			SysName: k,
			Setting: make(map[string]*core.SettingItem),
		}
		for k1, v1 := range v {
			this.ServiceMonitor.SysSetting[k].Setting[k1] = &core.SettingItem{
				ItemName: k1,
				IsWrite:  false,
				Data:     v1,
			}
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
	this.service.RegisterGO(Rpc_SetMonitorServiceSetting, this.Rpc_SetMonitorServiceSetting)
	this.service.RegisterGO(Rpc_SetMonitorModuleSetting, this.Rpc_SetMonitorModuleSetting)
	return
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
	memory, _ := this.Process.MemoryPercent()
	cpu, _ := this.Process.CPUPercent()
	this.ServiceMonitor.CpuUsed = cpu
	this.ServiceMonitor.MemoryUsed = float64(memory)
	this.ServiceMonitor.TotalGoroutine = runtime.NumGoroutine()
	this.ServiceMonitor.CurrPreWeight = this.service.GetPreWeight()
	return this.ServiceMonitor, ""
}

//读取服务监控信息
func (this *Monitor) Rpc_SetMonitorServiceSetting(key, value string) (result string, err string) {
	if f, ok := this.ServiceSetting[key]; ok {
		if e := f(value); e != nil {
			return "", fmt.Sprintf("modifier key:%s err:%s", key, e.Error())
		}
	} else {
		return "", fmt.Sprintf("no register key:%s modifier", key)
	}
	return "", ""
}

//读取服务监控信息
func (this *Monitor) Rpc_SetMonitorModuleSetting(module, key, value string) (result string, err string) {
	if m, ok := this.ModulesSetting[core.M_Modules(module)]; ok {
		if f, ok := m[key]; ok {
			if e := f(value); e != nil {
				return "", fmt.Sprintf("modifier key:%s err:%s", key, e.Error())
			}
		} else {
			return "", fmt.Sprintf("no register module:%s key:%s modifier", module, key)
		}
	} else {
		return "", fmt.Sprintf("no register module:%s modifier", module)
	}
	return "", ""
}
