package model

type LoadMetricParam struct {
}

type LoadMetricData struct {
	One, Five, Fifteen float64
}

type LoadMetricFormatData struct {
	One           float64
	Five          float64
	Fifteen       float64
	OneMinute     float64
	FiveMinute    float64
	FifteenMinute float64
}
