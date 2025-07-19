package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Add a retry wrapper function
type retryableFunc func(TTSRequest) (bool, error)

func withRetry(fn retryableFunc, name, indent string, maxRetries int) retryableFunc {
	return func(req TTSRequest) (bool, error) {
		var err error
		var ok bool
		for retryIndex := 0; retryIndex < maxRetries; retryIndex++ {
			LogInfo("%s%s begins", indent, name)
			beginTime := time.Now()
			ok, err = fn(req)
			timeCost := time.Since(beginTime)
			if err == nil {
				LogInfo("%s%s [%d]succeeded, took %.3f(s)", indent, name, retryIndex, timeCost.Seconds())
				return ok, nil
			}
			LogInfo("%s%s [%d]failed: %v, took %.3f(s), will retry", indent, name, retryIndex, err, timeCost.Seconds())
		}
		return ok, err
	}
}

func runWithIndent(fn func(TTSRequest) (bool, error), req TTSRequest, depth int, wg *sync.WaitGroup) {
	functionName := getFuncName(fn)
	go func() {
		defer wg.Done()
		_, err := fn(req)

		if err != nil {
			LogInfo("%s ends with error: %v", functionName, err)
			return
		}
		//		LogInfo("%s ends, success: %v", functionName, result)
	}()
}

func RetryWithBackoff(fn func() error, maxRetries int, initialInterval time.Duration) error {
	interval := initialInterval
	var lastErr error
	for i := 0; i < maxRetries; i++ {
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
	lang, found := GetLang(Language)
	if !found {
		fmt.Println("Language not found:", Language)
		return
	}
	req := NewTTSRequest(
		Content,
		lang.NameFUll,
		lang.Reader,
	)
	content := req.Content
	LogInfo("%s: [%s]", GetFlag(), content)
	ok, ttsErr := reqTTS(req)
	if ttsErr != nil || !ok {
		fmt.Println("TTS error:", ttsErr)
		os.Exit(1)
	}

	if len(content) > 64 {
		content = content[:64] + "..."
	}
	LogInfo("📂: %s", toHomeRelativePath(req.Dest))
	maxRetries := 10
	funcs := []func(TTSRequest) (bool, error){
		withRetry(playAudio, "main.playAudio", "    ", maxRetries),
	}
	if !DryRun {
		funcs = append(funcs, withRetry(AppendRecord, "main.AppendRecord", "  ", maxRetries))
		funcs = append(funcs, withRetry(uploadToR2, "main.uploadToR2", "", maxRetries))
	}

	var wg sync.WaitGroup
	wg.Add(len(funcs))
	for i, f := range funcs {
		runWithIndent(f, req, i, &wg) // Pass retryIndex (i)
	}
	wg.Wait()
	LogInfo("✅ %s: [%s]\n", GetFlag(), content)
}
