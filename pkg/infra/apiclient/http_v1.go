package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ContentTypeHeader   = "Content-Type"
	ApplicationJsonMIME = "application/json"
)

type HttpClient struct {
	client   *http.Client
	endpoint string
}

type HttpParam struct {
	Method  string
	BaseUrl string
	Query   map[string]string
	Header  map[string]string
	ReqBody interface{}
}

func NewHttpClient(endpoint string, timeout time.Duration) *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Timeout: timeout,
		},
		endpoint: strings.TrimSuffix(endpoint, "/"),
	}
}

func (client *HttpClient) HttpRequest(ctx context.Context, param *HttpParam) (*httpResponse, error) {
	url := client.handleUrl(param.BaseUrl, param.Query)
	resp, err := client.doHttpRequest(ctx, param.Method, url, param.Header, param.ReqBody)
	return resp, err
}

func (client *HttpClient) handleUrl(baseUrl string, query map[string]string) string {
	values := url.Values{}
	for key, value := range query {
		values.Set(key, value)
	}

	path := baseUrl
	if len(values) > 0 {
		path = fmt.Sprintf("%s?%s", baseUrl, values.Encode())
	}

	if strings.HasPrefix(path, "/") {
		return fmt.Sprintf("%s%s", client.endpoint, path)
	} else {
		return fmt.Sprintf("%s/%s", client.endpoint, path)
	}
}

func (client *HttpClient) doHttpRequest(ctx context.Context, method, url string, header map[string]string, body interface{}) (*httpResponse, error) {
	// TODO: 这里将所有请求都处理为 json，不一定合适
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal http request body error: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create http request error: %v", err)
	}

	for key, value := range header {
		req.Header.Set(key, value)
	}

	resp, err := client.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do http request error: %v", err)
	}

	defer closeHttpResponse(resp)

	return handleHttpResponse(resp)
}

func closeHttpResponse(resp *http.Response) {
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
}

func handleHttpResponse(resp *http.Response) (*httpResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read http response body error: %v", err)
	}

	url := ""
	method := ""
	if resp.Request != nil {
		method = resp.Request.Method
		url = resp.Request.URL.String()
	}

	return &httpResponse{method, url, resp.StatusCode, body}, nil

}

type httpResponse struct {
	reqMethod string
	reqUrl    string
	respCode  int
	respBody  []byte
}

func (r *httpResponse) Method() string {
	return r.reqMethod
}

func (r *httpResponse) Url() string {
	return r.reqUrl
}

func (r *httpResponse) StatusCode() int {
	return r.respCode
}

func (r *httpResponse) Body() []byte {
	return r.respBody
}
