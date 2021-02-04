package console

import (
	"fmt"
	"sync"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/cron"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"gonum.org/v1/gonum/floats"
)

//主机信息监控
type HostMonitorComp struct {
	cbase.ModuleCompBase
	module      *Console
	lock        sync.RWMutex
	hostInfo    *HostInfo
	cpuInfo     []*CpuInfo
	memoryInfo  *MemoryInfo
	hostMonitor *HostMonitor
}

func (this *HostMonitorComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, options)
	this.module = module.(*Console)
	this.getHostInfo()
	this.getCpuInfo()
	this.getMemoryInfo()
	this.hostMonitor = &HostMonitor{
		CpuUsageRate:    make([]float64, 60),
		MemoryUsageRate: make([]float64, 60),
	}
	cron.AddFunc("0 */1 * * * *", this.Monitor)         //每隔一分钟监听一次
	cron.AddFunc("1 0 */1 * * *", this.SaveMonitorData) //每小时保存一次数据
	return
}

//获取主机信息 只需一次即可
func (this *HostMonitorComp) getHostInfo() {
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
func (this *HostMonitorComp) getCpuInfo() {
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
			Flags:      fmt.Sprintf("%v", cpu.Flags),
			Microcode:  cpu.Microcode,
		}
	}
}

//获取Cpu信息 只需一次即可
func (this *HostMonitorComp) getMemoryInfo() {
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
func (this *HostMonitorComp) Monitor() {
	this.lock.Lock()
	defer this.lock.Unlock()
	Minute := time.Now().Minute()
	cpuinfo, _ := cpu.Percent(time.Second, false)
	if cpuinfo != nil && len(cpuinfo) > 0 {
		this.hostMonitor.CpuUsageRate[Minute] = cpuinfo[0]
	}
	memoryinfo, _ := mem.VirtualMemory()
	if memoryinfo != nil {
		this.hostMonitor.MemoryUsageRate[Minute] = memoryinfo.UsedPercent
	}
	log.Debugf("Monitor Minute:%d cpuinfo:%v memoryinfo:%v", Minute, this.hostMonitor.CpuUsageRate, this.hostMonitor.MemoryUsageRate)
}

func (this *HostMonitorComp) SaveMonitorData() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.module.cache.AddNewHostMonitor(this.hostMonitor)
	this.hostMonitor = &HostMonitor{
		CpuUsageRate:    make([]float64, 60),
		MemoryUsageRate: make([]float64, 60),
	}
	Hour := time.Now().Hour()
	log.Debugf("SaveMonitorData Hour:%d", Hour)
}

//读取主机监控数据
func (this *HostMonitorComp) GetHostMonitorData(queryTime QueryMonitorTime) (result *QueryHostMonitorDataResp) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	var (
		index          int
		step           int
		Minute         int
		totalMinute    int
		chartleng      int
		cpudata        []float64
		memorydata     []float64
		hostMonitor    []*HostMonitor
		stepcpudata    [][]float64
		stepmemorydata [][]float64
	)
	chartleng = 20
	stepcpudata = make([][]float64, chartleng)
	stepmemorydata = make([][]float64, chartleng)
	Minute = time.Now().Minute()
	result = &QueryHostMonitorDataResp{
		CurrCpuPer:    this.hostMonitor.CpuUsageRate[Minute],
		CurrMemoryPer: this.hostMonitor.MemoryUsageRate[Minute],
		Keys:          make([]string, chartleng),
		Cpu:           make([]float64, chartleng),
		Memory:        make([]float64, chartleng),
	}
	if queryTime == QueryMonitorTime_OneHour {
		step = 60 / chartleng
		totalMinute = 60
		hostMonitor, _ = this.module.cache.GetHostMonitor(1)
	} else if queryTime == QueryMonitorTime_SixHour {
		step = 60 / chartleng * 6
		totalMinute = 60 * 6
		hostMonitor, _ = this.module.cache.GetHostMonitor(6)
	} else if queryTime == QueryMonitorTime_OneDay {
		step = 60 / chartleng * 24
		totalMinute = 60 * 24
		hostMonitor, _ = this.module.cache.GetHostMonitor(24)
	} else {
		step = 60 / chartleng * 24 * 7
		totalMinute = 60 * 24 * 7
		hostMonitor, _ = this.module.cache.GetHostMonitor(24 * 7)
	}
	cpudata = make([]float64, totalMinute)
	memorydata = make([]float64, totalMinute)
	index = totalMinute - 1
	for i := Minute; i >= 0; i-- {
		cpudata[index] = this.hostMonitor.CpuUsageRate[i]
		memorydata[index] = this.hostMonitor.MemoryUsageRate[i]
		index--
	}
	if index > 0 {
		for i1 := 0; i1 < len(hostMonitor); i1++ {
			v := hostMonitor[i1]
			for i := 59; i >= 0; i-- {
				cpudata[index] = v.CpuUsageRate[i]
				memorydata[index] = v.MemoryUsageRate[i]
				index--
				if index < 0 {
					break
				}
			}
			if index < 0 {
				break
			}
		}
	}

	index = chartleng - 1
	for i := totalMinute; i > 0; i = i - step {
		start := i - step
		end := i
		stepcpudata[index] = cpudata[start:end]
		stepmemorydata[index] = memorydata[start:end]
		index--
	}
	for i := chartleng - 1; i >= 0; i-- {
		_time := time.Now().Add(time.Duration((chartleng-1-i)*-1*step) * time.Minute)
		result.Keys[i] = fmt.Sprintf("%d:%d", _time.Hour(), _time.Minute())
		result.Cpu[i] = floats.Max(stepcpudata[i])
		result.Memory[i] = floats.Max(stepmemorydata[i])
	}
	return
}
