package errorx

type CodeError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	err     error
}

func (e CodeError) GetCode() int {
	return e.Code
}

func (e CodeError) GetMessage() string {
	return e.Message
}

func (e CodeError) Error() string {
	return e.Message
}

func (e CodeError) Unwrap() error {
	return e.err
}

func NewCodeError(code int, err error) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(CodeError); ok {
		return err
	}

	return CodeError{
		Code:    code,
		Message: err.Error(),
		err:     err,
	}
}
