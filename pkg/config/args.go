package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/spf13/pflag"
	"github.com/zhasm/tts-reader/pkg/logger"
)

const (
	DEFAULT_LOG_LEVEL = "info"
)

var (
	Language    string
	Speed       float64 = 0.8
	Content     string
	Help        bool
	Version     bool
	VersionInfo string
	GenConfig   bool
	DryRun      bool
	OverWrite   bool
	LogLevel    string = DEFAULT_LOG_LEVEL
	ConfigFile  string
)

// Dynamic usage function that handles all flags
func customUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

	pflag.VisitAll(func(f *pflag.Flag) {
		// Determine type string
		var flagType string
		if f.Value.Type() == "bool" {
			flagType = ""
		} else {
			flagType = f.Value.Type()
		}

		if f.Shorthand != "" {
			fmt.Fprintf(os.Stderr, "  -%s, --%s %s\n", f.Shorthand, f.Name, flagType)
		} else {
			fmt.Fprintf(os.Stderr, "      --%s %s\n", f.Name, flagType)
		}
		fmt.Fprintf(os.Stderr, "    \t%s\n", f.Usage)
	})
}

func PrintHelp(code int) {
	pflag.Usage()
	os.Exit(code)
}

func PrintVersion() {
	if VersionInfo == "" {
		VersionInfo = "unknown"
	}
	fmt.Printf("tts-reader version %s\n", VersionInfo)
	os.Exit(0)
}

var parseOnce sync.Once

func ParseArgs() error {
	// Set custom usage function
	pflag.Usage = customUsage

	var parseErr error
	parseOnce.Do(func() {
		// Register flags with both short and long names using VarP
		pflag.StringVarP(&LogLevel, "log-level", "L", DEFAULT_LOG_LEVEL, "log level: debug(d), info(i), warn(w), error(e)")
		pflag.StringVar(&ConfigFile, "config", "", "config file path")
		pflag.StringVarP(&Language, "language", "l", "fr", "language ("+GetAllLangShortNamesStr()+")")
		pflag.Float64VarP(&Speed, "speed", "s", 0.8, "speed (float)")
		pflag.BoolVarP(&Help, "help", "h", false, "print help")
		pflag.BoolVar(&GenConfig, "gen-config", false, "generate default tts-langs.yml config file")
		pflag.BoolVarP(&Version, "version", "V", false, "show version info")
		pflag.BoolVarP(&DryRun, "dry-run", "d", false, "dry run mode (no changes will be made)")
		pflag.BoolVarP(&OverWrite, "over-write", "o", false, "force re-download even if file exists")

		pflag.Parse()
		// Validate log level
		LogLevel = strings.ToLower(LogLevel)
		// Map single-letter aliases to full log level names
		switch LogLevel {
		case "d":
			LogLevel = "debug"
		case "i":
			LogLevel = "info"
		case "w":
			LogLevel = "warn"
		case "e":
			LogLevel = "error"
		}
		validLogLevels := map[string]bool{
			"debug": true,
			"info":  true,
			"warn":  true,
			"error": true,
		}
		if !validLogLevels[LogLevel] {
			fmt.Fprintf(os.Stderr, "Error: invalid log level: %s\n", LogLevel)
			fmt.Fprintf(os.Stderr, "Valid log levels: debug(d), info(i), warn(w), error(e)\n")
			os.Exit(2)
		}
		// Set logger log level
		logger.SetLogLevel(LogLevel)
		// Positional argument (content)
		remaining := pflag.Args()
		if len(remaining) > 0 {
			Content = strings.TrimSpace(remaining[0])
		}
	})
	return parseErr
}

// ValidateAndHandleArgs checks for help/version flags, missing content, and language validity. Exits if any are triggered.
func ValidateAndHandleArgs() error {
	if Version {
		PrintVersion()
		return nil
	}
	if GenConfig {
		GenerateConfigFile()
		os.Exit(0)
	}
	if Help {
		PrintHelp(0)
		return nil
	}
	if Content == "" {
		// If no arguments at all were provided, show help.
		if len(os.Args) == 1 {
			PrintHelp(0)
			return fmt.Errorf("content argument is missing")
		}
		fmt.Fprintln(os.Stderr, "Error: content argument is missing.")
		PrintHelp(1)
		return fmt.Errorf("content argument is missing")
	}

	// Check if Content is an HTTP/HTTPS URL (case-insensitive)
	urlRegex := regexp.MustCompile(`(?i)https?://`)
	if urlRegex.MatchString(Content) {
		logger.LogError("Content is a URL: %s", Content)
		msg := "content must not contain http/https"
		return fmt.Errorf("%s", msg)
	}

	return nil
}

// ResetArgs resets all flag variables and parseOnce for testing
func ResetArgs() {
	LogLevel = DEFAULT_LOG_LEVEL
	Language = "fr"
	Speed = 0.8
	Content = ""
	Help = false
	Version = false
	VersionInfo = ""
	GenConfig = false
	DryRun = false
	OverWrite = false
	ConfigFile = ""
	parseOnce = sync.Once{}
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
}
