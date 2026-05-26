package runtime

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
)

type LocalRuntime struct {
}

func (l *LocalRuntime) ExecuteFunction(
	name string,
	payload []byte,
) (string, error) {

	path := fmt.Sprintf(
		"./build/%s",
		name,
	)

	log.Printf(
		"Executing compiled function: %s",
		name,
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		ExecutionTimeout,
	)

	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		path,
	)

	cmd.Stdin = bytes.NewReader(payload)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf(
			"function execution timed out",
		)
	}

	if err != nil {
		return "", fmt.Errorf(
			stderr.String(),
		)
	}

	return out.String(), nil
}
