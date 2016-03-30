package fakes

type FakeCloudError struct {
	t       string
	message string
}

func NewFakeCloudError(t, message string) FakeCloudError {
	return FakeCloudError{t: t, message: message}
}

func (e FakeCloudError) Type() string  { return e.t }
func (e FakeCloudError) Error() string { return e.message }

type FakeRetryableError struct {
	canRetry bool
	message  string
}

func NewFakeRetryableError(message string, canRetry bool) FakeRetryableError {
	return FakeRetryableError{message: message, canRetry: canRetry}
}

func (e FakeRetryableError) Error() string  { return e.message }
func (e FakeRetryableError) CanRetry() bool { return e.canRetry }
