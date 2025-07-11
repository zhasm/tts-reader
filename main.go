package main

import (
	"fmt"
	"strings"
	"sync"
)

func runWithIndent(fn func(TTSRequest) (bool, error), req TTSRequest, depth int, wg *sync.WaitGroup) {
	defer wg.Done()
	indent := strings.Repeat("  ", depth)
	functionName := getFuncName(fn)
	logger.Printf("%s%s begins", indent, functionName)
	// If you want to run nested functions, call runWithIndent with depth+1
	fn(req)
	logger.Printf("%s%s ends", indent, functionName)
}

func main() {

	req := NewTTSRequest(
		"il prenait, ils prenaient",
		"fr-FR",
		"fr-FR-DeniseNeural",
	)

	ok, err := reqTTS(req)
	if err != nil {
		fmt.Println("TTS error:", err)
	}
	if ok {
		funcs := []func(TTSRequest) (bool, error){
			playAudio,
			uploadToR2,
			AppendRecord,
		}
		var wg sync.WaitGroup
		wg.Add(len(funcs))
		for i, f := range funcs {
			go runWithIndent(f, req, i, &wg)
		}
		wg.Wait()
	}
}
