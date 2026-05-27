package api

import (
	"encoding/json"
	"net/http"

	"serverless-runtime/internal/metrics"
)

type MetricsResponse struct {
	TotalInvocations uint64 `json:"total_invocations"`

	FailedInvocations uint64 `json:"failed_invocations"`

	ActiveWorkers int64 `json:"active_workers"`

	QueueDepth int64 `json:"queue_depth"`

	AverageExecutionMs float64 `json:"average_execution_ms"`

	AverageQueueWaitMs float64 `json:"average_queue_wait_ms"`
}

func MetricsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	avgExecution := 0.0
	avgQueueWait := 0.0

	if metrics.TotalInvocations > 0 {

		avgExecution =
			float64(
				metrics.TotalExecutionTimeNs,
			) /
				float64(
					metrics.TotalInvocations,
				) /
				1e6

		avgQueueWait =
			float64(
				metrics.TotalQueueWaitTimeNs,
			) /
				float64(
					metrics.TotalInvocations,
				) /
				1e6
	}

	response := MetricsResponse{
		TotalInvocations: metrics.TotalInvocations,

		FailedInvocations: metrics.FailedInvocations,

		ActiveWorkers: metrics.ActiveWorkers,

		QueueDepth: metrics.QueueDepth,

		AverageExecutionMs: avgExecution,

		AverageQueueWaitMs: avgQueueWait,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).
		Encode(response)
}
