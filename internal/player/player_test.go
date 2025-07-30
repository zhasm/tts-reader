package player

import (
	"os"
	"testing"

	"github.com/zhasm/tts-reader/internal/tts"
)

func TestIsAudioFileValid(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "audio-*.mp3")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Write enough data to be valid (2000 bytes)
	if _, err := tmpfile.Write(make([]byte, 2000)); err != nil {
		t.Fatalf("Failed to write to tmpfile: %v", err)
	}
	tmpfile.Close()

	valid, err := tts.IsAudioFileValid(tmpfile.Name())
	if !valid || err != nil {
		t.Errorf("Expected valid audio file, got valid=%v, err=%v", valid, err)
	}

	// Test with too small file
	tmpfile2, _ := os.CreateTemp("", "audio-*.mp3")
	defer os.Remove(tmpfile2.Name())
	if _, err := tmpfile2.Write([]byte("a")); err != nil {
		t.Fatalf("Failed to write to tmpfile2: %v", err)
	}
	tmpfile2.Close()
	valid, err = tts.IsAudioFileValid(tmpfile2.Name())
	if valid || err == nil {
		t.Errorf("Expected invalid audio file, got valid=%v, err=%v", valid, err)
	}

	// Test with non-existent file
	valid, err = tts.IsAudioFileValid("nonexistent.mp3")
	if valid || err == nil {
		t.Errorf("Expected error for non-existent file, got valid=%v, err=%v", valid, err)
	}
}

func TestPlayAudio_NonExistent(t *testing.T) {
	req := tts.TTSRequest{Dest: "nonexistent.mp3"}
	ok, err := PlayAudio(req)
	if ok || err == nil {
		t.Errorf("Expected error for non-existent file, got ok=%v, err=%v", ok, err)
	}
}
