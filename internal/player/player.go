package player

import (
	"os/exec"

	"github.com/zhasm/tts-reader/internal/tts"
	"github.com/zhasm/tts-reader/pkg/logger"
)

func PlayAudio(req tts.TTSRequest) (bool, error) {
	file := req.Dest
	// Check if file exists and is valid
	if valid, err := tts.IsAudioFileValid(file); !valid {
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
