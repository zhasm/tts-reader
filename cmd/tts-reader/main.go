package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/spf13/pflag"
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
		interval := time.Second
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
			if retryIndex < maxRetries-1 {
				time.Sleep(interval)
				interval *= 2
			}
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

	if config.Version {
		config.PrintVersion()
	}

	if config.Help {
		config.PrintHelp(0)
	}

	// If no content is provided, we should exit, unless it's a "meta" command
	// that doesn't need content. Help and Version are already handled.
	if config.Content == "" {
		// If no arguments at all were provided, show help.
		if len(os.Args) == 1 {
			config.PrintHelp(0)
		}

		// A flag was passed, but no content. This is only ok for -v, -h, -V.
		// -h and -V are already handled. So if -v is the *only* flag, it's ok.
		isOnlyVerbose := false
		if pflag.NFlag() == 1 {
			pflag.Visit(func(f *pflag.Flag) {
				if f.Name == "verbose" {
					isOnlyVerbose = true
				}
			})
		}

		if isOnlyVerbose {
			os.Exit(0) // Successfully do nothing.
		}

		fmt.Fprintln(os.Stderr, "Error: content argument is missing.")
		config.PrintHelp(1)
	}

	lang, found := config.GetLang(config.Language)
	if !found {
		fmt.Println("Language not found:", config.Language)
		return
	}
	req := tts.NewTTSRequest(
		config.Content,
		lang.NameFUll,
		lang.Reader,
		config.Speed,
	)
	content := req.Content
	contentLen := len(content)

	if len(content) > MAX_CONTENT_LENGTH_TO_SHOW {
		content = content[:MAX_CONTENT_LENGTH_TO_SHOW] + "..."
	}
	content = fmt.Sprintf("%s [%s][%d]", config.GetFlagByName(config.Language), content, contentLen)
	logger.LogInfo("%s ⏰", content)
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
