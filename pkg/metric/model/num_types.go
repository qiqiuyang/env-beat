package model

import (
	"github.com/qiqiuyang/env-beat/pkg/metric/enums/cpuCountType"
)

type NumMetricParam struct {
	// CpuSourceType cpu取值来源
	CpuSourceType cpuCountType.CpuCountType
	// CpuSourceValue cpu来源值
	CpuSourceValue string
}

type NumMetricData struct {
	// cpuNum cpu核心数量
	CpuNum int
}

type NumMetricFormatData struct {
	// cpuNum cpu核心数量
	CpuNum int
}
