package model

type CpuMetricParam struct {
}

type CpuMetricData struct {
	Totals CPU

	// list carries the same data, broken down by CPU
	List []CPU

	// CPUInfo carries some data from /proc/cpuinfo
	CPUInfo []CPUInfo
}

type CpuMetricFormatData struct {
	TotalPct   float64 `json:"totalPct"`
	User       uint64  `json:"user,omitempty"`
	UserPct    float64 `json:"userPct"`
	Sys        uint64  `json:"system,omitempty"`
	SysPct     float64 `json:"sysPct"`
	Idle       uint64  `json:"idle,omitempty"`
	IdlePct    float64 `json:"idlePct"`
	Nice       uint64  `json:"nice,omitempty"` // Linux, Darwin, BSD
	NicePct    float64 `json:"nicePct"`
	Irq        uint64  `json:"irq,omitempty"` // Linux and openbsd
	IrqPct     float64 `json:"irqPct"`
	Wait       uint64  `json:"iowait,omitempty"` // Linux and AIX
	WaitPct    float64 `json:"waitPct"`
	SoftIrq    uint64  `json:"softirq,omitempty"` // Linux only
	SofaIrqPct float64 `json:"sofaIrqPct"`
	Stolen     uint64  `json:"steal,omitempty"` // Linux only
	StplenPct  float64 `json:"stplenPct"`
}

type CPU struct {
	User    uint64 `struct:"user,omitempty"`
	Sys     uint64 `struct:"system,omitempty"`
	Idle    uint64 `struct:"idle,omitempty"`
	Nice    uint64 `struct:"nice,omitempty"`    // Linux, Darwin, BSD
	Irq     uint64 `struct:"irq,omitempty"`     // Linux and openbsd
	Wait    uint64 `struct:"iowait,omitempty"`  // Linux and AIX
	SoftIrq uint64 `struct:"softirq,omitempty"` // Linux only
	Stolen  uint64 `struct:"steal,omitempty"`   // Linux only
}

func (cpu CPU) Total() uint64 {
	// it's generally safe to blindly sum these up,
	// As we're just trying to get a total of all CPU time.
	return cpu.User + cpu.Nice + cpu.Sys + cpu.Idle + cpu.Wait + cpu.Irq + cpu.SoftIrq + cpu.Stolen
}

type CPUInfo struct {
	ModelName   string
	ModelNumber string
	Mhz         float64
	PhysicalID  int
	CoreID      int
}
