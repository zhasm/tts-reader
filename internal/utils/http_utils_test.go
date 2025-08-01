package utils

import (
	"testing"
)

func TestBuildCurlCommand(t *testing.T) {
	cmd := buildCurlCommand("POST", "http://example.com", map[string]string{"A": "B"}, nil)
	if cmd == "" {
		t.Error("Expected non-empty curl command")
	}
}

func TestNewHTTPRequest_InvalidURL(t *testing.T) {
	_, err := NewHTTPRequest(":badmethod", "://bad-url", nil, nil)
	if err == nil {
		t.Error("Expected error for invalid URL/method")
	}
}
