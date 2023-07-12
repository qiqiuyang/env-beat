package cpu

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/joeshaw/multierror"
	"github.com/qiqiuyang/env-beat/pkg/logutil"
	"github.com/qiqiuyang/env-beat/pkg/metric"
	"github.com/qiqiuyang/env-beat/pkg/metric/enums/metricType"
	"github.com/qiqiuyang/env-beat/pkg/metric/model"
	"github.com/shopspring/decimal"
	"strings"
)

func init() {
	metric.Register[model.CpuMetricParam, model.CpuMetricData, model.CpuMetricFormatData](metricType.CpuMetric, &CpuMetricData{})
}

func (c *CpuMetricData) FetchMetricData(param model.CpuMetricParam) {
	cpuMetric, err := Get()
	if err != nil {
		if logutil.GetSugarLogger() != nil {
			logutil.GetSugarLogger().Errorf("Fetch failed, err: %s", err.Error())
		}
		return
	}

	// 获取cpu核心数据
	cpuCount, _ := metric.GetProcessor[model.NumMetricParam, model.NumMetricData, model.NumMetricFormatData](metricType.CpuNumMetric).GetMetricData()

	// 更新数据
	c.oldSample = c.NewMetric
	c.NewMetric = &cpuMetric
	c.count = cpuCount.CpuNum
}

func (c *CpuMetricData) GetMetricData() (result model.CpuMetricData, err error) {
	result.Totals = c.NewMetric.Totals
	result.List = c.NewMetric.List
	result.CPUInfo = c.NewMetric.CPUInfo
	return result, nil
}

func (c *CpuMetricData) GetFormatMetricData() (result model.CpuMetricFormatData, err error) {
	if c.NewMetric == nil || c.oldSample == nil {
		return result, errors.New("previous sample or current sample is nil. skip")
	}
	timeDelta := c.NewMetric.Totals.Total() - c.oldSample.Totals.Total()
	if timeDelta <= 0 {
		return result, errors.New("previous sample is newer than current sample")
	}

	result.TotalPct = createTotal(c.oldSample.Totals, c.NewMetric.Totals, timeDelta, c.count)
	reportOptMetric := func(current, previous uint64, norm int) float64 {
		if current > 0 {
			return cpuMetricTimeDelta(previous, current, timeDelta, norm)
		}
		return 0
	}

	// /proc/stat metrics
	result.UserPct = reportOptMetric(c.NewMetric.Totals.User, c.oldSample.Totals.User, c.count)
	result.SysPct = reportOptMetric(c.NewMetric.Totals.Sys, c.oldSample.Totals.Sys, c.count)
	result.IdlePct = reportOptMetric(c.NewMetric.Totals.Idle, c.oldSample.Totals.Idle, c.count)
	result.NicePct = reportOptMetric(c.NewMetric.Totals.Nice, c.oldSample.Totals.Nice, c.count)
	result.IrqPct = reportOptMetric(c.NewMetric.Totals.Irq, c.oldSample.Totals.Irq, c.count)
	result.WaitPct = reportOptMetric(c.NewMetric.Totals.Wait, c.oldSample.Totals.Wait, c.count)
	result.SofaIrqPct = reportOptMetric(c.NewMetric.Totals.SoftIrq, c.oldSample.Totals.SoftIrq, c.count)
	result.StplenPct = reportOptMetric(c.NewMetric.Totals.Stolen, c.oldSample.Totals.Stolen, c.count)

	result.User = c.NewMetric.Totals.User
	result.Sys = c.NewMetric.Totals.Sys
	result.Idle = c.NewMetric.Totals.Idle
	result.Nice = c.NewMetric.Totals.Nice
	result.Irq = c.NewMetric.Totals.Irq
	result.Wait = c.NewMetric.Totals.Wait
	result.SoftIrq = c.NewMetric.Totals.SoftIrq
	result.Stolen = c.NewMetric.Totals.Stolen

	return result, nil
}

func createTotal(prev, cur model.CPU, timeDelta uint64, numCPU int) float64 {
	idleTime := cpuMetricTimeDelta(prev.Idle, cur.Idle, timeDelta, numCPU)
	// Subtract wait time from total
	if cur.Wait > 0 {
		idleTime = idleTime + cpuMetricTimeDelta(prev.Wait, cur.Wait, timeDelta, numCPU)
	}
	pct, _ := decimal.NewFromFloat(float64(numCPU) - idleTime).Round(2).Float64()
	return pct
}

// cpuMetricTimeDelta is a helper used by fillTicks to calculate the delta between two CPU tick values
func cpuMetricTimeDelta(prev, current uint64, timeDelta uint64, numCPU int) float64 {
	cpuDelta := int64(current - prev)
	pct := float64(cpuDelta) / float64(timeDelta)
	pct, _ = decimal.NewFromFloat(pct * float64(numCPU)).Round(2).Float64()
	return pct
}

func scanStatFile(scanner *bufio.Scanner) (model.CpuMetricData, error) {
	cpuData, err := statScanner(scanner, parseCPULine)
	if err != nil {
		return model.CpuMetricData{}, fmt.Errorf("error scanning stat file: %w", err)
	}
	return cpuData, nil
}

func parseCPULine(line string) (model.CPU, error) {

	var errs multierror.Errors
	tryParseUint := func(name, field string) (v uint64) {
		u, err := touint64(field)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to parse %v: %s", name, field))
		} else {
			v = u
		}
		return v
	}

	cpuData := model.CPU{}
	fields := strings.Fields(line)

	cpuData.User = tryParseUint("user", fields[1])
	cpuData.Nice = tryParseUint("nice", fields[2])
	cpuData.Sys = tryParseUint("sys", fields[3])
	cpuData.Idle = tryParseUint("idle", fields[4])
	cpuData.Wait = tryParseUint("wait", fields[5])
	cpuData.Irq = tryParseUint("irq", fields[6])
	cpuData.SoftIrq = tryParseUint("softirq", fields[7])
	cpuData.Stolen = tryParseUint("stolen", fields[8])

	return cpuData, errs.Err()
}

func scanCPUInfoFile(scanner *bufio.Scanner) ([]model.CPUInfo, error) {
	return cpuinfoScanner(scanner)
}
