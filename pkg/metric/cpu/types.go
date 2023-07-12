package cpu

import "github.com/qiqiuyang/env-beat/pkg/metric/model"

type CpuMetricData struct {
	NewMetric *model.CpuMetricData
	oldSample *model.CpuMetricData
	count     int
}
