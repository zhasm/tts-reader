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
		if err == nil {
			t.Errorf("Expected non-zero exit code when no arguments are provided")
		}
		if !strings.Contains(string(output), "Usage of") {
			t.Errorf("Expected usage information, but got: %s", string(output))
		}
	})

	t.Run("only verbose flag", func(t *testing.T) {
		config.ResetArgs()
		cmd := exec.Command(buildPath, "-v")
		output, err := cmd.CombinedOutput()
		// Should exit 0 and do nothing
		if err != nil {
			t.Errorf("Expected no error when only -v is provided, but got: %v", err)
		}
		if len(output) != 0 && !strings.Contains(string(output), "Usage of") {
			t.Errorf("Expected no output or usage info, got: %s", string(output))
		}
	})

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

	t.Run("help flag", func(t *testing.T) {
		cmd := exec.Command(buildPath, "--help")
		output, err := cmd.CombinedOutput()
		if err != nil && !strings.Contains(string(output), "Usage of") {
			t.Errorf("Expected usage info with --help, got: %s", string(output))
		}
		if !strings.Contains(string(output), "Usage of") {
			t.Errorf("Expected usage information, but got: %s", string(output))
		}
	})

	t.Run("version flag", func(t *testing.T) {
		cmd := exec.Command(buildPath, "--version")
		output, err := cmd.CombinedOutput()
		if err != nil && !strings.Contains(string(output), "tts-reader version") {
			t.Errorf("Expected version info with --version, got: %s", string(output))
		}
		if !strings.Contains(string(output), "tts-reader version") {
			t.Errorf("Expected version information, but got: %s", string(output))
		}
	})

	t.Run("server mode with port", func(t *testing.T) {
		cmd := exec.Command(buildPath, "serve", "--port", "9999", "--help")
		output, err := cmd.CombinedOutput()
		if err != nil && !strings.Contains(string(output), "Usage of") {
			t.Errorf("Expected usage info with serve --help, got: %s", string(output))
		}
		if !strings.Contains(string(output), "Usage of") {
			t.Errorf("Expected usage information for server mode, but got: %s", string(output))
		}
	})

	t.Run("server mode with invalid port", func(t *testing.T) {
		cmd := exec.Command(buildPath, "serve", "--port", "99999")
		output, err := cmd.CombinedOutput()
		if err == nil {
			t.Errorf("Expected error with invalid port, but got nil")
		}
		if !strings.Contains(string(output), "invalid port") {
			t.Errorf("Expected invalid port error, but got: %s", string(output))
		}
	})
}
