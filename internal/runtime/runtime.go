package runtime

type Runtime interface {
	ExecuteFunction(
		name string,
		payload []byte,
	) (string, error)
}
