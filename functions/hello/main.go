package main

//represents deployed user code
import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Request struct {
	Name string `json:"name"`
}

func main() {

	input, err := io.ReadAll(os.Stdin) // subprocess input stream

	if err != nil {
		fmt.Println("failed to read stdin")
		return
	}

	var req Request

	err = json.Unmarshal(input, &req)

	if err != nil {
		fmt.Println("invalid json payload")
		return
	}

	fmt.Printf(
		"Hello %s from serverless runtime!\n",
		req.Name,
	)
}
