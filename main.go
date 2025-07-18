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
		result, err := fn(req)
		timeCost := time.Since(beginTime)

		if err != nil {
			LogInfo("%s%s ends with error: %v, took %.3f", indent, functionName, err, timeCost.Seconds())
		} else {
			LogInfo("%s%s ends, success: %v, took %.3f", indent, functionName, result, timeCost.Seconds())
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
		content := req.Content
		if len(content) > 64 {
			content = content[:64] + "..."
		}
		LogInfo("%s: [%s]", GetFlag(), content)
		LogInfo("📂: %s", toHomeRelativePath(req.Dest))
		funcs := []func(TTSRequest) (bool, error){
			playAudio,
		}
		if !DryRun {
			funcs = append(funcs, uploadToR2, AppendRecord)
		}
		var wg sync.WaitGroup
		wg.Add(len(funcs))
		for i, f := range funcs {
			runWithIndent(f, req, i, &wg)
		}
		wg.Wait()
		LogInfo("✅ %s [%s]\n", GetFlag(), content)
	}
}
