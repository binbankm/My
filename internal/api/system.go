package api

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type SystemInfo struct {
	Hostname        string  `json:"hostname"`
	OS              string  `json:"os"`
	Platform        string  `json:"platform"`
	PlatformVersion string  `json:"platformVersion"`
	KernelVersion   string  `json:"kernelVersion"`
	Uptime          uint64  `json:"uptime"`
	CPUCores        int     `json:"cpuCores"`
}

type SystemStats struct {
	CPU     []float64     `json:"cpu"`
	Memory  *mem.VirtualMemoryStat `json:"memory"`
	Disk    []*disk.UsageStat `json:"disk"`
	Network []net.IOCountersStat `json:"network"`
}

// GetSystemInfo returns basic system information
func GetSystemInfo(c *gin.Context) {
	hostInfo, err := host.Info()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	info := SystemInfo{
		Hostname:        hostInfo.Hostname,
		OS:              runtime.GOOS,
		Platform:        hostInfo.Platform,
		PlatformVersion: hostInfo.PlatformVersion,
		KernelVersion:   hostInfo.KernelVersion,
		Uptime:          hostInfo.Uptime,
		CPUCores:        runtime.NumCPU(),
	}

	c.JSON(http.StatusOK, info)
}

// GetSystemStats returns system resource statistics
func GetSystemStats(c *gin.Context) {
	// CPU usage
	cpuPercent, err := cpu.Percent(0, true)
	if err != nil {
		cpuPercent = []float64{}
	}

	// Memory usage
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get memory info"})
		return
	}

	// Disk usage
	partitions, _ := disk.Partitions(false)
	var diskUsages []*disk.UsageStat
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err == nil {
			diskUsages = append(diskUsages, usage)
		}
	}

	// Network stats
	netIO, _ := net.IOCounters(true)

	stats := SystemStats{
		CPU:     cpuPercent,
		Memory:  memInfo,
		Disk:    diskUsages,
		Network: netIO,
	}

	c.JSON(http.StatusOK, stats)
}
