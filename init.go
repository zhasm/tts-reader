package main

import (
	"os"
	"strings"
)

var TTS_API_KEY string
var TTS_PATH = os.Getenv("HOME") + "/icloud/0-tmp/tts"
var R2_DB_TOKEN string

// isTest returns true if the program is running under go test
func isTest() bool {
	// Check if any of the test flags are present
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}
	return false
}

func init() {
	// Only parse args if not running tests
	if !isTest() {
		ParseArgs()
	}
	TTS_API_KEY = os.Getenv("TTS_API_KEY")
	if TTS_API_KEY == "" {
		VPrintln("Warning: TTS_API_KEY environment variable is not set")
		VPrintln("Please set the TTS_API_KEY environment variable:")
		VPrintln("export TTS_API_KEY=your_api_key_here")
		os.Exit(1)
	} else {
		VPrintf("TTS_API_KEY loaded successfully (length: %d)\n", len(TTS_API_KEY))
	}

	R2_DB_TOKEN = os.Getenv("R2_DB_TOKEN")
	if R2_DB_TOKEN == "" {
		VPrintln("Warning: R2_DB_TOKEN environment variable is not set")
		VPrintln("Please set the R2_DB_TOKEN environment variable:")
		VPrintln("export R2_DB_TOKEN=your_token_here")
		os.Exit(2)
	} else {
		VPrintf("R2_DB_TOKEN loaded successfully (length: %d)\n", len(R2_DB_TOKEN))
	}
}
