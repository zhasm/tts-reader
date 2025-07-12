package main

import (
	"io"
	"net/http"
	"time"
)

const MAX_RETRY = 10

// newHTTPRequestWithRetry abstracts http.NewRequest with retry logic for request creation only.
func newHTTPRequestWithRetry(method, url string, body io.Reader, headers map[string]string) (*http.Request, error) {
	var req *http.Request
	var err error
	delay := 200 * time.Millisecond
	maxAttempts := MAX_RETRY
	curlCmd := buildCurlCommand(method, url, headers, body)
	VPrintf("Curl: %s", curlCmd)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		req, err = http.NewRequest(method, url, body)
		if err == nil {
			// Set headers if provided
			for k, v := range headers {
				req.Header.Set(k, v)
			}
			VPrintf("HTTP request for %s OK. at %d th try. ", url, attempt)
			return req, nil
		}
		if attempt < maxAttempts {
			// Exponential backoff: delay doubles each time
			VPrintf("HTTP request for %s failed. waiting for the %d th try. ", url, attempt)
			time.Sleep(delay)
			delay *= 2
		}
	}
	return nil, err
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
