package runtime

//runtime execution engine

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"
)

const ExecutionTimeout = 5 * time.Second // set timeout limit for function

func ExecuteFunction(name string) (string, error) {

	path := fmt.Sprintf("./functions/%s/main.go", name)

	log.Printf("Executing function: %s", name)

	ctx, cancel := context.WithTimeout( //create cancellable execution context
		context.Background(),
		ExecutionTimeout,
	)

	defer cancel() //defer means , run this function when current fun exists

	cmd := exec.CommandContext(
		ctx,
		"go",
		"run",
		path,
	)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("function execution timed out")
	}

	if err != nil {
		return "", fmt.Errorf(stderr.String())
	}

	log.Printf("Function %s executed successfully", name)

	return out.String(), nil
}
