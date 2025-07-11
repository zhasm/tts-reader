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
