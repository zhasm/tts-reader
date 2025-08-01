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

// NewHTTPRequest abstracts http.NewRequest with retry logic for request creation only.
func NewHTTPRequest(method, url string, body io.Reader, headers map[string]string) (*http.Request, error) {
	var req *http.Request
	var err error
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
	logger.LogDebug("Curl: %s", curlCmdForLog)

	// Create the HTTP request once (no retry)
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	// Set headers if provided
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	logger.LogDebug("HTTP request for %s OK.", url)
	return req, nil
}

// HTTPRequest logs the request details, sends the HTTP request, and returns the response and error.
// Sensitive headers (like Ocp-Apim-Subscription-Key) are hidden in logs.
func HTTPRequest(client *http.Client, httpReq *http.Request) (*http.Response, error) {
	logger.LogDebug("Request URL: %s", httpReq.URL.String())
	logger.LogDebug("Request Method: %s", httpReq.Method)
	logger.LogDebug("Request Headers:")
	for key, values := range httpReq.Header {
		for _, value := range values {
			if IsHiddenKey(key) {
				logger.LogDebug("  %s: [HIDDEN]", key)
			} else {
				logger.LogDebug("  %s: %s", key, value)
			}
		}
	}

	logger.LogDebug("Sending request...")
	var resp *http.Response
	var err error
	initialInterval := time.Second
	err = RetryWithBackoff(func(retryIdx int) error {
		resp, err = client.Do(httpReq)
		if err != nil || resp == nil {
			logger.LogWarn("HTTP request failed %d: %v", retryIdx, err)
			return err
		}
		return nil
	}, MAX_RETRY, initialInterval)
	if err != nil {
		logger.LogError("HTTP request failed after %d attempts: %v", MAX_RETRY, err)
		return nil, err
	}
	if resp == nil {
		logger.LogError("Error: Response is nil")
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
