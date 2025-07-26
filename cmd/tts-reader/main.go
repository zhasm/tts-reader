package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/zhasm/tts-reader/pkg/config"
)

type TTSRequest struct {
	Language string  `json:"language"`
	Speed    float64 `json:"speed"`
	Content  string  `json:"content"`
}

func ttsHandler(w http.ResponseWriter, r *http.Request) {
	var req TTSRequest

	switch r.Method {
	case http.MethodGet:
		req.Language = r.URL.Query().Get("language")
		req.Content = r.URL.Query().Get("content")
		speedStr := r.URL.Query().Get("speed")
		if speedStr != "" {
			if s, err := strconv.ParseFloat(speedStr, 64); err == nil {
				req.Speed = s
			} else {
				http.Error(w, "Invalid speed parameter", http.StatusBadRequest)
				return
			}
		}
	case http.MethodPost:
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set defaults if not provided
	if req.Language == "" {
		req.Language = config.Language
	}
	if req.Speed == 0 {
		req.Speed = config.Speed
	}
	if req.Content == "" {
		http.Error(w, "Missing content", http.StatusBadRequest)
		return
	}

	// Here you would call your TTS logic, e.g., runTTS(req.Language, req.Speed, req.Content)
	if err := RunWithAPI(req.Language, req.Speed, req.Content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "TTS processed: language=%s, speed=%.2f, content=%q\n", req.Language, req.Speed, req.Content)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		http.HandleFunc("/tts", ttsHandler)
		addr := "0.0.0.0:8080"
		fmt.Println("Listening on", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		return
	}

	if err := run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
