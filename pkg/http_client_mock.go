package httpclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Mock struct {
	RequestURL     string
	RequestHeader  http.Header
	RequestBody    string
	RequestMethod  string
	ResponseHeader http.Header
	ResponseBody   string
	ResponseStatus int
	Error          error
}

func (me *Mock) Do(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		me.RequestBody = string(body)
	}

	me.RequestMethod = req.Method
	me.RequestURL = req.URL.String()
	me.RequestHeader = req.Header

	if me.Error != nil {
		return nil, me.Error
	}

	response := ioutil.NopCloser(bytes.NewReader([]byte(me.ResponseBody)))

	return &http.Response{
		StatusCode: me.ResponseStatus,
		Body:       response,
		Header:     me.ResponseHeader,
	}, nil
}

func (me *Mock) Status(statusCode int) *Mock {
	me.ResponseStatus = statusCode
	return me
}

func (me *Mock) Body(body string) *Mock {
	me.ResponseBody = body
	return me
}

func (me *Mock) Err(err error) *Mock {
	me.Error = err
	return me
}

type MultMock struct {
	Mocks map[string]map[string]*Mock
}

func NewHTTPMultMock() *MultMock {
	return &MultMock{make(map[string]map[string]*Mock)}
}

func (me *MultMock) createMock(method, URL string) *Mock {
	methodMocks, ok := me.Mocks[method]

	if !ok {
		methodMocks = make(map[string]*Mock)
		me.Mocks[method] = methodMocks
	}
	mock := &Mock{}
	methodMocks[URL] = mock
	return mock
}

func (me *MultMock) Get(URL string) *Mock {
	return me.createMock("GET", URL)
}

func (me *MultMock) Put(URL string) *Mock {
	return me.createMock("PUT", URL)
}

func (me *MultMock) Post(URL string) *Mock {
	return me.createMock("POST", URL)
}

func (me *MultMock) Do(req *http.Request) (*http.Response, error) {
	method := req.Method
	URL := req.URL.String()

	methodMocks, ok := me.Mocks[method]
	if !ok {
		return nil, fmt.Errorf("No mock for [%s][%s]", method, URL)
	}
	pathMock, ok := methodMocks[URL]
	if !ok {
		return nil, fmt.Errorf("No mock for method [%s][%s]", method, URL)
	}

	return pathMock.Do(req)
}
