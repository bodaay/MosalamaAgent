package monitoring

import (
	"fmt"
	"time"

	nvml "github.com/mindprince/gonvml"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/v3/cpu"
)

type ResourceMonitor struct {
}

func NewResourceMonitor() *ResourceMonitor {
	return &ResourceMonitor{}
}
func (rm *ResourceMonitor) GetCPUUsage() (float64, error) {
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}
	if len(percentages) > 0 {
		return percentages[0], nil
	}
	return 0, fmt.Errorf("unable to get CPU usage")
}
func (rm *ResourceMonitor) GetMemoryUsage() (uint64, uint64, float64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, 0, err
	}
	return vmStat.Total, vmStat.Used, vmStat.UsedPercent, nil
}
func (rm *ResourceMonitor) GetDiskUsage(path string) (uint64, uint64, float64, error) {
	usageStat, err := disk.Usage(path)
	if err != nil {
		return 0, 0, 0, err
	}
	return usageStat.Total, usageStat.Used, usageStat.UsedPercent, nil
}

func (rm *ResourceMonitor) GetGPUUsage() ([]float64, error) {
	err := nvml.Initialize()
	if err != nil {
		return nil, err
	}
	defer nvml.Shutdown()

	deviceCount, err := nvml.DeviceCount()
	if err != nil {
		return nil, err
	}

	var gpuUsages []float64

	for i := uint(0); i < deviceCount; i++ {
		device, err := nvml.DeviceHandleByIndex(i)
		if err != nil {
			return nil, err
		}

		gpuUtil, _, err := device.UtilizationRates()
		if err != nil {
			return nil, err
		}

		gpuUsages = append(gpuUsages, float64(gpuUtil))
	}

	return gpuUsages, nil
}
func (rm *ResourceMonitor) StartMonitoring(interval time.Duration, stopChan <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			rm.collectAndLogMetrics()
		case <-stopChan:
			fmt.Println("Stopping monitoring")
			return
		}
	}
}

func (rm *ResourceMonitor) collectAndLogMetrics() {
	cpuUsage, err := rm.GetCPUUsage()
	if err != nil {
		fmt.Printf("Error collecting CPU usage: %v\n", err)
	} else {
		fmt.Printf("CPU Usage: %.2f%%\n", cpuUsage)
	}
	// Similarly, collect memory, disk, and GPU usage
}
