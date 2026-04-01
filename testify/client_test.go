package testify

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHTTPClient(t *testing.T) {
	t.Run("Should be ok", func(t *testing.T) {
		client := NewHTTPClient(t)

		c := client.(*httpClient)
		if assert.NotNil(t, c) {
			assert.Empty(t, c.callers)
			assert.NotNil(t, c.defaultResponse)
			assert.Equal(t, t, c.t)
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
		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		c := &httpClient{
			t: t,
			defaultResponse: &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       nil,
				Header:     header,
			},
		}

		caller := Caller{
			ExpectedRequest: nil,
			Response:        nil,
			Err:             nil,
		}

		c.AddCall(caller)

		assert.Len(t, c.callers, 1)
	})
}

func TestHttpClient_SetDefaultResponse(t *testing.T) {
	t.Run("Should return an error with a nil response parameter", func(t *testing.T) {
		client := NewHTTPClient(t)

		err := client.SetDefaultResponse(nil)
		require.ErrorIs(t, err, ErrNilResponse)
	})
	t.Run("Should be ok", func(t *testing.T) {
		client := NewHTTPClient(t)

		err := client.SetDefaultResponse(&http.Response{})
		require.NoError(t, err)
		assert.Equal(t, &http.Response{}, client.(*httpClient).defaultResponse)
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
		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		c := &httpClient{
			t: t,
			defaultResponse: &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       nil,
				Header:     header,
			},
		}
		c.AddCall(Caller{
			ExpectedRequest: req,
		})

		resp, err := c.RoundTrip(req)
		defer func() {
			if resp.Body != nil {
				assert.NoError(t, resp.Body.Close())
			}
		}()
		require.NoError(t, err)
		assert.Equal(t, c.defaultResponse, resp)
	})
	t.Run("Should be ok with custom response", func(t *testing.T) {
		req := &http.Request{}
		client := NewHTTPClient(t)
		client.AddCall(Caller{
			ExpectedRequest: req,
			Response:        &http.Response{StatusCode: http.StatusCreated},
		})

		resp, err := client.RoundTrip(req)
		defer func() {
			if resp.Body != nil {
				assert.NoError(t, resp.Body.Close())
			}
		}()
		require.NoError(t, err)
		assert.Equal(t, &http.Response{StatusCode: http.StatusCreated}, resp)
	})
	t.Run("Should be ok with custom error", func(t *testing.T) {
		req := &http.Request{}
		client := NewHTTPClient(t)
		client.AddCall(Caller{
			ExpectedRequest: req,
			Err:             assert.AnError,
		})

		resp, err := client.RoundTrip(req)
		assert.Nil(t, resp)
		require.ErrorIs(t, err, assert.AnError)
	})
}
