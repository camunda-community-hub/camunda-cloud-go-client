package client

func NewError(message string) error {
	return &ErrorString{message}
}

type ErrorString struct {
	message string
}

func (e *ErrorString) Error() string {
	return e.message
}
