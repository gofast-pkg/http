package http

import (
	"fmt"
	"io"
	"net/http"
)

const (
	maxReadableResponseSize = 4 * 1024
	unreadableResponseBody  = "failed to read response body"
)

// UnexpectedResponseBody return the content body in string format.
// The function close the body.
// Use this fonction to log data when an unexpected response occurred.
//
//nolint:nonamedreturns // need to complete message.
func UnexpectedResponseBody(resp *http.Response) (message string) {
	message = fmt.Sprintf("response [status: %s] [header: %v]", resp.Status, resp.Header)

	if resp.Body == nil {
		message = fmt.Sprintf("%s [body: nil]", message)

		return
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			message = fmt.Sprintf("%s: %s", message, err.Error())
		}
	}()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxReadableResponseSize))
	if err != nil {
		return fmt.Sprintf("%s [body: %s]", message, unreadableResponseBody)
	}

	message = fmt.Sprintf("%s [body: %s]", message, body)

	return message
}
