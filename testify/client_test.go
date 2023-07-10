package testify

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTPClient(t *testing.T) {
	t.Run("Should be ok", func(t *testing.T) {
		client := NewHTTPClient(t)

		c := client.(*httpClient)
		if assert.NotNil(t, c) {
			assert.Empty(t, c.callers)
			assert.NotNil(t, c.defaultResponse)
			assert.EqualValues(t, t, c.t)
		}
	})
}

func TestHttpClient_ExpectedCalls(t *testing.T) {
	t.Run("Should be ok", func(t *testing.T) {
		client := NewHTTPClient(t)

		client.ExpectedCalls()
	})
}

func TestHttpClient_AddCall(t *testing.T) {
	t.Run("Should be ok", func(t *testing.T) {
		client := NewHTTPClient(t)

		caller := Caller{
			ExpectedRequest: nil,
			Response:        nil,
			Err:             nil,
		}

		client.AddCall(caller)

		c := client.(*httpClient)
		if assert.NotNil(t, c) {
			assert.Len(t, c.callers, 1)
		}
	})
}

func TestHttpClient_SetDefaultResponse(t *testing.T) {
	t.Run("Should return an error with a nil response parameter", func(t *testing.T) {
		client := NewHTTPClient(t)

		err := client.SetDefaultResponse(nil)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrNilResponse)
		}
	})
	t.Run("Should be ok", func(t *testing.T) {
		client := NewHTTPClient(t)

		err := client.SetDefaultResponse(&http.Response{})
		if assert.NoError(t, err) {
			assert.EqualValues(t, &http.Response{}, client.(*httpClient).defaultResponse)
		}
	})
}

func TestHttpClient_Client(t *testing.T) {
	t.Run("Should be ok", func(t *testing.T) {
		client := NewHTTPClient(t)

		assert.NotNil(t, client.Client())
	})
}

func TestHttpClient_RoundTrip(t *testing.T) {
	t.Run("Should be ok with default response", func(t *testing.T) {
		req := &http.Request{}
		client := NewHTTPClient(t)
		client.AddCall(Caller{
			ExpectedRequest: req,
		})

		resp, err := client.RoundTrip(req)
		if assert.NoError(t, err) {
			if resp.Body != nil {
				defer resp.Body.Close()
			}
			assert.EqualValues(t, client.(*httpClient).defaultResponse, resp)
		}
	})
	t.Run("Should be ok with custom response", func(t *testing.T) {
		req := &http.Request{}
		client := NewHTTPClient(t)
		client.AddCall(Caller{
			ExpectedRequest: req,
			Response:        &http.Response{StatusCode: http.StatusCreated},
		})

		resp, err := client.RoundTrip(req)
		if assert.NoError(t, err) {
			if resp.Body != nil {
				defer resp.Body.Close()
			}
			assert.EqualValues(t, &http.Response{StatusCode: http.StatusCreated}, resp)
		}
	})
	t.Run("Should be ok with custom error", func(t *testing.T) {
		errCustom := errors.New("custom error")
		req := &http.Request{}
		client := NewHTTPClient(t)
		client.AddCall(Caller{
			ExpectedRequest: req,
			Err:             errCustom,
		})

		resp, err := client.RoundTrip(req)
		if assert.Error(t, err) {
			assert.Nil(t, resp)
			assert.ErrorIs(t, err, errCustom)
			assert.Nil(t, resp)
		}
		if resp != nil {
			defer resp.Body.Close()
		}
	})
}
