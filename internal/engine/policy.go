package engine

type Policy struct {
	CPUUpperThreshold    float64
	MemoryUpperThreshold float64
	MaxVerticalCPU       int64
	MaxVerticalMemory    int64
	QPSUpperThreshold    float64
}
