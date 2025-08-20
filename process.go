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

		var statusRes string
		status, e := p.Status()
		if e != nil || len(status) == 0 {
			statusRes = "unknown"
		} else {
			statusRes = status[0]
		}

		createTime, _ := p.CreateTime()
		cpuPercent, _ := p.CPUPercent()
		mem, e := p.MemoryInfo()

		var memInfo float32

		if e != nil {
			memInfo = 0.0
		} else {
			memInfo = float32(mem.RSS) / 1024 / 1024
		}
		numThreads, e := p.NumThreads()
		nice, e := p.Nice()
		ppid, _ := p.Ppid()

		pInfo := ProcessInfo{
			PID:        p.Pid,
			PPID:       ppid,
			Name:       name,
			Exe:        exe,
			Cmdline:    cmdline,
			Username:   username,
			Status:     statusRes,
			CPUPercent: cpuPercent,
			MemoryMB:   memInfo,
			CreateTime: createTime,
			NumThreads: numThreads,
			Nice:       nice,
		}

		console.Log(pInfo)

		processList = append(processList, pInfo)

	}

	console.Success(fmt.Sprintf("%d processes running", len(processList)))

	return processList
}
