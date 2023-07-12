package cpuCountType

type CpuCountType int

const (
	// CpuCountDefault 从默认的获取的获取方式获取
	CpuCountDefault CpuCountType = iota
	// CpuCountRuntime 从go的runtime获取
	CpuCountRuntime
	// CpuCountOnline online状态的cpu：表示可以被调度器使用。
	CpuCountOnline
	// CpuCountPresent present状态的cpu：表示已经被kernel接管。
	CpuCountPresent
	// CpuCountEnv 从系统环境变量读取cpu个数
	CpuCountEnv

	// 此外还有两种状态
	// possible状态的cpu：可理解为存在这个CPU资源，但还没有纳入Kernel的管理范围
	// active状态的cpu：表示可以被迁移migrate。
)
