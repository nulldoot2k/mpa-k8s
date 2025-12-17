package engine

type ScaleMode string

const (
	ScaleNone       ScaleMode = "None"
	ScaleVertical   ScaleMode = "Vertical"
	ScaleHorizontal ScaleMode = "Horizontal"
)

type ScaleDecision struct {
	Mode        ScaleMode
	NewCPU      *int64 // millicores
	NewMemory   *int64 // MiB
	NewReplicas *int32
	Reason      string
}
