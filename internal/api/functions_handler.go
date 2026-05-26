package api

import (
	"encoding/json"
	"net/http"

	"serverless-runtime/internal/registry"
)

func ListFunctionsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	functions :=
		registry.GlobalRegistry.List()

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).
		Encode(functions)
}
