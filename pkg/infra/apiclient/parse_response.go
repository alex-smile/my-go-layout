package apiclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

func parseResponse(resp *http.Response, requestErr error, object interface{}, clientName string) (*apiResponseResult, error) {
	// 检查请求错误
	if requestErr != nil {
		return nil, errors.Wrapf(requestErr, "request api error: [%s]", clientName)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	r := resp.Request

	// 读取响应内容
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "read response body error: [%s] [%s] [%s]", clientName, r.Method, getRequestUrl(r.URL))
	}

	// unmarshal 响应内容
	if object != nil {
		err = json.Unmarshal(respBody, object)
		if err != nil {
			return nil, errors.Wrapf(err, "unmarshal response body error: [%d] [%s] [%s] [%s] %s",
				resp.StatusCode, clientName, r.Method, getRequestUrl(r.URL), string(respBody))
		}
	}

	return &apiResponseResult{
		response:     resp,
		responseBody: respBody,
		request:      r,
		clientName:   clientName,
	}, nil
}

type apiResponseResult struct {
	// 响应信息
	response     *http.Response
	responseBody []byte
	// 请求信息
	request    *http.Request
	clientName string
}

func (r *apiResponseResult) ok() bool {
	return r.response.StatusCode >= 200 && r.response.StatusCode < 300
}

func (r *apiResponseResult) getStatusCode() int {
	return r.response.StatusCode
}

func (r *apiResponseResult) getResponseBody() []byte {
	return r.responseBody
}

func (r *apiResponseResult) getStatusError() error {
	if r.ok() {
		return nil
	}

	return fmt.Errorf("response status code error, [%d] [%s] [%s] [%s] %s",
		r.response.StatusCode, r.clientName, r.request.Method, getRequestUrl(r.request.URL), string(r.responseBody),
	)
}

func (r *apiResponseResult) unmarshalResponseBody(object interface{}) error {
	err := json.Unmarshal(r.responseBody, object)
	if err != nil {
		return errors.Wrapf(err, "unmarshal response body error: [%d] [%s] [%s] [%s] %s",
			r.response.StatusCode, r.clientName, r.request.Method, getRequestUrl(r.request.URL), string(r.responseBody),
		)
	}

	return nil
}

func getRequestUrl(u *url.URL) string {
	return fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.EscapedPath())
}
