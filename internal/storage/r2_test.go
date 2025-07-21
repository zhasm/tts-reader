package storage

import (
	"testing"

	"github.com/zhasm/tts-reader/internal/tts"
)

func TestUploadToR2_FileNotExist(t *testing.T) {
	req := tts.TTSRequest{Dest: "nonexistent.mp3", Md5: "abc"}
	ok, err := UploadToR2(req)
	if ok || err == nil {
		t.Errorf("Expected error for non-existent file, got ok=%v, err=%v", ok, err)
	}
}
