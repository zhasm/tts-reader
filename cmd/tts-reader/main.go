package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/zhasm/tts-reader/internal/player"
	"github.com/zhasm/tts-reader/internal/storage"
	"github.com/zhasm/tts-reader/internal/tts"
	"github.com/zhasm/tts-reader/internal/utils"
	"github.com/zhasm/tts-reader/pkg/config"
	"github.com/zhasm/tts-reader/pkg/logger"
)

const (
	MAX_CONTENT_LENGTH_TO_SHOW = 42
)

// Add a retry wrapper function
type retryableFunc func(tts.TTSRequest) (bool, error)

func withRetry(fn retryableFunc, name, indent string, maxRetries int) retryableFunc {
	return func(req tts.TTSRequest) (bool, error) {
		var err error
		var ok bool
		for retryIndex := range maxRetries {
			logger.LogInfo("%s%s begins", indent, name)
			beginTime := time.Now()
			ok, err = fn(req)
			timeCost := time.Since(beginTime)
			if err == nil {
				logger.LogInfo("%s%s [%d]succeeded, took %.3f(s)", indent, name, retryIndex, timeCost.Seconds())
				return ok, nil
			}
			logger.LogInfo("%s%s [%d]failed: %v, took %.3f(s), will retry", indent, name, retryIndex, err, timeCost.Seconds())
		}
		return ok, err
	}
}

func runWithIndent(fn func(tts.TTSRequest) (bool, error), req tts.TTSRequest, wg *sync.WaitGroup) {
	functionName := utils.GetFuncName(fn)
	go func() {
		defer wg.Done()
		_, err := fn(req)

		if err != nil {
			logger.LogInfo("%s ends with error: %v", functionName, err)
			return
		}
		//		logger.LogInfo("%s ends, success: %v", functionName, result)
	}()
}

func RetryWithBackoff(fn func() error, maxRetries int, initialInterval time.Duration) error {
	interval := initialInterval
	var lastErr error
	for i := range maxRetries {
		err := fn()
		if err == nil {
			return nil
		}
		lastErr = err
		if i < maxRetries-1 {
			time.Sleep(interval)
			interval *= 2
		}
	}
	return lastErr
}

func main() {
	// Initialize logger and config
	logger.Init()
	config.Init()

	lang, found := config.GetLang(config.Language)
	if !found {
		fmt.Println("Language not found:", config.Language)
		return
	}
	req := tts.NewTTSRequest(
		config.Content,
		lang.NameFUll,
		lang.Reader,
	)
	content := req.Content
	if len(content) > MAX_CONTENT_LENGTH_TO_SHOW {
		content = content[:MAX_CONTENT_LENGTH_TO_SHOW]
	}
	content = fmt.Sprintf("%s %s", config.GetFlag(), content)
	logger.LogInfo("%s ...", content)
	ok, ttsErr := tts.ReqTTS(req)
	if ttsErr != nil || !ok {
		fmt.Println("TTS error:", ttsErr)
		os.Exit(1)
	}

	logger.LogInfo("📂: %s", utils.ToHomeRelativePath(req.Dest))
	maxRetries := 10
	funcs := []func(tts.TTSRequest) (bool, error){
		withRetry(player.PlayAudio, "main.playAudio", "    ", maxRetries),
	}
	if !config.DryRun {
		funcs = append(funcs, withRetry(storage.AppendRecord, "main.AppendRecord", "  ", maxRetries))
		funcs = append(funcs, withRetry(storage.UploadToR2, "main.uploadToR2", "", maxRetries))
	}

	var wg sync.WaitGroup
	wg.Add(len(funcs))
	for _, f := range funcs {
		runWithIndent(f, req, &wg)
	}
	wg.Wait()
	logger.LogInfo("%s ✅ \n\n", content)
}
