package metrics

import (
	"sync/atomic"
)

var TotalInvocations uint64

var FailedInvocations uint64

var ActiveWorkers int64

var QueueDepth int64

var TotalExecutionTimeNs uint64

var TotalQueueWaitTimeNs uint64

func IncrementInvocations() {
	atomic.AddUint64(&TotalInvocations, 1)
} //concurrent Multiple workers increment metrics simultaneously.Without atomic ops:race conditions occur

func IncrementFailures() {

	atomic.AddUint64(&FailedInvocations, 1)
}

func IncrementActiveWorkers() {

	atomic.AddInt64(&ActiveWorkers, 1)
}

func DecrementActiveWorkers() {

	atomic.AddInt64(&ActiveWorkers, -1)
}

func SetQueueDepth(depth int) {
	atomic.StoreInt64(&QueueDepth, int64(depth))
}

func AddExecutionTime(ns uint64) {

	atomic.AddUint64(&TotalExecutionTimeNs, ns)
}

func AddQueueWaitTime(ns uint64) {

	atomic.AddUint64(&TotalQueueWaitTimeNs, ns)
}
