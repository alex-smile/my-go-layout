package apiclient

import (
	"context"
	"errors"
)

type APIErrorResponseV1 struct {
	Message string       `json:"message"`
	Errors  []APIErrorV1 `json:"errors"`
}

func (r APIErrorResponseV1) String() string {
	if len(r.Message) > 0 {
		return r.Message
	}

	if len(r.Errors) == 0 {
		return ""
	}

	return r.Errors[0].String()
}

type APIErrorV1 struct {
	Source  string `json:"source"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e APIErrorV1) String() string {
	return e.Message
}

func parseAPIErrorResponseV1(ctx context.Context, result *apiResponseResult) error {
	if result.ok() {
		return nil
	}

	var errResponse APIErrorResponseV1
	err := result.unmarshalResponseBody(&errResponse)
	message := errResponse.String()
	// 此处判定 err == nil，且错误消息非空，则认为获取到合适的错误消息
	if err == nil && message != "" {
		return errors.New(message)
	}

	return result.getStatusError()

}
