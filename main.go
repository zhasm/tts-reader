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
	logger.Printf("%s%s begins", indent, functionName)
	go func() {
		defer wg.Done()
		beginTime := time.Now()
		result, err := fn(req)
		timeCost := time.Since(beginTime)

		if err != nil {
			logger.Printf("%s%s ends with error: %v, took %.3f", indent, functionName, err, timeCost.Seconds())
		} else {
			logger.Printf("%s%s ends, success: %v, took %.3f", indent, functionName, result, timeCost.Seconds())
		}
	}()
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
		logger.Printf("📂: %s", req.Dest)
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
