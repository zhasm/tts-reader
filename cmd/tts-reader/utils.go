package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/zhasm/tts-reader/internal/tts"
	"github.com/zhasm/tts-reader/internal/utils"
	"github.com/zhasm/tts-reader/pkg/logger"
)

// RetryableFunc defines a function that can be retried
type RetryableFunc func(tts.TTSRequest) (bool, error)

// RetryConfig holds configuration for retry behavior
type RetryConfig struct {
	Name         string
	Indent       string
	MaxRetries   int
	InitialDelay time.Duration
}

// WithRetry wraps a function with retry logic and logging
func WithRetry(fn RetryableFunc, name, indent string, maxRetries int) RetryableFunc {
	config := RetryConfig{
		Name:         name,
		Indent:       indent,
		MaxRetries:   maxRetries,
		InitialDelay: time.Second,
	}
	return WithRetryConfig(fn, config)
}

// WithRetryConfig wraps a function with configurable retry logic
func WithRetryConfig(fn RetryableFunc, config RetryConfig) RetryableFunc {
	if config.MaxRetries <= 0 {
		config.MaxRetries = 1
	}
	if config.InitialDelay <= 0 {
		config.InitialDelay = time.Second
	}

	return func(req tts.TTSRequest) (bool, error) {
		delay := config.InitialDelay

		for attempt := 0; attempt < config.MaxRetries; attempt++ {
			result, err := executeAttempt(fn, req, config, attempt)

			// Success case - return immediately
			if err == nil {
				return result, nil
			}

			// Last attempt failed - return the error
			if attempt == config.MaxRetries-1 {
				return result, err
			}

			// Sleep before retry (exponential backoff)
			time.Sleep(delay)
			delay *= 2
		}

		// This should never be reached, but included for completeness
		return false, fmt.Errorf("retry exhausted")
	}
}

// executeAttempt runs a single attempt and handles logging
func executeAttempt(fn RetryableFunc, req tts.TTSRequest, config RetryConfig, attempt int) (bool, error) {
	logger.LogInfo("%s%s begins", config.Indent, config.Name)
	start := time.Now()

	result, err := fn(req)
	duration := time.Since(start).Seconds()

	if err == nil {
		logger.LogInfo("%s%s [%d] succeeded, took %.3f(s)",
			config.Indent, config.Name, attempt, duration)
		return result, nil
	}

	logger.LogInfo("%s%s [%d] failed: %v, took %.3f(s), will retry",
		config.Indent, config.Name, attempt, err, duration)
	return result, err
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
