package runtime

import "time"

const (
	ExecutionTimeout = 5 * time.Second

	MemoryLimit = "128m"

	CPULimit = "0.5"

	PidsLimit = "64"
)
