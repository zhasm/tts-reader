package config

import (
	"flag"
	"os"
	"testing"
)

func TestIsTest(t *testing.T) {
	os.Args = []string{"cmd", "-test.v"}
	if !isTest() {
		t.Error("Expected isTest to return true for test flag")
	}
	os.Args = []string{"cmd"}
	if isTest() {
		t.Error("Expected isTest to return false for no test flag")
	}
}

func TestInit(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Setenv("TTS_API_KEY", "test-key")
	os.Setenv("R2_DB_TOKEN", "test-token")
	defer os.Unsetenv("TTS_API_KEY")
	defer os.Unsetenv("R2_DB_TOKEN")
	// This will call os.Exit on error, so we can't test the error path easily
	// Just test the happy path
	Init()
	if TTS_API_KEY != "test-key" {
		t.Errorf("Expected TTS_API_KEY to be 'test-key', got '%s'", TTS_API_KEY)
	}
	if R2_DB_TOKEN != "test-token" {
		t.Errorf("Expected R2_DB_TOKEN to be 'test-token', got '%s'", R2_DB_TOKEN)
	}
}
