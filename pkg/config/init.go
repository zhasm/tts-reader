package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zhasm/tts-reader/pkg/logger"
)

var (
	TTS_PATH = GetStringFromEnvOrDefault("TTS_PATH", "~/icloud/0-tmp/tts")
)

var TTS_API_KEY string
var R2_DB_TOKEN string

func resolveDefaultValue(defaultValue string) string {
	if len(defaultValue) >= 2 && defaultValue[:2] == "~/" {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, defaultValue[2:])
		}
	}
	return defaultValue
}

func GetIntFromEnvOrDefault(envVar, defaultValue string) int {
	if value := os.Getenv(envVar); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	if intValue, err := strconv.Atoi(resolveDefaultValue(defaultValue)); err == nil {
		return intValue
	}
	return 0
}

func GetStringFromEnvOrDefault(envVar, defaultValue string) string {
	if value := os.Getenv(envVar); value != "" {
		return value
	}
	return resolveDefaultValue(defaultValue)
}

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
			logger.LogError("Error parsing args:", err)
			os.Exit(1)
		}
	}

	TTS_API_KEY = os.Getenv("TTS_API_KEY")
	if TTS_API_KEY == "" {
		logger.LogError("Warning: TTS_API_KEY environment variable is not set")
		logger.LogError("Please set the TTS_API_KEY environment variable:")
		logger.LogError("export TTS_API_KEY=your_api_key_here")
		os.Exit(1)
	} else {
		logger.LogDebug("TTS_API_KEY loaded successfully (length: %d)", len(TTS_API_KEY))
	}

	R2_DB_TOKEN = os.Getenv("R2_DB_TOKEN")
	if R2_DB_TOKEN == "" {
		logger.LogError("Warning: R2_DB_TOKEN environment variable is not set")
		logger.LogError("Please set the R2_DB_TOKEN environment variable:")
		logger.LogError("export R2_DB_TOKEN=your_token_here")
		os.Exit(2)
	} else {
		logger.LogDebug("R2_DB_TOKEN loaded successfully (length: %d)", len(R2_DB_TOKEN))
	}
}
