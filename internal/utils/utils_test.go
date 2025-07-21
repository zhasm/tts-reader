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

func TestGetFuncName(t *testing.T) {
	name := GetFuncName(TestGetFuncName)
	if name == "" {
		t.Error("Expected non-empty function name")
	}
}
