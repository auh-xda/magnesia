package main

import (
	"fmt"

	"github.com/auh-xda/magnesia/helpers/console"
	"github.com/shirou/gopsutil/v3/process"
)

func (agent Magnesia) ProcessList() []ProcessInfo {
	console.Info("getting the process list")

	processes, _ := process.Processes()

	var processList []ProcessInfo

	for _, p := range processes {
		name, _ := p.Name()
		exe, _ := p.Exe()
		cmdline, _ := p.Cmdline()
		username, _ := p.Username()
		status, _ := p.Status()
		ppid, _ := p.Ppid()
		createTime, _ := p.CreateTime()
		cpuPercent, _ := p.CPUPercent()
		memInfo, _ := p.MemoryInfo()
		numThreads, _ := p.NumThreads()
		nice, _ := p.Nice()

		processList = append(processList, ProcessInfo{
			PID:        p.Pid,
			PPID:       ppid,
			Name:       name,
			Exe:        exe,
			Cmdline:    cmdline,
			Username:   username,
			Status:     status[0],
			CPUPercent: cpuPercent,
			MemoryMB:   float32(memInfo.RSS) / 1024 / 1024,
			CreateTime: createTime,
			NumThreads: numThreads,
			Nice:       nice,
		})
	}

	console.Success(fmt.Sprintf("%d processes running", len(processList)))

	return processList
}
