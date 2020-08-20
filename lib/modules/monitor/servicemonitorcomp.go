package monitor

import (
	"os"
	"runtime"

	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/cron"
	"github.com/shirou/gopsutil/process"
)

//进程统计组件
type ServiceMonitorComp struct {
	cbase.ModuleCompBase
	module       *Monitor
	service      base.IClusterService
	Process      *process.Process
	MonitorNum   uint32
	MonitorTotal uint32
}

func (this *ServiceMonitorComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, setting map[string]interface{}) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, setting)
	this.module = module.(*Monitor)
	this.service = service.(base.IClusterService)
	this.MonitorNum = 0
	this.MonitorTotal = 360
	this.Process, err = process.NewProcess(int32(os.Getpid()))
	if err == nil {
		pname, _ := this.Process.Name()
		this.module.ServiceMonitor.Pid = this.Process.Pid
		this.module.ServiceMonitor.Pname = pname
		this.module.ServiceMonitor.MemoryInfo = make([]float32, this.MonitorTotal)
		this.module.ServiceMonitor.CpuInfo = make([]float64, this.MonitorTotal)
		this.module.ServiceMonitor.TotalGoroutine = make([]int, this.MonitorTotal)
	}
	return
}

func (this *ServiceMonitorComp) Start() (err error) {
	err = this.ModuleCompBase.Start()
	cron.AddFunc("@every 1m", this.RefreshMonitorInfo)
	return
}

func (this *ServiceMonitorComp) RefreshMonitorInfo() {
	Memory, _ := this.Process.MemoryPercent()
	Cpu, _ := this.Process.CPUPercent()
	if this.MonitorNum >= this.MonitorTotal {
		this.module.ServiceMonitor.ServicePreWeight = this.service.GetPreWeight()
		this.module.ServiceMonitor.TotalGoroutine = append(this.module.ServiceMonitor.TotalGoroutine, runtime.NumGoroutine())
		this.module.ServiceMonitor.MemoryInfo = append(this.module.ServiceMonitor.MemoryInfo, Memory)
		this.module.ServiceMonitor.CpuInfo = append(this.module.ServiceMonitor.CpuInfo, Cpu)
		this.module.ServiceMonitor.TotalGoroutine = this.module.ServiceMonitor.TotalGoroutine[1:]
		this.module.ServiceMonitor.MemoryInfo = this.module.ServiceMonitor.MemoryInfo[1:]
		this.module.ServiceMonitor.CpuInfo = this.module.ServiceMonitor.CpuInfo[1:]
	} else {
		this.module.ServiceMonitor.ServicePreWeight = this.service.GetPreWeight()
		this.module.ServiceMonitor.TotalGoroutine[this.MonitorNum] = runtime.NumGoroutine()
		this.module.ServiceMonitor.MemoryInfo[this.MonitorNum] = Memory
		this.module.ServiceMonitor.CpuInfo[this.MonitorNum] = Cpu
		this.MonitorNum++
	}
}
