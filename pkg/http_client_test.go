package httpclient_test

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	httpclient "gitlab.neoway.com.br/sebrae-sp-ddm-data-collector/pkg/http_client"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

var (
	httpMockUpExpected    = httpclient.NewHTTPMultMock()
	url                   = "https://www.url.com"
	mockedRequestExpected = &httpclient.Mock{
		RequestURL:     "",
		RequestHeader:  http.Header(nil),
		RequestBody:    "",
		RequestMethod:  "",
		ResponseHeader: http.Header(nil),
		ResponseBody:   "",
		ResponseStatus: 0, Error: error(nil),
	}
	errFoo                        = errors.New("internal server error")
	mockedRequestExpectedWithBody = &httpclient.Mock{
		RequestURL:     "",
		RequestHeader:  http.Header(nil),
		RequestBody:    "",
		RequestMethod:  "",
		ResponseHeader: http.Header(nil),
		ResponseBody:   `{"test":"body"}`,
		ResponseStatus: 500,
		Error:          errFoo,
	}
	mock = &httpclient.Mock{
		RequestURL:     url,
		RequestHeader:  http.Header(nil),
		RequestBody:    `{"test":"request body"}`,
		RequestMethod:  "GET",
		ResponseHeader: http.Header(nil),
		ResponseBody:   `{"test":"response body"}`,
		ResponseStatus: 200,
	}
	mockWithError = &httpclient.Mock{
		RequestURL:     url,
		RequestHeader:  http.Header(nil),
		RequestBody:    `{"test":"request body"}`,
		RequestMethod:  "GET",
		ResponseHeader: http.Header(nil),
		ResponseBody:   `{"test":"response body"}`,
		ResponseStatus: 500,
		Error:          errFoo,
	}
	mockWithNoBody = &httpclient.Mock{
		RequestURL:     url,
		RequestHeader:  http.Header(nil),
		RequestMethod:  "GET",
		ResponseHeader: http.Header(nil),
		ResponseStatus: 500,
		Error:          errFoo,
	}
	httpResponseExpected = &http.Response{
		Status:           "",
		StatusCode:       200,
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           http.Header(nil),
		Body:             ioutil.NopCloser(bytes.NewReader([]byte(`{"test":"response body"}`))),
		ContentLength:    0,
		TransferEncoding: []string(nil),
		Close:            false,
		Uncompressed:     false,
		Trailer:          http.Header(nil),
		Request:          (*http.Request)(nil),
		TLS:              (*tls.ConnectionState)(nil),
	}
)

func TestNewHTTPClient(t *testing.T) {
	client := httpclient.NewHTTPClient(time.Minute)
	assert.NotNil(t, client)
}
func TestHTTPClient_Mock(t *testing.T) {
	t.Run("HTTP Client Mock execute Do command returns a response, and test Body e Status command", func(t *testing.T) {
		// Create Test Request
		req, err := http.NewRequest(http.MethodGet, "http://test/v1", strings.NewReader("test payload"))
		assert.NoError(t, err)
		assert.NotNil(t, req)

		// Create Mock
		httpCliMock := &httpclient.Mock{
			ResponseBody:   "return done",
			ResponseStatus: 200,
		}

		response, err := httpCliMock.Do(req)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, httpCliMock.ResponseStatus, response.StatusCode)

		bodyCliMock := httpCliMock.Body(httpCliMock.ResponseBody)
		assert.NotNil(t, bodyCliMock)
		assert.Equal(t, bodyCliMock.ResponseBody, httpCliMock.ResponseBody)

		statusCliMock := httpCliMock.Status(httpCliMock.ResponseStatus)
		assert.NotNil(t, statusCliMock)
		assert.Equal(t, statusCliMock.ResponseStatus, httpCliMock.ResponseStatus)
	})

	t.Run("HTTP Client Mock execute Do command finds struct error. And Test Err command", func(t *testing.T) {
		expectedError := "Failed to authenticate"

		// Create Test Request
		req, err := http.NewRequest(http.MethodGet, "http://test/v1", strings.NewReader("test payload"))
		assert.NoError(t, err)
		assert.NotNil(t, req)

		// Create Mock
		httpCliMock := &httpclient.Mock{
			Error: fmt.Errorf("%v", expectedError),
		}

		response, err := httpCliMock.Do(req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err.Error())

		errorCliMock := httpCliMock.Err(err)
		assert.NotNil(t, errorCliMock)
		assert.Equal(t, errorCliMock.Error, httpCliMock.Error)
	})
}

func TestNewHTTPMultMock(t *testing.T) {
	httpMockUp := httpclient.NewHTTPMultMock()
	assert.Equal(t, httpMockUpExpected, httpMockUp)
}

func TestRequestVerbs(t *testing.T) {
	t.Run("GET", func(t *testing.T) {
		httpMockUp := httpclient.NewHTTPMultMock()
		mockedRequest := httpMockUp.Get(url)
		assert.Equal(t, mockedRequestExpected, mockedRequest)
	})
	t.Run("POST", func(t *testing.T) {
		httpMockUp := httpclient.NewHTTPMultMock()
		mockedRequest := httpMockUp.Post(url)
		assert.Equal(t, mockedRequestExpected, mockedRequest)
	})
	t.Run("PUT", func(t *testing.T) {
		httpMockUp := httpclient.NewHTTPMultMock()
		mockedRequest := httpMockUp.Put(url)
		assert.Equal(t, mockedRequestExpected, mockedRequest)
	})
}

func TestBodyRequest(t *testing.T) {
	httpMockUp := httpclient.NewHTTPMultMock()
	mockedRequest := httpMockUp.Post(url)
	mockedRequest.Body(`{"test":"body"}`)
	mockedRequest.Err(errFoo)
	mockedRequest.Status(500)
	assert.Equal(t, mockedRequestExpectedWithBody, mockedRequest)
}
func TestDoRequest(t *testing.T) {
	request, _ := http.NewRequest("GET", url, nil)
	response, err := mock.Do(request)
	assert.Equal(t, httpResponseExpected, response)
	assert.Nil(t, err)
}
func TestDoRequestWithErrorResponse(t *testing.T) {
	request, _ := http.NewRequest("GET", url, nil)
	response, err := mockWithError.Do(request)
	assert.Nil(t, response)
	assert.Equal(t, errFoo, err)
}
func TestDoRequestWithNoBody(t *testing.T) {
	request2, _ := http.NewRequest("GET", url, nil)
	stringReader := strings.NewReader(`{"test":"request body"}`)
	stringReadCloser := io.NopCloser(stringReader)
	request2.Body = stringReadCloser
	response, err := mockWithNoBody.Do(request2)
	assert.Nil(t, response)
	assert.Equal(t, errFoo, err)
}

func TestDoRequestMocked(t *testing.T) {
	request, _ := http.NewRequest("GET", url, nil)
	httpMockUp := httpclient.NewHTTPMultMock()
	mockedRequest := httpMockUp.Get(url)
	mockedRequest.Status(200)
	mockedRequest.Body(`{"test":"response body"}`)

	response, err := httpMockUp.Do(request)
	assert.Equal(t, httpResponseExpected, response)
	assert.Nil(t, err)
}
