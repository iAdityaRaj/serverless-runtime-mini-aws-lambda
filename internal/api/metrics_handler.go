package api

import (
	"encoding/json"
	"net/http"

	"serverless-runtime/internal/metrics"
)

type MetricsResponse struct {
	TotalInvocations uint64 `json:"total_invocations"`

	FailedInvocations uint64 `json:"failed_invocations"`
}

func MetricsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	response := MetricsResponse{
		TotalInvocations: metrics.TotalInvocations,

		FailedInvocations: metrics.FailedInvocations,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).
		Encode(response)
}
