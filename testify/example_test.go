package testify_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofast-pkg/http/testify"
	"github.com/stretchr/testify/assert"
)

const (
	urlTestRequest1 = "http://example.com/request_1"
	urlTestRequest2 = "http://example.com/request_2"
)

// testedCode simulates a code that makes two http requests
func testedCode(ctx context.Context, c *http.Client) error {
	var err error
	var r1 *http.Request
	var r2 *http.Request
	var resp1 *http.Response
	var resp2 *http.Response

	if r1, err = http.NewRequestWithContext(ctx, http.MethodGet, urlTestRequest1, nil); err != nil {
		return err
	}
	if r2, err = http.NewRequestWithContext(ctx, http.MethodGet, urlTestRequest2, nil); err != nil {
		return err
	}

	if resp1, err = c.Do(r1); err != nil {
		return err
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusOK {
		return errors.New("unexpected status code")
	}

	if resp2, err = c.Do(r2); err != nil {
		return err
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusCreated {
		return errors.New("unexpected status code")
	}

	return nil
}

func ExampleNewHTTPClient() {
	var err error
	var expectedR1 *http.Request
	var expectedR2 *http.Request
	t := new(testing.T)
	client := testify.NewHTTPClient(t)
	ctx := context.Background()

	if expectedR1, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		urlTestRequest1,
		nil); err != nil {
		t.Fatal(err)
	}
	if expectedR2, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		urlTestRequest2,
		nil); err != nil {
		t.Fatal(err)
	}

	client.AddCall(testify.Caller{
		ExpectedRequest: expectedR1,
	})
	client.AddCall(testify.Caller{
		ExpectedRequest: expectedR2,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
		},
	})

	err = testedCode(ctx, client.Client())
	if assert.NoError(t, err) {
		fmt.Println("no error")
		client.ExpectedCalls()
	}
	// Output:
	// no error
}
