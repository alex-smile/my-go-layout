package apiclient

type APIError struct {
	statusCode int
	message    string
}

func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{statusCode: statusCode, message: message}
}

func (e *APIError) Error() string {
	return e.message
}

func (e *APIError) GetStatusCode() int {
	return e.statusCode
}
