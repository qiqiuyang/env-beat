package metricType

type MetricType int

const (
	// CpuNumMetric
	CpuNumMetric MetricType = iota + 1
	// CpuMetric
	CpuMetric
	// MemoryMetric
	MemoryMetric
	LoadMetric
)
