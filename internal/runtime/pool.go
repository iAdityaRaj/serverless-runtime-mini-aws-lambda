package runtime

import (
	"log"
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

	for i := 0; i < workerCount; i++ { // worker waits forever for jobs, persistent warm worker

		go p.worker(i) // creates lightweight thread
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

		log.Printf(
			"Worker %d executing %s",
			id,
			job.FunctionName,
		)

		output, err := p.Runtime.ExecuteFunction(
			job.FunctionName,
			job.Payload,
		)

		job.ResultChan <- InvocationResult{
			Output: output,
			Error:  err,
		}
	}
}
