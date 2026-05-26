package runtime

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
)

type DockerRuntime struct {
}

func (d *DockerRuntime) ExecuteFunction(
	name string,
	payload []byte,
) (string, error) {

	imageName := fmt.Sprintf(
		"%s-function",
		name,
	)

	log.Printf(
		"Executing containerized function: %s",
		name,
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		ExecutionTimeout,
	)

	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"docker",
		"run",
		"-i",   // interactive mode
		"--rm", // remove the container after execution

		"--memory",
		MemoryLimit,

		"--cpus", // cgroups
		CPULimit,

		"--pids-limit",
		PidsLimit,

		imageName,
	)

	cmd.Stdin = bytes.NewReader(payload)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	log.Println("Starting docker container")
	err := cmd.Run()
	log.Println("Docker container finished")

	if ctx.Err() ==
		context.DeadlineExceeded {
		log.Println("Container execution timed out")
		return "",
			fmt.Errorf(
				"container execution timed out",
			)
	}

	if err != nil {

		return "",
			fmt.Errorf(
				stderr.String(),
			)
	}

	return out.String(), nil
}
