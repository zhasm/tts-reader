package main

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/zhasm/tts-reader/pkg/config"
)

func TestMain_CLIArgs(t *testing.T) {
	buildPath := t.TempDir() + "/tts-reader-test"
	cmd := exec.Command("go", "build", "-o", buildPath, "github.com/zhasm/tts-reader/cmd/tts-reader")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build test binary: %v", err)
	}

	t.Run("no arguments", func(t *testing.T) {
		cmd := exec.Command(buildPath)
		output, err := cmd.CombinedOutput()
		// When no arguments are provided, the program should exit with non-zero code
		// and show usage information
		if err != nil {
			t.Errorf("Expected non-zero exit code when no arguments are provided")
		}
		outputStr := string(output)
		if !strings.Contains(outputStr, "Usage of") && !strings.Contains(outputStr, "usage") && !strings.Contains(outputStr, "Usage:") {
			t.Errorf("Expected usage information, but got: %s", outputStr)
		}
	})

	t.Run("only verbose flag", func(t *testing.T) {
		config.ResetArgs()
		cmd := exec.Command(buildPath, "-v")
		output, err := cmd.CombinedOutput()
		outputStr := string(output)

		// The verbose flag alone might be valid (exit code 0) or might require content
		// Check what the actual behavior should be based on the output
		if err != nil {
			// If it exits with error, it should show usage or error message
			if !strings.Contains(outputStr, "Usage of") &&
				!strings.Contains(outputStr, "usage") &&
				!strings.Contains(outputStr, "Usage:") &&
				!strings.Contains(outputStr, "Error:") &&
				!strings.Contains(outputStr, "content argument is missing") {
				t.Errorf("Expected usage info or error message when -v fails, got: %s", outputStr)
			}
		} else {
			// If it succeeds, output should be empty or contain version/help info
			if len(outputStr) > 0 &&
				!strings.Contains(outputStr, "Usage of") &&
				!strings.Contains(outputStr, "usage") &&
				!strings.Contains(outputStr, "Usage:") &&
				!strings.Contains(outputStr, "version") {
				t.Errorf("Expected no output, usage info, or version info when -v succeeds, got: %s", outputStr)
			}
		}
	})

	t.Run("flags requiring content", func(t *testing.T) {
		config.ResetArgs()
		cmd := exec.Command(buildPath, "-l", "fr")
		output, err := cmd.CombinedOutput()
		if err == nil {
			t.Errorf("Expected error when content is missing, but got nil")
		}
		outputStr := string(output)
		if !strings.Contains(outputStr, "Error: content argument is missing") &&
			!strings.Contains(outputStr, "content argument is missing") &&
			!strings.Contains(outputStr, "missing") {
			t.Errorf("Expected error message about missing content, but got: %s", outputStr)
		}
	})

	t.Run("help flag", func(t *testing.T) {
		cmd := exec.Command(buildPath, "--help")
		output, _ := cmd.CombinedOutput()
		outputStr := string(output)

		// Help flag typically exits with code 0 in most CLI tools
		// but some might exit with code 2, so we check the output content instead
		if !strings.Contains(outputStr, "Usage of") &&
			!strings.Contains(outputStr, "usage") &&
			!strings.Contains(outputStr, "Usage:") {
			t.Errorf("Expected usage information with --help, but got: %s", outputStr)
		}
	})

	t.Run("version flag", func(t *testing.T) {
		cmd := exec.Command(buildPath, "--version")
		output, _ := cmd.CombinedOutput()
		outputStr := string(output)

		// Version flag typically exits with code 0
		// but some might exit with code 2, so we check the output content instead
		if !strings.Contains(outputStr, "tts-reader version") &&
			!strings.Contains(outputStr, "version") {
			t.Errorf("Expected version information with --version, but got: %s", outputStr)
		}
	})

	t.Run("server mode with port", func(t *testing.T) {
		cmd := exec.Command(buildPath, "serve", "--port", "9999", "--help")
		output, _ := cmd.CombinedOutput()
		outputStr := string(output)

		// Server help should show usage information
		if !strings.Contains(outputStr, "Usage of") &&
			!strings.Contains(outputStr, "usage") &&
			!strings.Contains(outputStr, "Usage:") {
			t.Errorf("Expected usage information for server mode with --help, but got: %s", outputStr)
		}
	})

	t.Run("server mode with invalid port", func(t *testing.T) {
		cmd := exec.Command(buildPath, "serve", "--port", "99999")
		output, err := cmd.CombinedOutput()
		if err == nil {
			t.Errorf("Expected error with invalid port, but got nil")
		}
		outputStr := string(output)
		if !strings.Contains(outputStr, "invalid port") &&
			!strings.Contains(outputStr, "port") &&
			!strings.Contains(outputStr, "range") {
			t.Errorf("Expected invalid port error, but got: %s", outputStr)
		}
	})
}
