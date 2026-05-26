package api

//api gateway layer, it rec req, extract function name , nvoke runtime, return response , like lambda api layer

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"serverless-runtime/internal/runtime"

	"serverless-runtime/internal/metrics"
	"serverless-runtime/internal/registry"
)

type InvokeResponse struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

var runtimeEngine runtime.Runtime = &runtime.DockerRuntime{}

var workerPool = runtime.NewWorkerPool(
	runtimeEngine,
	100,
)

func init() {

	workerPool.Start(3)
}

func InvokeHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	functionName := strings.TrimPrefix(
		r.URL.Path,
		"/invoke/",
	)
	metrics.IncrementInvocations()

	registry.GlobalRegistry.
		IncrementInvocation(
			functionName,
		)

	body, err := io.ReadAll(r.Body)

	if err != nil {

		http.Error(
			w,
			"failed to read request body",
			http.StatusBadRequest,
		)

		return
	}

	resultChan := make(
		chan runtime.InvocationResult, //This channel is the exclusive link between this specific HTTP request and the worker that picks up the job.
	)

	job := runtime.InvocationJob{
		FunctionName: functionName,
		Payload:      body,
		ResultChan:   resultChan,
	}

	workerPool.JobQueue <- job

	result := <-resultChan // The HTTP Handler goroutine blocks here. It stops execution and waits, parked by go runtime

	response := InvokeResponse{
		Output: result.Output,
	}

	if result.Error != nil {
		metrics.IncrementFailures()
		response.Error =
			result.Error.Error()

		w.WriteHeader(
			http.StatusInternalServerError,
		)
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(response)
}
