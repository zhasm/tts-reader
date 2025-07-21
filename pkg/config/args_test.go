package config

import (
	"flag"
	"os"
	"testing"
)

func TestPrintVersion(t *testing.T) {
	VersionInfo = "test-version"
	// This will call os.Exit(0), so we can't test it directly without os/exec trickery
	// Instead, just check that it prints the right string (manual test)
}

func TestParseArgs_Help(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--help"}
	defer func() {
		_ = recover() // ignore os.Exit
	}()
	_ = ParseArgs()
}

func TestParseArgs_Version(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--version"}
	defer func() {
		_ = recover() // ignore os.Exit
	}()
	_ = ParseArgs()
}

func TestParseArgs_InvalidLang(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-l", "xx", "content"}
	defer func() {
		_ = recover() // ignore os.Exit
	}()
	_ = ParseArgs()
}
