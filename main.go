package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

func runWithIndent(fn func(TTSRequest) (bool, error), req TTSRequest, depth int, wg *sync.WaitGroup) {
	indent := strings.Repeat("  ", depth)
	functionName := getFuncName(fn)
	LogInfo("%s%s begins", indent, functionName)
	go func() {
		defer wg.Done()
		beginTime := time.Now()
		err := RetryWithBackoff(
			func() error {
				result, err := fn(req)
				if err != nil {
					LogInfo("%s%s ends with error: %v, will retry", indent, functionName, err)
					return err
				}
				LogInfo("%s%s ends, success: %v", indent, functionName, result)
				return nil
			},
			5,             // N: max retries
			2*time.Second, // M: initial interval
		)
		timeCost := time.Since(beginTime)
		if err != nil {
			LogInfo("%s%s failed after retries: %v, took %.3f", indent, functionName, err, timeCost.Seconds())
		} else {
			LogInfo("%s%s succeeded, took %.3f", indent, functionName, timeCost.Seconds())
		}
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
	ok, ttsErr := reqTTS(req)
	if ttsErr != nil {
		fmt.Println("TTS error:", ttsErr)
	}
	if ok {
		content := req.Content
		if len(content) > 64 {
			content = content[:64] + "..."
		}
		LogInfo("%s: [%s]", GetFlag(), content)
		LogInfo("📂: %s", toHomeRelativePath(req.Dest))
		funcs := []func(TTSRequest) (bool, error){
			uploadToR2,
			AppendRecord,
			playAudio,
		}
		var wg sync.WaitGroup
		wg.Add(len(funcs))
		for i, f := range funcs {
			runWithIndent(f, req, i, &wg)
		}
		wg.Wait()
	}
}
