package engine

type Engine struct {
	Policy Policy
}

func (e *Engine) Decide(m MetricsSnapshot) ScaleDecision {

	if m.CPUUtilization > e.Policy.CPUUpperThreshold {
		cpu := e.Policy.MaxVerticalCPU
		return ScaleDecision{
			Mode:   ScaleVertical,
			NewCPU: &cpu,
			Reason: "CPU utilization high, scale vertically",
		}
	}

	if m.QPS > e.Policy.QPSUpperThreshold {
		replicas := int32(3)
		return ScaleDecision{
			Mode:        ScaleHorizontal,
			NewReplicas: &replicas,
			Reason:      "High QPS, scale horizontally",
		}
	}

	return ScaleDecision{
		Mode:   ScaleNone,
		Reason: "Metrics within thresholds",
	}
}
