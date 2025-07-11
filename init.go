package main

import (
	"os"
)

var TTS_API_KEY string
var TTS_PATH = os.Getenv("HOME") + "/icloud/0-tmp/tts"
var R2_DB_TOKEN string

func init() {
	ParseArgs()
	TTS_API_KEY = os.Getenv("TTS_API_KEY")
	if TTS_API_KEY == "" {
		VPrintln("Warning: TTS_API_KEY environment variable is not set")
		VPrintln("Please set the TTS_API_KEY environment variable:")
		VPrintln("export TTS_API_KEY=your_api_key_here")
	} else {
		VPrintf("TTS_API_KEY loaded successfully (length: %d)\n", len(TTS_API_KEY))
	}

	R2_DB_TOKEN = os.Getenv("R2_DB_TOKEN")
	if R2_DB_TOKEN == "" {
		VPrintln("Warning: R2_DB_TOKEN environment variable is not set")
		VPrintln("Please set the R2_DB_TOKEN environment variable:")
		VPrintln("export R2_DB_TOKEN=your_token_here")
	} else {
		VPrintf("R2_DB_TOKEN loaded successfully (length: %d)\n", len(R2_DB_TOKEN))
	}
}
