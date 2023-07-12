package memory

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/qiqiuyang/env-beat/pkg/logutil"
	"github.com/qiqiuyang/env-beat/pkg/metric"
	"github.com/qiqiuyang/env-beat/pkg/metric/enums/metricType"
	"github.com/qiqiuyang/env-beat/pkg/metric/model"
	"github.com/shopspring/decimal"
	"io"
	"os"
	"strconv"
	"strings"
)

func init() {
	metric.Register[model.MemoryMetricParam, model.MemoryMetricData, model.MemoryMetricFormatData](metricType.MemoryMetric, &MemoryMetricData{})
}

func (n *MemoryMetricData) FetchMetricData(param model.MemoryMetricParam) {
	base, err := get()
	if err != nil {
		if logutil.GetSugarLogger() != nil {
			logutil.GetSugarLogger().Errorf("Fetch failed, err: %s", err.Error())
		}
		n.errMsg = err
		return
	}
	n.Total = base.Total
	n.Used = base.Used
	n.Free = base.Free
	n.Cached = base.Cached
	n.Actual = base.Actual
	n.Swap = base.Swap
}

func (n *MemoryMetricData) GetMetricData() (result model.MemoryMetricData, err error) {
	result.Total = n.Total
	result.Used = n.Used
	result.Free = n.Free
	result.Cached = n.Cached
	result.Actual = n.Actual
	result.Swap = n.Swap
	return result, n.errMsg
}

func (n *MemoryMetricData) GetFormatMetricData() (result model.MemoryMetricFormatData, err error) {
	result.Total = n.Total
	result.Used = n.Used
	result.Free = n.Free
	result.Cached = n.Cached
	result.Actual = n.Actual
	result.Swap = n.Swap
	return result, n.errMsg
}

func (base *MemoryMetricData) fillPercentages() {
	// Add percentages
	// In theory, `Used` and `Total` are available everywhere, so assume values are good.
	if base.Total > 0 {
		percUsed := float64(base.Used.Bytes) / float64(base.Total)
		base.Used.Pct, _ = decimal.NewFromFloat(percUsed).Round(2).Float64()

		actualPercUsed := float64(base.Actual.Used.Bytes) / float64(base.Total)
		base.Actual.Used.Pct, _ = decimal.NewFromFloat(actualPercUsed).Round(2).Float64()
	}

	if base.Swap.Total > 0 {
		perc := float64(base.Swap.Used.Bytes) / float64(base.Swap.Total)
		base.Swap.Used.Pct, _ = decimal.NewFromFloat(perc).Round(2).Float64()
	}
}

// get is the linux implementation for fetching Memory data
func get() (model.MemoryMetricData, error) {
	table, err := ParseMeminfo()
	if err != nil {
		return model.MemoryMetricData{}, fmt.Errorf("error fetching meminfo: %w", err)
	}

	memData := MemoryMetricData{}

	var free, cached uint64
	var ok bool
	if total, ok := table["MemTotal"]; ok {
		memData.Total = total
	}
	if free, ok = table["MemFree"]; ok {
		memData.Free = free
	}
	if cached, ok = table["Cached"]; ok {
		memData.Cached = cached
	}

	// overlook parsing issues here
	// On the very small chance some of these don't exist,
	// It's not the end of the world
	buffers := table["Buffers"]

	if memAvail, ok := table["MemAvailable"]; ok {
		// MemAvailable is in /proc/meminfo (kernel 3.14+)
		memData.Actual.Free = memAvail
	} else {
		// in the future we may want to find another way to do this.
		// "MemAvailable" and other more derivied metrics
		// Are very relative, and can be unhelpful in cerntain workloads
		// We may want to find a way to more clearly express to users
		// where a certain value is coming from and what it represents

		// The use of `cached` here is particularly concerning,
		// as under certain intense DB server workloads, the cached memory can be quite large
		// and give the impression that we've passed memory usage watermark
		memData.Actual.Free = free + buffers + cached
	}

	memData.Used.Bytes = memData.Total - memData.Free
	memData.Actual.Used.Bytes = memData.Total - memData.Actual.Free

	// Populate swap data
	swapTotal, okST := table["SwapTotal"]
	if okST {
		memData.Swap.Total = swapTotal
	}
	swapFree, okSF := table["SwapFree"]
	if okSF {
		memData.Swap.Free = swapFree
	}

	if okSF && okST {
		memData.Swap.Used.Bytes = swapTotal - swapFree
	}
	memData.fillPercentages()
	return memData.MemoryMetricData, nil

}

// ParseMeminfo parses the contents of /proc/meminfo into a hashmap
func ParseMeminfo() (map[string]uint64, error) {
	table := map[string]uint64{}

	meminfoPath := "/proc/meminfo"
	err := readFile(meminfoPath, func(line string) bool {
		fields := strings.Split(line, ":")

		if len(fields) != 2 {
			return true // skip on errors
		}

		valueUnit := strings.Fields(fields[1])
		value, err := strconv.ParseUint(valueUnit[0], 10, 64)
		if err != nil {
			return true // skip on errors
		}

		if len(valueUnit) > 1 && valueUnit[1] == "kB" {
			value *= 1024
		}
		table[fields[0]] = value

		return true
	})
	return table, err
}

func readFile(file string, handler func(string) bool) error {
	contents, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", file, err)
	}

	reader := bufio.NewReader(bytes.NewBuffer(contents))

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if !handler(string(line)) {
			break
		}
	}

	return nil
}
