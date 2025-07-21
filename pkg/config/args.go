package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/pflag"
	"github.com/zhasm/tts-reader/pkg/logger"
)

var (
	Verbose     bool
	Language    string
	Speed       float64 = 0.8
	Content     string
	Help        bool
	Version     bool
	VersionInfo string
	DryRun      bool
)

// flagMapping maps short flags to their corresponding long flags
var flagMapping = map[string]string{
	"v": "verbose",
	"l": "language",
	"s": "speed",
	"h": "help",
	"V": "version",
}

// Dynamic usage function that groups short and long flags
func customUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

	// Visit all registered flags and group them
	flagGroups := make(map[string][]string)
	flagDescriptions := make(map[string]string)

	pflag.VisitAll(func(f *pflag.Flag) {
		// Find the corresponding short/long flag
		var shortFlag, longFlag string
		for short, long := range flagMapping {
			if f.Name == short {
				shortFlag = short
				longFlag = long
				break
			} else if f.Name == long {
				shortFlag = short
				longFlag = long
				break
			}
		}

		if shortFlag != "" {
			groupKey := shortFlag + "," + longFlag
			flagGroups[groupKey] = []string{shortFlag, longFlag}
			flagDescriptions[groupKey] = f.Usage
		}
	})

	// Print grouped flags
	for groupKey, flags := range flagGroups {
		shortFlag := flags[0]
		longFlag := flags[1]
		description := flagDescriptions[groupKey]

		// Get the flag value to determine if it's a string/float/bool
		var flagType string
		pflag.VisitAll(func(f *pflag.Flag) {
			if f.Name == shortFlag || f.Name == longFlag {
				switch f.Value.String() {
				case "true", "false":
					flagType = ""
				default:
					if strings.Contains(f.Name, "language") {
						flagType = "string"
					} else if strings.Contains(f.Name, "speed") {
						flagType = "float"
					}
				}
			}
		})

		if flagType != "" {
			fmt.Fprintf(os.Stderr, "  -%s, --%s %s\n", shortFlag, longFlag, flagType)
			fmt.Fprintf(os.Stderr, "    \t%s\n", description)
		} else {
			fmt.Fprintf(os.Stderr, "  -%s, --%s\t%s\n", shortFlag, longFlag, description)
		}
	}
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

	// Initialize supportedLangs first
	initSupportedLangs()

	var parseErr error
	parseOnce.Do(func() {
		// Register flags with both short and long names using VarP
		pflag.BoolVarP(&Verbose, "verbose", "v", false, "verbose mode")
		pflag.StringVarP(&Language, "language", "l", "fr", "language ("+strings.Join(supportedLangs, ", ")+")")
		pflag.Float64VarP(&Speed, "speed", "s", 0.8, "speed (float)")
		pflag.BoolVarP(&Help, "help", "h", false, "print help")
		pflag.BoolVarP(&Version, "version", "V", false, "show version info")
		pflag.BoolVarP(&DryRun, "dry-run", "d", false, "dry run mode (no changes will be made)")

		pflag.Parse()
		// Set logger verbose flag
		logger.SetVerbose(Verbose)
		// Positional argument (content)
		remaining := pflag.Args()
		if len(remaining) > 0 {
			Content = remaining[0]
		}
	})
	return parseErr
}

// ResetArgs resets all flag variables and parseOnce for testing
func ResetArgs() {
	Verbose = false
	Language = "fr"
	Speed = 0.8
	Content = ""
	Help = false
	Version = false
	VersionInfo = ""
	DryRun = false
	parseOnce = sync.Once{}
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
}
