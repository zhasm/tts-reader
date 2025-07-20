package player

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/zhasm/tts-reader/internal/tts"
	"github.com/zhasm/tts-reader/pkg/logger"
)

func isAudioFileValid(file string) (bool, error) {
	if _, err := os.Stat(file); err != nil {
		logger.VPrintf("Warning: Audio file does not exist: %s\n", file)
		logger.VPrintf("Error: %v\n", err)
		return false, err
	}

	// Check if file has minimum size (not empty/corrupted)
	if fileInfo, err := os.Stat(file); err == nil {
		if fileInfo.Size() < 1000 {
			logger.VPrintf("Warning: Audio file appears to be corrupted or empty (size: %d bytes)\n", fileInfo.Size())
			return false, fmt.Errorf("file too small: %d bytes", fileInfo.Size())
		}
	}
	return true, nil
}

func PlayAudio(req tts.TTSRequest) (bool, error) {
	file := req.Dest
	// Check if file exists and is valid
	if valid, err := isAudioFileValid(file); !valid {
		logger.VPrintf("Audio file validation failed: %v\n", err)
		return false, err
	}

	// Play audio with ffplay command in background
	logger.VPrintf("Playing audio: %s\n", file)
	cmd := exec.Command("ffplay", "-hide_banner", "-loglevel", "panic", "-nodisp", "-autoexit", file)

	// Start the command in background (non-blocking)
	if err := cmd.Start(); err != nil {
		logger.VPrintf("Error starting audio playback: %v\n", err)
		logger.VPrintf("Make sure ffplay is installed and available in PATH\n")
		return false, err
	}

	logger.VPrintf("Audio playback started in background\n")
	return true, nil
}
