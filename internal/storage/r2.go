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
		logger.LogWarn("File does not exist: %s", filename)
		return false, err
	}

	if fileInfo.Size() == 0 {
		logger.LogWarn("File is empty: %s", filename)
		return false, fmt.Errorf("file is empty")
	}

	// Check if rclone is available
	if _, err := exec.LookPath("rclone"); err != nil {
		logger.LogError("rclone not found in PATH: %v", err)
		return false, err
	}

	// Upload file to R2
	logger.LogDebug("Uploading %s to R2...", filename)
	cmd := exec.Command("rclone", "copy", filename, "r2:tts/")

	uploadErr := utils.RetryWithBackoff(func() error {
		err := cmd.Run()
		if err != nil {
			logger.LogWarn("Upload failed: %v", err)
		}
		return err
	}, utils.MAX_RETRY, 1*time.Second)
	if uploadErr != nil {
		logger.LogError("Upload failed after retries: %v", uploadErr)
		return false, uploadErr
	}

	logger.LogDebug("Successfully uploaded %s to R2", filename)

	//	copy mp3 url to clipboard
	url := fmt.Sprintf("%s/%s.mp3", R2_URL_PREFIX, req.Md5)
	cmd = exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(url)
	if err := cmd.Run(); err != nil {
		logger.LogWarn("copy url to clipboard error: %v", err)
		return false, err
	}

	return true, nil
}
