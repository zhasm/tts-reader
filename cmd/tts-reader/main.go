package main

import (
	"fmt"
	"os"

	"github.com/zhasm/tts-reader/pkg/config"
)

var (
	MAX_CONTENT_LENGTH_TO_SHOW = config.GetFromEnvOrDefault("MAX_CONTENT_LENGTH_TO_SHOW", "42")
)

func main() {
	if err := run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
