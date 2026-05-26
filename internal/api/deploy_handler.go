package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"serverless-runtime/internal/registry"
)

func DeployHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req DeployRequest

	err := json.NewDecoder(
		r.Body,
	).Decode(&req)

	if err != nil {

		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)

		return
	}

	functionDir := fmt.Sprintf(
		"./functions/%s",
		req.Name,
	)

	err = os.MkdirAll(
		functionDir,
		0755,
	)

	if err != nil {

		http.Error(
			w,
			"failed to create function directory",
			http.StatusInternalServerError,
		)

		return
	}

	mainFile := fmt.Sprintf(
		"%s/main.go",
		functionDir,
	)

	err = os.WriteFile(
		mainFile,
		[]byte(req.Code),
		0644,
	)

	if err != nil {

		http.Error(
			w,
			"failed to write source code",
			http.StatusInternalServerError,
		)

		return
	}

	dockerfile := `
FROM golang:1.24-alpine

WORKDIR /app

COPY main.go .

RUN go build -o function main.go

CMD ["./function"]
`

	dockerfilePath := fmt.Sprintf(
		"%s/Dockerfile",
		functionDir,
	)

	err = os.WriteFile(
		dockerfilePath,
		[]byte(dockerfile),
		0644,
	)

	if err != nil {

		http.Error(
			w,
			"failed to create Dockerfile",
			http.StatusInternalServerError,
		)

		return
	}

	imageName := fmt.Sprintf(
		"%s-function",
		req.Name,
	)

	buildCmd := exec.Command(
		"docker",
		"build",
		"-t",
		imageName,
		functionDir,
	)

	output, err := buildCmd.CombinedOutput()

	if err != nil {

		http.Error(
			w,
			fmt.Sprintf(
				"docker build failed:\n%s",
				string(output),
			),
			http.StatusInternalServerError,
		)

		return
	}

	registry.GlobalRegistry.Register(
		registry.Function{
			Name: req.Name,

			Image: imageName,

			Runtime: "docker",

			CreatedAt: time.Now(),
		},
	)

	w.WriteHeader(http.StatusCreated)

	w.Write(
		[]byte(
			"function deployed successfully",
		),
	)
}
