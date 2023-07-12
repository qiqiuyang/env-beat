package load

import (
	"github.com/qiqiuyang/env-beat/pkg/logutil"
	"github.com/qiqiuyang/env-beat/pkg/metric"
	"github.com/qiqiuyang/env-beat/pkg/metric/enums/metricType"
	"github.com/qiqiuyang/env-beat/pkg/metric/model"
	"github.com/shopspring/decimal"
	"os"
	"strconv"
	"strings"
)

func init() {
	metric.Register[model.LoadMetricParam, model.LoadMetricData, model.LoadMetricFormatData](metricType.LoadMetric, &LoadMetricData{})
}

func (l *LoadMetricData) FetchMetricData(param model.LoadMetricParam) {
	err := l.get()
	if err != nil {
		if logutil.GetSugarLogger() != nil {
			logutil.GetSugarLogger().Errorf("Fetch failed, err: %s", err.Error())
		}
	}

}

func (l *LoadMetricData) GetMetricData() (result model.LoadMetricData, err error) {
	result.One = l.One
	result.Five = l.Five
	result.Fifteen = l.Fifteen
	return result, nil
}

func (l *LoadMetricData) GetFormatMetricData() (result model.LoadMetricFormatData, err error) {
	cpuCount, _ := metric.GetProcessor[model.NumMetricParam, model.NumMetricData, model.NumMetricFormatData](metricType.CpuNumMetric).GetMetricData()
	result.One, _ = decimal.NewFromFloat(l.One).Round(4).Float64()
	result.Five, _ = decimal.NewFromFloat(l.Five).Round(4).Float64()
	result.Fifteen, _ = decimal.NewFromFloat(l.Fifteen).Round(4).Float64()
	result.OneMinute, _ = decimal.NewFromFloat(l.One / float64(cpuCount.CpuNum)).Round(4).Float64()
	result.FiveMinute, _ = decimal.NewFromFloat(l.Five / float64(cpuCount.CpuNum)).Round(4).Float64()
	result.FifteenMinute, _ = decimal.NewFromFloat(l.Fifteen / float64(cpuCount.CpuNum)).Round(4).Float64()
	return result, nil
}

func (l *LoadMetricData) get() error {
	line, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return nil
	}

	fields := strings.Fields(string(line))

	l.One, _ = strconv.ParseFloat(fields[0], 64)
	l.Five, _ = strconv.ParseFloat(fields[1], 64)
	l.Fifteen, _ = strconv.ParseFloat(fields[2], 64)

	return nil
}
