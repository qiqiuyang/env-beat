package metric

import "github.com/qiqiuyang/env-beat/pkg/metric/model"

type MetricParam interface {
	model.NumMetricParam | model.CpuMetricParam | model.MemoryMetricParam | model.LoadMetricParam
}

type MetricData interface {
	model.NumMetricData | model.CpuMetricData | model.MemoryMetricData | model.LoadMetricData
}

type MetricFormatData interface {
	model.NumMetricFormatData | model.CpuMetricFormatData | model.MemoryMetricFormatData | model.LoadMetricFormatData
}

type MetricService[T MetricParam, R MetricData, F MetricFormatData] interface {
	// FetchMetricData 获取metric数据
	FetchMetricData(request T)
	// GetMetricData 获取metric数据metric原始数据
	GetMetricData() (result R, err error)
	// FormatMetricData 格式化并返回值
	GetFormatMetricData() (result F, err error)
}
