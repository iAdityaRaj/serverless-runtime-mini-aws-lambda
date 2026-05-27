package runtime

import (
	"log"
	"time"

	"serverless-runtime/internal/metrics"
)

type WorkerPool struct {
	Runtime Runtime

	JobQueue chan InvocationJob
}

func NewWorkerPool(
	r Runtime,
	queueSize int,
) *WorkerPool {

	return &WorkerPool{
		Runtime: r,

		JobQueue: make(
			chan InvocationJob,
			queueSize,
		),
	}
}

func (p *WorkerPool) Start(
	workerCount int,
) {

	for i := 0; i < workerCount; i++ {

		go p.worker(i)
	}
}

func (p *WorkerPool) worker(
	id int,
) {

	log.Printf(
		"Worker %d started",
		id,
	)

	for job := range p.JobQueue {

		p.executeJob(
			id,
			job,
		)
	}
}

func (p *WorkerPool) executeJob(
	id int,
	job InvocationJob,
) {

	metrics.SetQueueDepth(
		len(p.JobQueue),
	)

	metrics.IncrementActiveWorkers()

	defer metrics.DecrementActiveWorkers()

	waitTime := time.Since(
		job.QueuedAt,
	)

	metrics.AddQueueWaitTime(
		uint64(waitTime.Nanoseconds()),
	)

	start := time.Now()

	log.Printf(
		"Worker %d executing %s",
		id,
		job.FunctionName,
	)

	output, err := p.Runtime.ExecuteFunction(
		job.FunctionName,
		job.Payload,
	)

	duration := time.Since(start)

	metrics.AddExecutionTime(
		uint64(duration.Nanoseconds()),
	)

	job.ResultChan <- InvocationResult{
		Output: output,
		Error:  err,
	}
}
