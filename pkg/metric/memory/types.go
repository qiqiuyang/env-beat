package memory

import "github.com/qiqiuyang/env-beat/pkg/metric/model"

type MemoryMetricData struct {
	model.MemoryMetricData
	errMsg error
}
