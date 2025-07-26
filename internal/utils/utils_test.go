package utils

import (
	"os"
	"testing"
)

func TestToHomeRelativePath(t *testing.T) {
	home, _ := os.UserHomeDir()
	path := home + "/test"
	rel := ToHomeRelativePath(path)
	if rel[:1] != "~" {
		t.Errorf("Expected path to start with ~, got %s", rel)
	}
}
