package main

import (
	"fmt"
	"os"
	"sync"

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

func main() {
	// Initialize logger and config
	logger.Init()
	config.Init()

	config.ValidateAndHandleArgs()

	lang, found := config.GetLang(config.Language)
	if !found {
		fmt.Println("Language not found:", config.Language)
		return
	}

	// generate tts requst struct, no Internet request yet.
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
	defer logger.LogInfo("%s ✅ \n\n", content)
	ok, ttsErr := tts.ReqTTS(req)
	if ttsErr != nil || !ok {
		fmt.Println("TTS error:", ttsErr)
		os.Exit(1)
	}

	logger.LogInfo("📂: %s", utils.ToHomeRelativePath(req.Dest))

	funcs := []func(tts.TTSRequest) (bool, error){
		WithRetry(player.PlayAudio, "main.playAudio", "    ", utils.MAX_RETRY),
	}
	if !config.DryRun {
		funcs = append(funcs, WithRetry(storage.AppendRecord, "main.AppendRecord", "  ", utils.MAX_RETRY))
		funcs = append(funcs, WithRetry(storage.UploadToR2, "main.uploadToR2", "", utils.MAX_RETRY))
	}

	var wg sync.WaitGroup
	wg.Add(len(funcs))
	for _, f := range funcs {
		runWithIndent(f, req, &wg)
	}
	wg.Wait()
}
