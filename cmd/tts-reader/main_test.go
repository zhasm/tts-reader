package main

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/zhasm/tts-reader/pkg/config"
)

// Test for the main function's argument handling
func TestMain(t *testing.T) {
	// Build the binary to ensure it's up to date
	buildPath := t.TempDir() + "/tts-reader-test"
	cmd := exec.Command("go", "build", "-o", buildPath, "github.com/zhasm/tts-reader/cmd/tts-reader")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build test binary: %v", err)
	}

	// Test case: no arguments
	t.Run("no arguments", func(t *testing.T) {
		cmd := exec.Command(buildPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if exitErr.ExitCode() != 0 {
					t.Errorf("Expected exit code 0 when no arguments are provided, but got %d", exitErr.ExitCode())
				}
			} else {
				t.Errorf("Expected no error when no arguments are provided, but got: %v", err)
			}
		}
		if !strings.Contains(string(output), "Usage of") {
			t.Errorf("Expected usage information, but got: %s", string(output))
		}
	})

	// Test case: flags requiring content
	t.Run("flags requiring content", func(t *testing.T) {
		config.ResetArgs()
		cmd := exec.Command(buildPath, "-l", "fr")
		output, err := cmd.CombinedOutput()
		if err == nil {
			t.Errorf("Expected error when content is missing, but got nil")
		}
		if !strings.Contains(string(output), "Error: content argument is missing") {
			t.Errorf("Expected error message about missing content, but got: %s", string(output))
		}
	})
}
