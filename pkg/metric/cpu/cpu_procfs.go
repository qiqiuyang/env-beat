package cpu

import (
	"bufio"
	"fmt"
	"github.com/qiqiuyang/env-beat/pkg/metric/model"
	"os"
	"strconv"
	"strings"
)

// Get returns a metrics object for CPU data
func Get() (model.CpuMetricData, error) {
	path := "/proc/stat"
	fd, err := os.Open(path)
	defer func() {
		_ = fd.Close()
	}()
	if err != nil {
		return model.CpuMetricData{}, fmt.Errorf("error opening file %s: %w", path, err)
	}

	metrics, err := scanStatFile(bufio.NewScanner(fd))
	if err != nil {
		return model.CpuMetricData{}, fmt.Errorf("scanning stat file: %w", err)
	}

	cpuInfoPath := "/proc/cpuinfo"
	cpuInfoFd, err := os.Open(cpuInfoPath)
	if err != nil {
		return model.CpuMetricData{}, fmt.Errorf("opening '%s': %w", cpuInfoPath, err)
	}
	defer cpuInfoFd.Close()

	cpuInfo, err := scanCPUInfoFile(bufio.NewScanner(cpuInfoFd))
	metrics.CPUInfo = cpuInfo

	return metrics, err
}

func cpuinfoScanner(scanner *bufio.Scanner) ([]model.CPUInfo, error) {
	cpuInfos := []model.CPUInfo{}
	current := model.CPUInfo{}
	// On my tests the order the cores appear on /proc/cpuinfo
	// is the same as on /proc/stats, this means it matches our
	// current 'system.core.id' metric. This information
	// is also the same as the 'processor' line on /proc/cpuinfo.
	coreID := 0
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, ":")
		if len(split) != 2 {
			// A blank line its a separation between CPUs
			// even the last CPU contains one blank line at the end
			cpuInfos = append(cpuInfos, current)
			current = model.CPUInfo{}
			coreID++

			continue
		}

		k, v := split[0], split[1]
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		switch k {
		case "model":
			current.ModelNumber = v
		case "model name":
			current.ModelName = v
		case "physical id":
			id, err := strconv.Atoi(v)
			if err != nil {
				return []model.CPUInfo{}, fmt.Errorf("parsing physical ID: %w", err)
			}
			current.PhysicalID = id
		case "core id":
			id, err := strconv.Atoi(v)
			if err != nil {
				return []model.CPUInfo{}, fmt.Errorf("parsing core ID: %w", err)
			}
			current.CoreID = id
		case "cpu MHz":
			mhz, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return []model.CPUInfo{}, fmt.Errorf("parsing CPU %d Mhz: %w", coreID, err)
			}
			current.Mhz = mhz
		}
	}

	return cpuInfos, nil
}

// statScanner iterates through a /proc/stat entry, reading both the global lines and per-CPU lines, each time calling lineReader, which implements the OS-specific code for parsing individual lines
func statScanner(scanner *bufio.Scanner, lineReader func(string) (model.CPU, error)) (model.CpuMetricData, error) {
	cpuData := model.CpuMetricData{}
	var err error

	for scanner.Scan() {
		text := scanner.Text()
		// Check to see if this is the global CPU line
		if isCPUGlobalLine(text) {
			cpuData.Totals, err = lineReader(text)
			if err != nil {
				return model.CpuMetricData{}, fmt.Errorf("error parsing global CPU line: %w", err)
			}
		}
		if isCPULine(text) {
			perCPU, err := lineReader(text)
			if err != nil {
				return model.CpuMetricData{}, fmt.Errorf("error parsing CPU line: %w", err)
			}
			cpuData.List = append(cpuData.List, perCPU)

		}
	}
	return cpuData, nil
}

func isCPUGlobalLine(line string) bool {
	if len(line) > 4 && line[0:4] == "cpu " {
		return true
	}
	return false
}

func isCPULine(line string) bool {
	if len(line) > 3 && line[0:3] == "cpu" && line[3] != ' ' {
		return true
	}
	return false
}

func touint64(val string) (uint64, error) {
	return strconv.ParseUint(val, 10, 64)
}
