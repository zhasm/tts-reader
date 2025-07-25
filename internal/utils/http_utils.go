package utils

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"time"

	"github.com/zhasm/tts-reader/pkg/logger"
)

var (
	HIDDEN_KEYS = []string{"Ocp-Apim-Subscription-Key", "Authorization"}
)

const (
	MAX_RETRY = 5
)

// contains checks if a string is in a slice of strings
func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

func IsHiddenKey(key string) bool {
	return contains(HIDDEN_KEYS, key)
}

// NewHTTPRequestWithRetry abstracts http.NewRequest with retry logic for request creation only.
func NewHTTPRequestWithRetry(method, url string, body io.Reader, headers map[string]string) (*http.Request, error) {
	var req *http.Request
	var err error
	delay := 200 * time.Millisecond
	maxAttempts := MAX_RETRY
	curlCmd := buildCurlCommand(method, url, headers, body)
	curlCmdForLog := curlCmd
	for _, key := range HIDDEN_KEYS {
		// Compile the regex pattern
		regex := regexp.MustCompile(fmt.Sprintf("%s: \\w+", regexp.QuoteMeta(key)))
		// Create the replacement string
		replace := fmt.Sprintf("%s: [HIDDEN]", key)
		// Replace all occurrences in the curlCmdForLog
		curlCmdForLog = regex.ReplaceAllString(curlCmdForLog, replace)
	}
	logger.VPrintf("Curl: %s", curlCmdForLog)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		req, err = http.NewRequest(method, url, body)
		if err == nil {
			// Set headers if provided
			for k, v := range headers {
				req.Header.Set(k, v)
			}
			logger.VPrintf("HTTP request for %s OK. at %d th try. ", url, attempt)
			return req, nil
		}
		if attempt < maxAttempts {
			// Exponential backoff: delay doubles each time
			logger.VPrintf("HTTP request for %s failed. waiting for the %d th try. ", url, attempt)
			time.Sleep(delay)
			delay *= 2
		}
	}
	return nil, err
}

// HTTPRequest logs the request details, sends the HTTP request, and returns the response and error.
// Sensitive headers (like Ocp-Apim-Subscription-Key) are hidden in logs.
func HTTPRequest(client *http.Client, httpReq *http.Request) (*http.Response, error) {
	logger.VPrintf("Request URL: %s\n", httpReq.URL.String())
	logger.VPrintf("Request Method: %s\n", httpReq.Method)
	logger.VPrintf("Request Headers:\n")
	for key, values := range httpReq.Header {
		for _, value := range values {
			if IsHiddenKey(key) {
				logger.VPrintf("  %s: [HIDDEN]\n", key)
			} else {
				logger.VPrintf("  %s: %s\n", key, value)
			}
		}
	}

	logger.VPrintf("Sending request...\n")
	var resp *http.Response
	var err error
	initialInterval := time.Second
	err = RetryWithBackoff(func() error {
		resp, err = client.Do(httpReq)
		if err != nil || resp == nil {
			logger.VPrintf("HTTP request failed: %v\n", err)
			return err
		}
		return nil
	}, MAX_RETRY, initialInterval)
	if err != nil {
		logger.VPrintf("HTTP request failed after %d attempts: %v\n", MAX_RETRY, err)
		return nil, err
	}
	if resp == nil {
		logger.VPrintf("Error: Response is nil\n")
		return nil, fmt.Errorf("response is nil")
	}
	return resp, nil
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
