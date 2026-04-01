package http

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type errReader struct{}

func (e errReader) Read(_ []byte) (int, error) {
	return 0, assert.AnError
}

type closeErrorBody struct {
	io.Reader
}

func (c closeErrorBody) Close() error {
	return assert.AnError
}

func TestUnexpectedResponseBody(t *testing.T) {
	t.Run("should print an error message when close body fails", func(t *testing.T) {
		body := closeErrorBody{
			Reader: strings.NewReader("hello"),
		}

		resp := &http.Response{
			Status: "400 Bad Request",
			Header: http.Header{},
			Body:   body,
		}

		result := UnexpectedResponseBody(resp)

		expectedMessage := "response [status: 400 Bad Request] [header: map[]] " +
			"[body: hello]: assert.AnError general error for testing"
		assert.Equal(t, expectedMessage, result)
	})

	t.Run("should print an error when the response body is unreadable", func(t *testing.T) {
		resp := &http.Response{
			Status: "500 Internal Server Error",
			Header: http.Header{},
			Body:   io.NopCloser(errReader{}),
		}

		result := UnexpectedResponseBody(resp)

		expectedMessage := "response [status: 500 Internal Server Error] [header: map[]] " +
			"[body: failed to read response body]"
		assert.Equal(t, expectedMessage, result)
	})
	t.Run("should be ok", func(t *testing.T) {
		body := io.NopCloser(strings.NewReader("hello world"))

		resp := &http.Response{
			Status: "400 Bad Request",
			Header: http.Header{"X-Test": []string{"value"}},
			Body:   body,
		}

		result := UnexpectedResponseBody(resp)

		expectedMessage := "response [status: 400 Bad Request] [header: map[X-Test:[value]]] " +
			"[body: hello world]"
		assert.Equal(t, expectedMessage, result)
	})
	t.Run("should be ok with an http.NoBody", func(t *testing.T) {
		resp := &http.Response{
			Status: "400 Bad Request",
			Header: http.Header{"X-Test": []string{"value"}},
			Body:   http.NoBody,
		}

		result := UnexpectedResponseBody(resp)

		expectedMessage := "response [status: 400 Bad Request] [header: map[X-Test:[value]]] " +
			"[body: ]"
		assert.Equal(t, expectedMessage, result)
	})
	t.Run("should be ok with a nil Body", func(t *testing.T) {
		resp := &http.Response{
			Status: "400 Bad Request",
			Header: http.Header{"X-Test": []string{"value"}},
			Body:   nil,
		}

		result := UnexpectedResponseBody(resp)

		expectedMessage := "response [status: 400 Bad Request] [header: map[X-Test:[value]]] " +
			"[body: nil]"
		assert.Equal(t, expectedMessage, result)
	})
}
