package main

import (
	"fmt"
	"os"
	"os/exec"
)

func isAudioFileValid(file string) (bool, error) {
	if _, err := os.Stat(file); err != nil {
		VPrintf("Warning: Audio file does not exist: %s\n", file)
		VPrintf("Error: %v\n", err)
		return false, err
	}

	// Check if file has minimum size (not empty/corrupted)
	if fileInfo, err := os.Stat(file); err == nil {
		if fileInfo.Size() < 1000 {
			VPrintf("Warning: Audio file appears to be corrupted or empty (size: %d bytes)\n", fileInfo.Size())
			return false, fmt.Errorf("file too small: %d bytes", fileInfo.Size())
		}
	}
	return true, nil
}

func playAudio(req TTSRequest) (bool, error) {
	file := req.Dest
	// Check if file exists and is valid
	if valid, err := isAudioFileValid(file); !valid {
		VPrintf("Audio file validation failed: %v\n", err)
		return false, err
	}

	// Play audio with ffplay command
	VPrintf("Playing audio: %s\n", file)
	cmd := exec.Command("ffplay", "-hide_banner", "-loglevel", "panic", "-nodisp", "-autoexit", file)

	// Run the command
	if err := cmd.Run(); err != nil {
		VPrintf("Error playing audio: %v\n", err)
		VPrintf("Make sure ffplay is installed and available in PATH\n")
		return false, err
	}

	VPrintf("Audio playback completed\n")
	return true, nil
}
