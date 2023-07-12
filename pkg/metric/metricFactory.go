package metric

import (
	"github.com/qiqiuyang/env-beat/pkg/metric/enums/metricType"
	"sync"
)

var (
	metricMap sync.Map
)

func Register[T MetricParam, R MetricData, F MetricFormatData](code metricType.MetricType, metricService MetricService[T, R, F]) {
	metricMap.Store(code, metricService)
}

func GetProcessor[T MetricParam, R MetricData, F MetricFormatData](code metricType.MetricType) MetricService[T, R, F] {
	if processor, ok := metricMap.Load(code); ok {
		return processor.(MetricService[T, R, F])
	}

	return nil
}
