package main

import (
	"fmt"
	"github.com/qiqiuyang/env-beat/beanFactory"
	"github.com/qiqiuyang/env-beat/pkg/metric/enums/cpuCountType"
	"github.com/qiqiuyang/env-beat/pkg/metric/model"
	"time"
)

func main() {
	beanFactory.CpuNumMetric.FetchMetricData(model.NumMetricParam{
		CpuSourceType:  cpuCountType.CpuCountDefault,
		CpuSourceValue: "",
	})

	result, _ := beanFactory.CpuNumMetric.GetMetricData()
	fmt.Println("Hello, result !", result)

	beanFactory.CpuMetric.FetchMetricData(model.CpuMetricParam{})
	result1, _ := beanFactory.CpuMetric.GetMetricData()
	fmt.Println("Hello, result1!", result1)

	result2, _ := beanFactory.CpuMetric.GetFormatMetricData()
	fmt.Println("Hello, result2!", result2)

	time.Sleep(1000000)
	beanFactory.CpuMetric.FetchMetricData(model.CpuMetricParam{})
	result3, _ := beanFactory.CpuMetric.GetFormatMetricData()
	fmt.Println("Hello, result3!", result3)

	beanFactory.MemoryMetric.FetchMetricData(model.MemoryMetricParam{})
	result4, _ := beanFactory.MemoryMetric.GetFormatMetricData()
	fmt.Println("Hello, result4!", result4)

	beanFactory.LoadMetric.FetchMetricData(model.LoadMetricParam{})
	result5, _ := beanFactory.LoadMetric.GetFormatMetricData()
	fmt.Println("Hello, result5!", result5)

}
