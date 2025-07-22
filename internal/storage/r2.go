package storage

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/zhasm/tts-reader/internal/tts"
	"github.com/zhasm/tts-reader/internal/utils"
	"github.com/zhasm/tts-reader/pkg/logger"
)

const (
	R2_URL_PREFIX = "https://pub-c6b11003307646e98afc7540d5f09c41.r2.dev"
)

func UploadToR2(req tts.TTSRequest) (bool, error) {
	// Check if file exists and is not empty
	filename := req.Dest
	fileInfo, err := os.Stat(filename)
	if err != nil {
		logger.VPrintf("File does not exist: %s\n", filename)
		return false, err
	}

	if fileInfo.Size() == 0 {
		logger.VPrintf("File is empty: %s\n", filename)
		return false, fmt.Errorf("file is empty")
	}

	// Check if rclone is available
	if _, err := exec.LookPath("rclone"); err != nil {
		logger.VPrintf("rclone not found in PATH: %v\n", err)
		return false, err
	}

	// Upload file to R2
	logger.VPrintf("Uploading %s to R2...\n", filename)
	cmd := exec.Command("rclone", "copy", filename, "r2:tts/")

	uploadErr := utils.RetryWithBackoff(func() error {
		err := cmd.Run()
		if err != nil {
			logger.VPrintf("Upload failed: %v\n", err)
		}
		return err
	}, 10, 1*time.Second)
	if uploadErr != nil {
		logger.VPrintf("Upload failed after retries: %v\n", uploadErr)
		return false, uploadErr
	}

	logger.VPrintf("Successfully uploaded %s to R2\n", filename)

	//	copy mp3 url to clipboard
	url := fmt.Sprintf("%s/%s.mp3", R2_URL_PREFIX, req.Md5)
	cmd = exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(url)
	if err := cmd.Run(); err != nil {
		logger.VPrintf("copy url to clipboard error: %v\n", err)
		return false, err
	}

	return true, nil
}
