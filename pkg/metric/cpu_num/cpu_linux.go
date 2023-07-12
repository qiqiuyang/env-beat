package cpu_num

import (
	"errors"
	"fmt"
	"github.com/qiqiuyang/env-beat/pkg/logutil"
	"github.com/qiqiuyang/env-beat/pkg/metric"
	"github.com/qiqiuyang/env-beat/pkg/metric/enums/cpuCountType"
	"github.com/qiqiuyang/env-beat/pkg/metric/enums/metricType"
	"github.com/qiqiuyang/env-beat/pkg/metric/model"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func init() {
	metric.Register[model.NumMetricParam, model.NumMetricData, model.NumMetricFormatData](metricType.CpuNumMetric, &NumMetricData{})
}

func (n *NumMetricData) FetchMetricData(param model.NumMetricParam) {
	count, exists, err := getCPU(param)
	if err != nil {
		if logutil.GetSugarLogger() != nil {
			logutil.GetSugarLogger().Errorf("cpu核心解析出错err: %s", err.Error())
		}
		n.CpuNum = runtime.NumCPU()
	}
	if !exists {
		if logutil.GetSugarLogger() != nil {
			logutil.GetSugarLogger().Errorf("cpu核心解析失败，将使用runtime.NumCPU返回")
		}
		n.CpuNum = runtime.NumCPU()
		//return runtime.NumCPU()
	}

	n.CpuNum = count
}

func (n *NumMetricData) GetMetricData() (result model.NumMetricData, err error) {
	result.CpuNum = n.CpuNum
	return result, nil
}

func (n *NumMetricData) GetFormatMetricData() (result model.NumMetricFormatData, err error) {
	result.CpuNum = n.CpuNum
	return result, nil
}

func getCPU(param model.NumMetricParam) (int, bool, error) {
	// cpuPath 解析cpu核心的文件位置
	cpuPath := ""

	switch param.CpuSourceType {
	case cpuCountType.CpuCountEnv:
		{
			value, ok := os.LookupEnv(param.CpuSourceValue)
			if !ok {
				return -1, false, nil
			}

			num, err := strconv.Atoi(value)
			if err != nil {
				return -1, false, fmt.Errorf("error Atoi value %s: %w", value, err)
			}

			return num, true, nil
		}
	case cpuCountType.CpuCountRuntime:
		{
			return runtime.NumCPU(), true, nil
		}
	case cpuCountType.CpuCountOnline, cpuCountType.CpuCountPresent:
		{
			cpuPath = param.CpuSourceValue
			break
		}
	default:
		{
			_, isPresent := os.LookupEnv("LINUX_CPU_COUNT_PRESENT")
			cpuPath = "/sys/devices/system/cpu/online"
			if isPresent {
				cpuPath = "/sys/devices/system/cpu/present"
			}
		}
	}

	rawFile, err := os.ReadFile(cpuPath)

	if errors.Is(err, os.ErrNotExist) {
		return -1, false, nil
	}
	if err != nil {
		return -1, false, fmt.Errorf("error reading file %s: %w", cpuPath, err)
	}

	cpuCount, err := parseCPUList(string(rawFile))
	if err != nil {
		return -1, false, fmt.Errorf("error parsing file %s: %w", cpuPath, err)
	}
	return cpuCount, true, nil
}

// parse the weird list files we get from sysfs
func parseCPUList(raw string) (int, error) {
	// 形如：0-15,0-15,1
	listPart := strings.Split(raw, ",")
	count := 0
	for _, value := range listPart {
		// 形如：0-15
		if strings.Contains(value, "-") {
			var first, last int
			_, err := fmt.Sscanf(value, "%d-%d", &first, &last)
			if err != nil {
				continue
			}
			count += (last - first) + 1

		} else {
			// 形如：1
			count++
		}
	}
	return count, nil
}
