package dispatcher

type Dispatcher interface {
	// Dispatch interprets request bytes, executes request, captures response and return response bytes.
	// It panics if built-in errors fail to serialize.
	Dispatch([]byte) []byte
}
