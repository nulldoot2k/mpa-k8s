package engine

import "testing"

func TestVerticalScalingDecision(t *testing.T) {
	engine := Engine{
		Policy: Policy{
			CPUUpperThreshold: 80,
			MaxVerticalCPU:    2000,
		},
	}

	metrics := MetricsSnapshot{
		CPUUtilization: 90,
	}

	decision := engine.Decide(metrics)

	if decision.Mode != ScaleVertical {
		t.Fatalf("expected ScaleVertical, got %s", decision.Mode)
	}

	if decision.NewCPU == nil {
		t.Fatalf("expected NewCPU to be set")
	}
}
