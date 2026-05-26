package runtime

type InvocationJob struct {
	FunctionName string // tells what to do
	Payload      []byte // holds data

	ResultChan chan InvocationResult // asynchronus response channel, Worker sends result back through channel
}

type InvocationResult struct {
	Output string
	Error  error
}
