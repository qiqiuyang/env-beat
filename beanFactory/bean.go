package beanFactory

import (
	"github.com/qiqiuyang/env-beat/pkg/metric"
	_ "github.com/qiqiuyang/env-beat/pkg/metric/cpu"
	_ "github.com/qiqiuyang/env-beat/pkg/metric/cpu_num"
	"github.com/qiqiuyang/env-beat/pkg/metric/enums/metricType"
	_ "github.com/qiqiuyang/env-beat/pkg/metric/load"
	_ "github.com/qiqiuyang/env-beat/pkg/metric/memory"
	"github.com/qiqiuyang/env-beat/pkg/metric/model"
)

var (
	CpuNumMetric = metric.GetProcessor[model.NumMetricParam, model.NumMetricData, model.NumMetricFormatData](metricType.CpuNumMetric)
	CpuMetric    = metric.GetProcessor[model.CpuMetricParam, model.CpuMetricData, model.CpuMetricFormatData](metricType.CpuMetric)
	MemoryMetric = metric.GetProcessor[model.MemoryMetricParam, model.MemoryMetricData, model.MemoryMetricFormatData](metricType.MemoryMetric)
	LoadMetric   = metric.GetProcessor[model.LoadMetricParam, model.LoadMetricData, model.LoadMetricFormatData](metricType.LoadMetric)
)
