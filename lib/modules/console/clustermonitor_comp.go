package console

import (
	"fmt"

	"sync"
	"time"

	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib/modules/monitor"
	"github.com/liwei1dao/lego/sys/cron"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/registry"
	"gonum.org/v1/gonum/floats"
)

//集群信息监控
type ClusterMonitorComp struct {
	cbase.ModuleCompBase
	service        base.IClusterService
	module         IConsole
	lock           sync.RWMutex
	servicemonitor map[string]*core.ServiceMonitor
	clusterMonitor map[string]*ClusterMonitor
}

func (this *ClusterMonitorComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, options)
	this.service = service.(base.IClusterService)
	this.module = module.(IConsole)
	this.servicemonitor = make(map[string]*core.ServiceMonitor)
	this.clusterMonitor = make(map[string]*ClusterMonitor)
	cron.AddFunc("0 */1 * * * *", this.Monitor)
	cron.AddFunc("1 0 */1 * * *", this.SaveMonitorData) //每小时保存一次数据
	return
}

func (this *ClusterMonitorComp) Monitor() {
	this.lock.Lock()
	defer this.lock.Unlock()
	Minute := time.Now().Minute()
	services := registry.GetAllServices()
	for _, s := range services {
		if result, err := this.service.RpcInvokeById(s.Id, monitor.Rpc_GetServiceMonitorInfo, true); err == nil {
			this.servicemonitor[s.Id] = result.(*core.ServiceMonitor)
			if _, ok := this.clusterMonitor[s.Id]; !ok {
				this.clusterMonitor[s.Id] = &ClusterMonitor{
					CpuUsageRate:    make([]float64, 60),
					MemoryUsageRate: make([]float64, 60),
					GoroutineUsed:   make([]float64, 60),
					PreWeight:       make([]float64, 60),
				}
			}
			this.clusterMonitor[s.Id].CpuUsageRate[Minute] = this.servicemonitor[s.Id].CpuUsed
			this.clusterMonitor[s.Id].MemoryUsageRate[Minute] = this.servicemonitor[s.Id].MemoryUsed
			this.clusterMonitor[s.Id].GoroutineUsed[Minute] = float64(this.servicemonitor[s.Id].TotalGoroutine)
			this.clusterMonitor[s.Id].PreWeight[Minute] = this.servicemonitor[s.Id].CurrPreWeight
		}
	}
	log.Debugf("Monitor Minute:%d clusterMonitor:%v", Minute, this.clusterMonitor)
}

func (this *ClusterMonitorComp) SaveMonitorData() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.module.Cache().AddNewClusterMonitor(this.clusterMonitor)
	for k, _ := range this.clusterMonitor {
		this.clusterMonitor[k] = &ClusterMonitor{
			CpuUsageRate:    make([]float64, 60),
			MemoryUsageRate: make([]float64, 60),
			GoroutineUsed:   make([]float64, 60),
			PreWeight:       make([]float64, 60),
		}
	}
}

func (this *ClusterMonitorComp) GetClusterMonitorDataResp(queryTime QueryMonitorTime) (result map[string]map[string]interface{}) {
	result = make(map[string]map[string]interface{})
	for k, v := range this.servicemonitor {
		result[k] = make(map[string]interface{})
		result[k]["Info"] = v
		result[k]["Monitor"] = this.getClusterMonitorData(v, queryTime)
	}
	return
}

//读取主机监控数据
func (this *ClusterMonitorComp) getClusterMonitorData(service *core.ServiceMonitor, queryTime QueryMonitorTime) (result *ClusterMonitorData) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	var (
		index             int
		step              int
		Minute            int
		totalMinute       int
		chartleng         int
		cpudata           []float64
		memorydata        []float64
		goroutinedata     []float64
		preWeightdata     []float64
		clusterMonitor    []*ClusterMonitor
		stepcpudata       [][]float64
		stepmemorydata    [][]float64
		stepgoroutinedata [][]float64
		steppreWeightdata [][]float64
	)
	chartleng = 20
	stepcpudata = make([][]float64, chartleng)
	stepmemorydata = make([][]float64, chartleng)
	stepgoroutinedata = make([][]float64, chartleng)
	steppreWeightdata = make([][]float64, chartleng)
	Minute = time.Now().Minute()
	result = &ClusterMonitorData{
		CurrCpuPer:    service.CurrPreWeight,
		CurrMemoryPer: service.MemoryUsed,
		CurrGoroutine: float64(service.TotalGoroutine),
		CurrPreWeight: service.CurrPreWeight,
		Keys:          make([]string, chartleng),
		Cpu:           make([]float64, chartleng),
		Memory:        make([]float64, chartleng),
		Goroutine:     make([]float64, chartleng),
		PreWeight:     make([]float64, chartleng),
	}
	if queryTime == QueryMonitorTime_OneHour {
		step = 60 / chartleng
		totalMinute = 60
		clusterMonitor, _ = this.module.Cache().GetClusterMonitor(service.ServiceId, 1)
	} else if queryTime == QueryMonitorTime_SixHour {
		step = 60 / chartleng * 6
		totalMinute = 60 * 6
		clusterMonitor, _ = this.module.Cache().GetClusterMonitor(service.ServiceId, 6)
	} else if queryTime == QueryMonitorTime_OneDay {
		step = 60 / chartleng * 24
		totalMinute = 60 * 24
		clusterMonitor, _ = this.module.Cache().GetClusterMonitor(service.ServiceId, 24)
	} else {
		step = 60 / chartleng * 24 * 7
		totalMinute = 60 * 24 * 7
		clusterMonitor, _ = this.module.Cache().GetClusterMonitor(service.ServiceId, 24*7)
	}
	cpudata = make([]float64, totalMinute)
	memorydata = make([]float64, totalMinute)
	goroutinedata = make([]float64, totalMinute)
	preWeightdata = make([]float64, totalMinute)
	index = totalMinute - 1
	for i := Minute; i >= 0; i-- {
		cpudata[index] = this.clusterMonitor[service.ServiceId].CpuUsageRate[i]
		memorydata[index] = this.clusterMonitor[service.ServiceId].MemoryUsageRate[i]
		goroutinedata[index] = this.clusterMonitor[service.ServiceId].GoroutineUsed[i]
		preWeightdata[index] = this.clusterMonitor[service.ServiceId].PreWeight[i]
		index--
	}
	if index > 0 {
		for i1 := 0; i1 < len(clusterMonitor); i1++ {
			v := clusterMonitor[i1]
			for i := 59; i >= 0; i-- {
				cpudata[index] = v.CpuUsageRate[i]
				memorydata[index] = v.MemoryUsageRate[i]
				goroutinedata[index] = v.GoroutineUsed[i]
				preWeightdata[index] = v.PreWeight[i]
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
		stepgoroutinedata[index] = goroutinedata[start:end]
		steppreWeightdata[index] = preWeightdata[start:end]
		index--
	}
	for i := chartleng - 1; i >= 0; i-- {
		_time := time.Now().Add(time.Duration((chartleng-1-i)*-1*step) * time.Minute)
		result.Keys[i] = fmt.Sprintf("%d:%d", _time.Hour(), _time.Minute())
		result.Cpu[i] = floats.Max(stepcpudata[i])
		result.Memory[i] = floats.Max(stepmemorydata[i])
		result.Goroutine[i] = floats.Max(stepgoroutinedata[i])
		result.PreWeight[i] = floats.Max(steppreWeightdata[i])
	}
	return
}
