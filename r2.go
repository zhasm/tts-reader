package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func uploadToR2(req TTSRequest) (bool, error) {
	// Check if file exists and is not empty
	filename := req.Dest
	fileInfo, err := os.Stat(filename)
	if err != nil {
		VPrintf("File does not exist: %s\n", filename)
		return false, err
	}

	if fileInfo.Size() == 0 {
		VPrintf("File is empty: %s\n", filename)
		return false, fmt.Errorf("file is empty")
	}

	// Check if rclone is available
	if _, err := exec.LookPath("rclone"); err != nil {
		VPrintf("rclone not found in PATH: %v\n", err)
		return false, err
	}

	// Upload file to R2
	VPrintf("Uploading %s to R2...\n", filename)
	cmd := exec.Command("rclone", "copy", filename, "r2:tts/")

	if err := cmd.Run(); err != nil {
		VPrintf("Upload failed: %v\n", err)
		return false, err
	}

	VPrintf("Successfully uploaded %s to R2\n", filename)

	//	copy mp3 url to clipboard
	url := fmt.Sprintf("https://pub-c6b11003307646e98afc7540d5f09c41.r2.dev/%s.mp3", req.Md5)
	cmd = exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(url)
	if err := cmd.Run(); err != nil {
		VPrintf("copy url to clipboard error: %v\n", err)
		return false, err
	}

	return true, nil
}
