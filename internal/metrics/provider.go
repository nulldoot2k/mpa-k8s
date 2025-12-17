package engine

type MetricsSnapshot struct {
	CPUUtilization    float64 // %
	MemoryUtilization float64 // %
	QPS               float64
	LatencyP99Ms      float64
}
