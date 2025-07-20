package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

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

	flag.VisitAll(func(f *flag.Flag) {
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
		flag.VisitAll(func(f *flag.Flag) {
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
	flag.Usage()
	os.Exit(code)
}

func PrintVersion() {
	if VersionInfo == "" {
		VersionInfo = "unknown"
	}
	fmt.Printf("tts-reader version %s\n", VersionInfo)
	os.Exit(0)
}

func ParseArgs() {
	// Set custom usage function
	flag.Usage = customUsage

	// Initialize supportedLangs first
	for _, l := range Langs {
		supportedLangs = append(supportedLangs, l.Name)
	}
	// Then register flags with both short and long names
	flag.BoolVar(&Verbose, "v", false, "verbose mode")
	flag.BoolVar(&Verbose, "verbose", false, "verbose mode")
	flag.StringVar(&Language, "l", "fr", "language ("+strings.Join(supportedLangs, ", ")+")")
	flag.StringVar(&Language, "language", "fr", "language ("+strings.Join(supportedLangs, ", ")+")")
	flag.Float64Var(&Speed, "s", 0.8, "speed (float)")
	flag.Float64Var(&Speed, "speed", 0.8, "speed (float)")
	flag.BoolVar(&Help, "h", false, "print help")
	flag.BoolVar(&Help, "help", false, "print help")
	flag.BoolVar(&Version, "V", false, "show version info")
	flag.BoolVar(&Version, "version", false, "show version info")

	flag.Parse()
	// Positional argument (content)
	remaining := flag.Args()
	if len(remaining) > 0 {
		Content = remaining[0]
	}

	// Validate language
	if !IsSupportedLang(Language) {
		langNames := make([]string, len(Langs))
		for i, l := range Langs {
			langNames[i] = l.Name
		}
		fmt.Printf("Invalid language. Choose from: %v\n", langNames)
		os.Exit(1)
	}

	if Help {
		PrintHelp(0)
	}

	if Version {
		PrintVersion()
	}

	if len(Content) == 0 {
		fmt.Println("Content arg missing~ ")
		PrintHelp(1)
	}

	// Print the parsed arguments
	logger.VPrintln("Parsed arguments:")
	logger.VPrintf("  Verbose : %v\n", Verbose)
	logger.VPrintf("  Language: %s\n", Language)
	logger.VPrintf("  Speed   : %v\n", Speed)
	logger.VPrintf("  Content : %s\n", Content)
}
