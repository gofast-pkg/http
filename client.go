// Package http provides tooling around http in Go configured for production usage
package http

import (
	"net"
	"net/http"
	"time"
)

// Constants for http client configuration ready for production
const (
	timeoutInMs                   = 10000
	dialerTimeoutInSecond         = 30
	dialerKeepAliveInSecond       = 30
	maxIdleConns                  = 100
	idleConnTimeoutInSecond       = 10
	tlsHandshakeTimeoutInSecond   = 10
	expectContinueTimeoutInSecond = 1
	maxIdleConnsPerHost           = 100
)

// NewClient returns a new http client configured for production usage
func NewClient() *http.Client {
	return &http.Client{
		Timeout: time.Duration(timeoutInMs) * time.Millisecond,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   dialerTimeoutInSecond * time.Second,
				KeepAlive: dialerKeepAliveInSecond * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          maxIdleConns,
			IdleConnTimeout:       idleConnTimeoutInSecond * time.Second,
			TLSHandshakeTimeout:   tlsHandshakeTimeoutInSecond * time.Second,
			ExpectContinueTimeout: expectContinueTimeoutInSecond * time.Second,
			ForceAttemptHTTP2:     true,
			MaxIdleConnsPerHost:   maxIdleConnsPerHost,
		},
	}
}
