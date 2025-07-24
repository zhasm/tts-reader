package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/zhasm/tts-reader/internal/player"
	"github.com/zhasm/tts-reader/internal/storage"
	"github.com/zhasm/tts-reader/internal/tts"
	"github.com/zhasm/tts-reader/internal/utils"
	"github.com/zhasm/tts-reader/pkg/config"
	"github.com/zhasm/tts-reader/pkg/logger"
)

func run() error {
	initLoggerAndConfig()

	if err := config.ValidateAndHandleArgs(); err != nil {
		return fmt.Errorf("argument validation failed: %w", err)
	}

	lang, found := config.GetLang(config.Language)
	if !found {
		return fmt.Errorf("language not found: %s", config.Language)
	}

	req := createTTSRequest(lang)
	logContentPreview(req)

	if ok, err := tts.ReqTTS(req); err != nil || !ok {
		return fmt.Errorf("TTS request failed: %w", err)
	}

	logger.LogInfo("📂: %s", utils.ToHomeRelativePath(req.Dest))

	funcs := buildProcessingPipeline()

	return runFunctionsConcurrently(funcs, req)
}

func initLoggerAndConfig() {
	logger.Init()
	config.Init()
}

func createTTSRequest(lang config.Lang) tts.TTSRequest {
	return tts.NewTTSRequest(
		config.Content,
		lang.NameFUll,
		lang.Reader,
		config.Speed,
	)
}

func logContentPreview(req tts.TTSRequest) string {
	content := req.Content
	contentLen := len(content)

	if contentLen > MAX_CONTENT_LENGTH_TO_SHOW {
		content = content[:MAX_CONTENT_LENGTH_TO_SHOW] + "..."
	}
	//	content =
	return fmt.Sprintf("%s [%s][%d]", config.GetFlagByName(config.Language), content, contentLen)

}

func buildProcessingPipeline() []func(tts.TTSRequest) (bool, error) {
	funcs := []func(tts.TTSRequest) (bool, error){
		player.PlayAudio,
	}

	if !config.DryRun {
		funcs = append([]func(tts.TTSRequest) (bool, error){
			storage.UploadToR2,
			storage.AppendRecord,
		}, funcs...)
	}
	return funcs
}

func runFunctionsConcurrently(funcs []func(tts.TTSRequest) (bool, error), req tts.TTSRequest) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(funcs))
	content := logContentPreview(req)

	logger.LogInfo("%s ⏰", content)
	defer logger.LogInfo("%s ✅ \n\n", content)

	// Function name mapping for logging - matches the expected log output
	funcNames := []string{
		"main.uploadToR2",
		"main.AppendRecord",
		"main.playAudio",
	}

	// Indentation levels based on call hierarchy - matches expected nesting
	indentLevels := []string{
		"  ",     // main.uploadToR2 - no indent
		"    ",   // main.AppendRecord - 2 spaces
		"      ", // main.playAudio - 4 spaces
	}

	wg.Add(len(funcs))
	for i, f := range funcs {
		funcName := funcNames[i]
		indent := ""
		if i < len(indentLevels) {
			indent = indentLevels[i]
		}

		// Log function start
		logger.LogInfo("%s%s begins", indent, funcName)
		go func(i int, f func(tts.TTSRequest) (bool, error)) {
			defer wg.Done()

			start := time.Now()

			ok, err := f(req)

			// Calculate duration
			duration := time.Since(start).Seconds()

			if err != nil || !ok {
				logger.LogInfo("%s%s [%d] failed, took %.3f(s)", indent, funcName, i, duration)
				errChan <- fmt.Errorf("function %d failed: %w", i, err)
			} else {
				logger.LogInfo("%s%s [%d] succeeded, took %.3f(s)", indent, funcName, i, duration)
			}
		}(i, f)
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}
