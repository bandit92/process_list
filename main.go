package main

import (
    "sort"
    "strings"
	"fmt"
	"flag"
    "github.com/alexeyco/simpletable"
    "github.com/shirou/gopsutil/process"
    "path/filepath"
)

type HostProcs struct {
	Name       string   `json:"name,omitempty"`
	PID        int32    `json:"pid,omitempty"`
	PPID       int32    `json:"ppid,omitempty"`
	Mem        float32  `json:"memory,omitempty"`
	Nice       int32    `json:"nice,omitempty"`
	Priority   int32    `json:"priority,omitempty"`
	CPU        float64  `json:"cpu,omitempty"`
	User       string   `json:"user,omitempty"`
	Cwd        string   `json:"cwd,omitempty"`
	Env        []string `json:"env,omitempty"`
	Cmdline    string   `json:"cmdline,omitempty"`
	CreateTime int64    `json:"createtime,omitempty"`
	Exe        string   `json:"exe,omitempty"`
}


func ListProcess() []HostProcs {
	var procs []HostProcs
	pids, err := process.Pids()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	for _, p := range pids {
		pid, _ := process.NewProcess(p)
		name, _ := pid.Name()
		nice, _ := pid.Nice()
		priority, _ := pid.IOnice()
		mem, _ := pid.MemoryPercent()
		cpu, _ := pid.CPUPercent()
		ppid, _ := pid.Ppid()
		user, _ := pid.Username()
		cwd, _ := pid.Cwd()
		env, _ := pid.Environ()
		exe, _ := pid.Exe()
		cmdline, _ := pid.Cmdline()
		createtime, _ := pid.CreateTime()
		pcs := HostProcs{Name: name,
			PID:        p,
			PPID:       ppid,
			Nice:       nice,
			Priority:   priority,
			Mem:        mem,
			CPU:        cpu,
			Cwd:        cwd,
			Env:        env,
			Cmdline:    cmdline,
			CreateTime: createtime,
			Exe:        exe,
			User:       user,
		}
		procs = append(procs, pcs)
	}
	return procs
}

func reverseSortUser(procs []HostProcs) {
	sort.Slice(procs, func(i, j int) bool {
		return procs[i].User > procs[j].User
	})
}

func sortByUser(procs []HostProcs) {
	sort.Slice(procs, func(i, j int) bool {
		return procs[i].User < procs[j].User
	})
}

// Find all processes running for a given User
func filterByUser(procs []HostProcs, user string) []HostProcs {
	var filterProcs []HostProcs
	for _, p := range procs {
		if p.User == user {
			filterProcs = append(filterProcs, p)
		}
	}
	return filterProcs
}

// Find a specific PID
func filterByPID(procs []HostProcs, pid int32) []HostProcs {
	var filterProcs []HostProcs
	for _, p := range procs {
		if p.PID == pid {
			filterProcs = append(filterProcs, p)
		} //end
	} //end
	return filterProcs
} //emd

// Lazy Match a specific process(es) by name
func filterByName(procs []HostProcs, name string) []HostProcs {
	var filterProcs []HostProcs
	for _, p := range procs {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(name)) {
			filterProcs = append(filterProcs, p)
		}
	}
	return filterProcs
} //end

// Find all processes with Parent Process ID
func fitlerByPPID(procs []HostProcs, ppid int32) []HostProcs {
	var filterProcs []HostProcs
	for _, p := range procs {
		if p.PPID == ppid {
			filterProcs = append(filterProcs, p)
		}
	} //end
	return filterProcs
}

//Print General Table
func printTable(procs []HostProcs) string {
	tbl := simpletable.New()
	if len(procs) == 0 {
		return "[!] No Processes Found!"
	} else {
		tbl.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: "PID"},
				{Align: simpletable.AlignLeft, Text: "Name"},
				{Align: simpletable.AlignLeft, Text: "PPID"},
				{Align: simpletable.AlignLeft, Text: "USER"},
				{Align: simpletable.AlignLeft, Text: "EXE"},
				{Align: simpletable.AlignLeft, Text: "Env"},
                {Align: simpletable.AlignLeft, Text: "Command Line"},
			},
		}
		for _, row := range procs {

			r := []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", row.PID)},
				{Align: simpletable.AlignLeft, Text: row.Name},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", row.PPID)},
				{Align: simpletable.AlignCenter, Text: row.User},
            	{Align: simpletable.AlignLeft, Text: filepath.Base(row.Exe)},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", row.Env)},
                {Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", row.Cmdline)},
			}
			tbl.Body.Cells = append(tbl.Body.Cells, r)
		}
	} //end
	return tbl.String()
}

func main(){
	pidPtr := flag.Int("pid",-1,"PID for process to search")
	ppidPtr := flag.Int("ppid",-1,"PPID for parent process to search")
	namePtr := flag.String("name","","Process name to search")
	userPtr := flag.String("user","","User to search for processes under")
	flag.Parse()

	user := *userPtr
	pid := int32(*pidPtr)
	ppid := int32(*ppidPtr)
	name := *namePtr

	// Grab all processes
	procs := ListProcess()


	if pid != -1 {
		procs = filterByPID(procs,pid)
	}//end
	if ppid != -1 {
		procs = filterByPID(procs, ppid)
	}//end
	if user != "" {
		procs = filterByUser(procs, user)
	}//end
	if name != "" {
		procs = filterByName(procs,name)
	}

	data := printTable(procs)
	fmt.Println(data)
}
