package api

//api gateway layer, it rec req, extract function name , nvoke runtime, return response , like lambda api layer

import (
	"encoding/json"
	"net/http"
	"strings"

	"serverless-runtime/internal/runtime"
)

type InvokeResponse struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

func InvokeHandler(w http.ResponseWriter, r *http.Request) {

	functionName := strings.TrimPrefix(r.URL.Path, "/invoke/")

	output, err := runtime.ExecuteFunction(functionName)

	response := InvokeResponse{
		Output: output,
	}

	if err != nil {
		response.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}
