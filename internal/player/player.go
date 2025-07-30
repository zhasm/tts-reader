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
		logger.LogError("Warning: Audio file does not exist: %s", file)
		logger.LogError("Error: %v", err)
		return false, err
	}

	// Check if file has minimum size (not empty/corrupted)
	if fileInfo, err := os.Stat(file); err == nil {
		if fileInfo.Size() < 1000 {
			logger.LogWarn("Warning: Audio file appears to be corrupted or empty (size: %d bytes)", fileInfo.Size())
			return false, fmt.Errorf("file too small: %d bytes", fileInfo.Size())
		}
	}
	return true, nil
}

func PlayAudio(req tts.TTSRequest) (bool, error) {
	file := req.Dest
	// Check if file exists and is valid
	if valid, err := isAudioFileValid(file); !valid {
		logger.LogError("Audio file validation failed: %v", err)
		return false, err
	}

	// Play audio with ffplay command in background
	logger.LogDebug("Playing audio: %s", file)
	cmd := exec.Command("ffplay", "-hide_banner", "-loglevel", "panic", "-nodisp", "-autoexit", file)

	// Start the command in background (non-blocking)
	if err := cmd.Start(); err != nil {
		logger.LogError("Error starting audio playback: %v", err)
		logger.LogError("Make sure ffplay is installed and available in PATH")
		return false, err
	}

	logger.LogDebug("Audio playback started in background")
	return true, nil
}
