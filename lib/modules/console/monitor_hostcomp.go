package console

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/cron"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

//主机信息监控
type MonitorHostcomp struct {
	cbase.ModuleCompBase
	hostInfo   *HostInfo
	cpuInfo    []*CpuInfo
	memoryInfo *MemoryInfo
}

func (this *MonitorHostcomp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, options)
	this.GetHostInfo()
	this.GetCpuInfo()
	this.GeMemoryInfo()
	cron.AddFunc("0 */1 * * * ?", this.Monitor) //每隔一分钟监听一次
	return
}

//获取主机信息 只需一次即可
func (this *MonitorHostcomp) GetHostInfo() {
	InfoStat, err := host.Info()
	if err != nil {
		log.Errorf("GetHostInfo err:%v", err)
		return
	}
	this.hostInfo = &HostInfo{
		HostID:               InfoStat.HostID,
		Hostname:             InfoStat.Hostname,
		Uptime:               InfoStat.Uptime,
		BootTime:             InfoStat.BootTime,
		Procs:                InfoStat.Procs,
		OS:                   InfoStat.OS,
		Platform:             InfoStat.Platform,
		PlatformFamily:       InfoStat.PlatformFamily,
		PlatformVersion:      InfoStat.PlatformVersion,
		KernelArch:           InfoStat.Hostname,
		VirtualizationSystem: InfoStat.VirtualizationSystem,
		VirtualizationRole:   InfoStat.VirtualizationRole,
	}
}

//获取Cpu信息 只需一次即可
func (this *MonitorHostcomp) GetCpuInfo() {
	cpus, err := cpu.Info()
	if err != nil {
		log.Errorf("GetCpuInfo err:%v", err)
		return
	}
	this.cpuInfo = make([]*CpuInfo, len(cpus))
	for i, cpu := range cpus {
		this.cpuInfo[i] = &CpuInfo{
			CPU:        cpu.CPU,
			VendorID:   cpu.VendorID,
			Family:     cpu.Family,
			Model:      cpu.Model,
			Stepping:   cpu.Stepping,
			PhysicalID: cpu.PhysicalID,
			CoreID:     cpu.CoreID,
			Cores:      cpu.Cores,
			ModelName:  cpu.ModelName,
			Mhz:        cpu.Mhz,
			CacheSize:  cpu.CacheSize,
			Flags:      cpu.Flags,
			Microcode:  cpu.Microcode,
		}
	}
}

//获取Cpu信息 只需一次即可
func (this *MonitorHostcomp) GeMemoryInfo() {
	memory, err := mem.VirtualMemory()
	if err != nil {
		log.Errorf("GeMemoryInfo err:%v", err)
		return
	}
	this.memoryInfo = &MemoryInfo{
		Total:       memory.Total,
		Available:   memory.Available,
		Used:        memory.Used,
		UsedPercent: memory.UsedPercent,
		Free:        memory.Free,
		Active:      memory.Active,
		Inactive:    memory.Inactive,
		Wired:       memory.Wired,
		Laundry:     memory.Laundry,
	}
}



//服务器监控
func (this *MonitorHostcomp) Monitor() {

}
