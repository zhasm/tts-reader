package main

import (
	"io"
	"net/http"
	"time"
)

const MAX_RETRY = 10

// newHTTPRequest abstracts http.NewRequest with retry logic for request creation only.
func newHTTPRequest(method, url string, body io.Reader, headers map[string]string) (*http.Request, error) {
	var req *http.Request
	var lastErr error

	curlCmd := buildCurlCommand(method, url, headers, body)
	VPrintf("Curl: %s", curlCmd)

	err := RetryWithBackoff(
		func() error {
			var err error
			req, err = http.NewRequest(method, url, body)
			if err == nil {
				for k, v := range headers {
					req.Header.Set(k, v)
				}
				VPrintf("HTTP request for %s OK.", url)
			} else {
				VPrintf("HTTP request for %s failed: %v", url, err)
			}
			lastErr = err
			return err
		},
		MAX_RETRY,
		200*time.Millisecond,
	)
	if err != nil {
		return nil, lastErr
	}
	return req, nil
}

// buildCurlCommand constructs an equivalent curl command for debugging purposes
func buildCurlCommand(method, url string, headers map[string]string, body io.Reader) string {
	curlCmd := "curl"

	// Add method (only if not GET, since GET is default)
	if method != "GET" {
		curlCmd += " -X " + method
	}

	// Add headers
	for k, v := range headers {
		curlCmd += " -H '" + k + ": " + v + "'"
	}

	// Add body if present (for POST, PUT, etc.)
	if body != nil && (method == "POST" || method == "PUT" || method == "PATCH") {
		// Note: This is a simplified version. For actual body content,
		// you might want to read and escape the body content properly
		curlCmd += " -d '...'"
	}

	// Add URL
	curlCmd += " '" + url + "'"

	return curlCmd
}
