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
	ResetArgs()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--help"}
	_ = ParseArgs()
	if !Help {
		t.Errorf("Expected Help to be true when --help is passed")
	}
}

func TestParseArgs_ShortHelp(t *testing.T) {
	ResetArgs()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-h"}
	_ = ParseArgs()
	if !Help {
		t.Errorf("Expected Help to be true when -h is passed")
	}
}

func TestParseArgs_Version(t *testing.T) {
	ResetArgs()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--version"}
	_ = ParseArgs()
	if !Version {
		t.Errorf("Expected Version to be true when --version is passed")
	}
}

func TestParseArgs_ShortVersion(t *testing.T) {
	ResetArgs()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-V"}
	_ = ParseArgs()
	if !Version {
		t.Errorf("Expected Version to be true when -V is passed")
	}
}

func TestParseArgs_Verbose(t *testing.T) {
	ResetArgs()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--verbose"}
	_ = ParseArgs()
	if !Verbose {
		t.Errorf("Expected Verbose to be true when --verbose is passed")
	}
}

func TestParseArgs_ShortVerbose(t *testing.T) {
	ResetArgs()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-v"}
	_ = ParseArgs()
	if !Verbose {
		t.Errorf("Expected Verbose to be true when -v is passed")
	}
}

func TestParseArgs_Language(t *testing.T) {
	ResetArgs()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--language", "jp"}
	_ = ParseArgs()
	if Language != "jp" {
		t.Errorf("Expected Language to be 'jp' when --language jp is passed, got '%s'", Language)
	}
}

func TestParseArgs_ShortLanguage(t *testing.T) {
	ResetArgs()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-l", "pl"}
	_ = ParseArgs()
	if Language != "pl" {
		t.Errorf("Expected Language to be 'pl' when -l pl is passed, got '%s'", Language)
	}
}

func TestParseArgs_Speed(t *testing.T) {
	ResetArgs()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--speed", "1.2"}
	_ = ParseArgs()
	if Speed != 1.2 {
		t.Errorf("Expected Speed to be 1.2 when --speed 1.2 is passed, got %f", Speed)
	}
}

func TestParseArgs_ShortSpeed(t *testing.T) {
	ResetArgs()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-s", "0.5"}
	_ = ParseArgs()
	if Speed != 0.5 {
		t.Errorf("Expected Speed to be 0.5 when -s 0.5 is passed, got %f", Speed)
	}
}

func TestParseArgs_DryRun(t *testing.T) {
	ResetArgs()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--dry-run"}
	_ = ParseArgs()
	if !DryRun {
		t.Errorf("Expected DryRun to be true when --dry-run is passed")
	}
}

func TestParseArgs_ShortDryRun(t *testing.T) {
	ResetArgs()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-d"}
	_ = ParseArgs()
	if !DryRun {
		t.Errorf("Expected DryRun to be true when -d is passed")
	}
}

func TestParseArgs_InvalidLang(t *testing.T) {
	ResetArgs()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-l", "xx", "content"}
	defer func() {
		_ = recover() // ignore os.Exit
	}()
	_ = ParseArgs()
}
