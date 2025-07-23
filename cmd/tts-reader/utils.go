package main

import (
	"sync"
	"time"

	"github.com/zhasm/tts-reader/internal/tts"
	"github.com/zhasm/tts-reader/internal/utils"
	"github.com/zhasm/tts-reader/pkg/logger"
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
