package config

import (
	"os"
	"strings"

	"github.com/zhasm/tts-reader/pkg/logger"
)

const (
	TTS_SUB_PATH = "/icloud/0-tmp/tts"
)

var TTS_API_KEY string
var TTS_PATH = os.Getenv("HOME") + TTS_SUB_PATH
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

func Init() {
	// Only parse args if not running tests
	if !isTest() {
		if err := ParseArgs(); err != nil {
			logger.VPrintln("Error parsing args:", err)
			os.Exit(1)
		}
	}
	TTS_API_KEY = os.Getenv("TTS_API_KEY")
	if TTS_API_KEY == "" {
		logger.VPrintln("Warning: TTS_API_KEY environment variable is not set")
		logger.VPrintln("Please set the TTS_API_KEY environment variable:")
		logger.VPrintln("export TTS_API_KEY=your_api_key_here")
		os.Exit(1)
	} else {
		logger.VPrintf("TTS_API_KEY loaded successfully (length: %d)\n", len(TTS_API_KEY))
	}

	R2_DB_TOKEN = os.Getenv("R2_DB_TOKEN")
	if R2_DB_TOKEN == "" {
		logger.VPrintln("Warning: R2_DB_TOKEN environment variable is not set")
		logger.VPrintln("Please set the R2_DB_TOKEN environment variable:")
		logger.VPrintln("export R2_DB_TOKEN=your_token_here")
		os.Exit(2)
	} else {
		logger.VPrintf("R2_DB_TOKEN loaded successfully (length: %d)\n", len(R2_DB_TOKEN))
	}
}
