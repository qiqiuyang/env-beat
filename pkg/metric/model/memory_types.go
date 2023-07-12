package model

type MemoryMetricParam struct {
}

type MemoryMetricData struct {
	Total uint64       `struct:"total,omitempty"`
	Used  UsedMemStats `struct:"used,omitempty"`

	Free   uint64 `struct:"free,omitempty"`
	Cached uint64 `struct:"cached,omitempty"`
	// "Actual" values are, technically, a linux-only concept
	// For better or worse we've expanded it to include "derived"
	// Memory values on other platforms, which we should
	// probably keep for the sake of backwards compatibility
	// However, because the derived value varies from platform to platform,
	// We may want to more precisely document what these mean.
	Actual ActualMemoryMetrics `struct:"actual,omitempty"`

	// Swap metrics
	Swap SwapMetrics `struct:"swap,omitempty"`
}

type MemoryMetricFormatData struct {
	Total uint64       `struct:"total,omitempty"`
	Used  UsedMemStats `struct:"used,omitempty"`

	Free   uint64 `struct:"free,omitempty"`
	Cached uint64 `struct:"cached,omitempty"`
	// "Actual" values are, technically, a linux-only concept
	// For better or worse we've expanded it to include "derived"
	// Memory values on other platforms, which we should
	// probably keep for the sake of backwards compatibility
	// However, because the derived value varies from platform to platform,
	// We may want to more precisely document what these mean.
	Actual ActualMemoryMetrics `struct:"actual,omitempty"`

	// Swap metrics
	Swap SwapMetrics `struct:"swap,omitempty"`
}

// UsedMemStats wraps used.* memory metrics
type UsedMemStats struct {
	Pct   float64 `struct:"pct,omitempty"`
	Bytes uint64  `struct:"bytes,omitempty"`
}

// ActualMemoryMetrics wraps the actual.* memory metrics
type ActualMemoryMetrics struct {
	Free uint64       `struct:"free,omitempty"`
	Used UsedMemStats `struct:"used,omitempty"`
}

// SwapMetrics wraps swap.* memory metrics
type SwapMetrics struct {
	Total uint64       `struct:"total,omitempty"`
	Used  UsedMemStats `struct:"used,omitempty"`
	Free  uint64       `struct:"free,omitempty"`
}
