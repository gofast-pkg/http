package testify

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	ErrNilResponse = errors.New("response cannot be nil")
)

// HTTPClient type for call(s) http testing with a http.Client without execute on the network call
// It is useful to test a code that makes http calls
type HTTPClient interface {
	// Client returns a http.Client
	Client() *http.Client
	// AddCall adds a Caller to the list of callers
	// The callers list is ordered, so the first caller added will be the first to be executed
	// By default, if the caller does not have a response or error set,
	// a default response will be returned.
	// The default response return a 200 OK with a empty body
	// and a Content-Type header set to application/json.
	AddCall(call Caller)
	// SetDefaultResponse sets a default response to be returned by the http.Client
	// If no response or error are set.
	SetDefaultResponse(response *http.Response) error
	// RoundTrip method to implement http.RoundTripper interface
	RoundTrip(req *http.Request) (*http.Response, error)
	// ExpectedCalls asserts that all calls were executed
	// You need to call this method at the end of your test to ensure that all calls were executed
	ExpectedCalls()
}

type httpClient struct {
	defaultResponse *http.Response
	t               *testing.T
	callers         []Caller
}

// Caller type for call(s) http testing with a http.Client without execute on the network call
// You need set the ExpectedRequest field.
// This fields should be asserted to have equal values:
// - Body
// - ContentLength
// - Form
// - Header
// - Method
// - URL
// - context.Context
type Caller struct {
	ExpectedRequest *http.Request
	Response        *http.Response
	Err             error
}

func NewHTTPClient(t *testing.T) HTTPClient {
	header := make(http.Header)
	header.Set("Content-Type", "application/json")

	return &httpClient{
		t: t,
		defaultResponse: &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       nil,
			Header:     header,
		},
	}
}

func (c *httpClient) ExpectedCalls() {
	assert.EqualValues(
		c.t,
		len(c.callers),
		0,
		"Expect all calls to be executed, but some calls are missing, %v",
		c.callers)
}

func (c *httpClient) AddCall(call Caller) {
	c.callers = append(c.callers, call)
}

func (c *httpClient) SetDefaultResponse(response *http.Response) error {
	if response == nil {
		return ErrNilResponse
	}
	c.defaultResponse = response

	return nil
}

func (c *httpClient) Client() *http.Client {
	return &http.Client{
		Transport: c,
	}
}

// RoundTrip method to implement http.RoundTripper interface
func (c *httpClient) RoundTrip(req *http.Request) (*http.Response, error) {
	call := c.currentCall()

	c.assertRequest(call, req)
	if call.Err != nil {
		return nil, call.Err
	}
	if call.Response != nil {
		return call.Response, nil
	}

	return c.defaultResponse, nil
}

func (c *httpClient) currentCall() Caller {
	call := c.callers[0]
	c.callers = c.callers[1:]

	return call
}

func (c *httpClient) assertRequest(call Caller, req *http.Request) {
	assert.EqualValues(c.t, req.Body, call.ExpectedRequest.Body)
	assert.EqualValues(c.t, req.ContentLength, call.ExpectedRequest.ContentLength)
	assert.EqualValues(c.t, req.Form, call.ExpectedRequest.Form)
	assert.EqualValues(c.t, req.Header, call.ExpectedRequest.Header)
	assert.EqualValues(c.t, req.Method, call.ExpectedRequest.Method)
	assert.EqualValues(c.t, req.URL, call.ExpectedRequest.URL)
	assert.EqualValues(c.t, req.Context(), call.ExpectedRequest.Context())
}
