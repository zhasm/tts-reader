package utils

import (
	"os"
	"strings"
	"time"
)

// Convert absolute path to relative path from home directory
func ToHomeRelativePath(absPath string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return absPath // Return original path if we can't get home directory
	}

	// Check if the path starts with the home directory
	if strings.HasPrefix(absPath, homeDir) {
		// Replace home directory with ~
		relativePath := strings.Replace(absPath, homeDir, "~", 1)
		return relativePath
	}

	return absPath // Return original path if it's not under home directory
}

// RetryWithBackoff retries the provided function with exponential backoff.
func RetryWithBackoff(fn func() error, maxRetries int, initialInterval time.Duration) error {
	interval := initialInterval
	var lastErr error
	for i := range maxRetries {
		err := fn()
		if err == nil {
			return nil
		}
		lastErr = err
		if i < maxRetries-1 {
			time.Sleep(interval)
			interval *= 2
		}
	}
	return lastErr
}
