package host

import (
	"fmt"
	"os"
	"runtime"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

type Disk struct {
	Device     string
	FSType     string
	TotalSpace uint64
	UsedSpace  uint64
	FreeSpace  uint64
}

type InfoProcess struct {
	IsRunning     bool
	Name          string
	MemoryPercent float32
	CPUPercent    float64
}

type SystemInfo struct {
	CPU      int    `json:"cpu"`
	Memory   string `json:"memory"`
	HostName string `json:"hostname"`
}

func GetHostInfo() []SystemInfo {
	var response []SystemInfo

	numCPU := runtime.NumCPU()
	virtualMem, _ := mem.VirtualMemory()
	totalRAM := float64(virtualMem.Total) / (1024 * 1024 * 1024)
	totalMemory := fmt.Sprintf("%.2f", totalRAM)
	hostname, _ := os.Hostname()
	response = append(response, SystemInfo{
		CPU:      numCPU,
		Memory:   totalMemory,
		HostName: hostname,
	})
	return response

}

func GetInfoProcess() []InfoProcess {
	p, _ := process.Processes()
	var response []InfoProcess

	for _, pr := range p {
		Name, _ := pr.Name()
		CPUPercent, _ := pr.CPUPercent()
		MemoryPercent, _ := pr.MemoryPercent()
		IsRunning, _ := pr.IsRunning()

		response = append(response, InfoProcess{
			Name:          Name,
			CPUPercent:    CPUPercent,
			MemoryPercent: MemoryPercent,
			IsRunning:     IsRunning,
		})

	}
	return response

}

func GetDisks() []Disk {
	var response []Disk
	partitions, _ := disk.Partitions(true)

	for _, part := range partitions {
		diskStat, _ := disk.Usage(part.Device)
		response = append(response, Disk{
			Device:     part.Device,
			FSType:     part.Fstype,
			TotalSpace: diskStat.Total / 1024 / 1024 / 1024,
			UsedSpace:  diskStat.Used / 1024 / 1024 / 1024,
			FreeSpace:  diskStat.Free / 1024 / 1024 / 1024,
		})
	}
	return response
}
