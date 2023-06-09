package cerrors

type CodedError struct {
	Code int
	Err  error
}

func New(code int, err error) *CodedError {
	return &CodedError{Code: code, Err: err}
}

func (e *CodedError) Error() string {
	return e.Err.Error()
}
