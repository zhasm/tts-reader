package storage

import (
	"testing"

	"github.com/zhasm/tts-reader/internal/tts"
)

func TestAppendRecord_UnsupportedLang(t *testing.T) {
	req := tts.TTSRequest{Lang: "xx", Content: "test", Dest: "dummy.mp3", Md5: "abc"}
	ok, err := AppendRecord(req)
	if ok || err == nil {
		t.Errorf("Expected error for unsupported language, got ok=%v, err=%v", ok, err)
	}
}

func TestAppendRecord_FileNotExist(t *testing.T) {
	req := tts.TTSRequest{Lang: "fr", Content: "test", Dest: "nonexistent.mp3", Md5: "abc"}
	ok, err := AppendRecord(req)
	if ok || err == nil {
		t.Errorf("Expected error for non-existent file, got ok=%v, err=%v", ok, err)
	}
}
