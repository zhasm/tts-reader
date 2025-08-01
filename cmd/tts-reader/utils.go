package main

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/zhasm/tts-reader/internal/player"
	"github.com/zhasm/tts-reader/internal/storage"
	"github.com/zhasm/tts-reader/internal/tts"
	"github.com/zhasm/tts-reader/internal/utils"
	"github.com/zhasm/tts-reader/pkg/config"
	"github.com/zhasm/tts-reader/pkg/logger"
	"golang.org/x/term"
)

const (
	INDENT_DEFAULT = 23
)

func run() error {
	startAll := time.Now()
	initLoggerAndConfig()
	if err := config.ValidateAndHandleArgs(); err != nil {
		return fmt.Errorf("argument validation failed: %w", err)
	}

	lang, found := config.GetLang(config.Language)
	if !found {
		return fmt.Errorf("language not found: %s", config.Language)
	}

	req := createTTSRequest(lang)
	content := logContentPreview(req)

	if ok, err := config.ValidateLangRegex(config.Language, config.Content); err != nil {
		return fmt.Errorf("language validation failed: %w", err)
	} else if !ok {
		return fmt.Errorf("language validation failed: %s", config.Language)
	}

	logger.LogInfo("%s", MsgWithIcon(content, "⏰"))
	logger.LogInfo("📂: %s", utils.ToHomeRelativePath(req.Dest))
	logger.LogInfo("⌛️ TTS request in progress...")

	defer func() {
		logger.LogInfo("%s", MsgWithIcon(content, "✅"))
		logger.LogInfo("Total time taken: %.3f(s)\n", time.Since(startAll).Seconds())
	}()

	start := time.Now()
	if ok, err := tts.ReqTTS(req); err != nil || !ok {
		return fmt.Errorf("TTS request failed: %w", err)
	}
	duration := time.Since(start).Seconds()
	logger.LogInfo("✅ TTS request completed, took %.3f(s)", duration)

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

	wg.Add(len(funcs))
	for i, f := range funcs {
		funcName := GetFuncName(f)
		indent := strings.Repeat("  ", i) + "  " // 2 spaces per level

		// Log function start
		logger.LogInfo("%s%s begins", indent, funcName)
		go func(i int, f func(tts.TTSRequest) (bool, error), funcName, indent string) {
			defer wg.Done()

			start := time.Now()

			ok, err := f(req)

			// Calculate duration
			duration := time.Since(start).Seconds()

			if err != nil || !ok {
				logger.LogWarn("%s%s [%d] failed, took %.3f(s)", indent, funcName, i, duration)
				errChan <- fmt.Errorf("function %d failed: %w", i, err)
			} else {
				logger.LogInfo("%s%s succeeded, took %.3f(s)", indent, funcName, duration)
			}
		}(i, f, funcName, indent)
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

func GetWindowWidth() (int, error) {
	fd := int(os.Stdout.Fd())

	// Check if the file descriptor refers to a terminal.
	if !term.IsTerminal(fd) {
		return 0, fmt.Errorf("not running in a terminal")
	}

	// Get the terminal size.
	width, _, err := term.GetSize(fd)
	if err != nil {
		fmt.Printf("Error getting terminal size: %v\n", err)
		return 0, err
	}

	return width, nil
}

func MsgWithIcon(content, icon string) string {
	defaultStr := content + " " + icon
	width, err := GetWindowWidth()
	if err != nil {
		return defaultStr
	}

	n := max(width-len(content)-INDENT_DEFAULT, 0)
	spaces := strings.Repeat(" ", n)
	return fmt.Sprintf("%s%s%s", content, spaces, icon)
}

func GetFuncName(i any) string {
	ret := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	if !strings.Contains(ret, "/") {
		return ret
	} else {
		segments := strings.Split(ret, "/")
		return segments[len(segments)-1]
	}
}
