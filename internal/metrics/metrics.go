package metrics

import (
	"sync/atomic"
)

var TotalInvocations uint64

var FailedInvocations uint64

func IncrementInvocations() {

	atomic.AddUint64( //concurrent Multiple workers increment metrics simultaneously.Without atomic ops:race conditions occur
		&TotalInvocations,
		1,
	)
}

func IncrementFailures() {

	atomic.AddUint64(
		&FailedInvocations,
		1,
	)
}
