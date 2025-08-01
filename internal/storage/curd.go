package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/zhasm/tts-reader/internal/tts"
	"github.com/zhasm/tts-reader/internal/utils"
	"github.com/zhasm/tts-reader/pkg/config"
	"github.com/zhasm/tts-reader/pkg/logger"
)

const (
	CRUD_HOST = "https://tts-server.rex-zhasm6886.workers.dev/api/item"
)

func AppendRecord(req tts.TTSRequest) (bool, error) {

	// Normalize language code
	lang := ""
	for _, l := range config.Langs {
		if req.Lang == l.Name || req.Lang == l.NameFUll {
			lang = l.Name
			break
		}
	}
	if lang == "" {
		logger.LogError("Unsupported language: %s", req.Lang)
		return false, fmt.Errorf("unsupported language: %s", req.Lang)
	}

	// Get file size in KB
	fileInfo, err := os.Stat(req.Dest)
	if err != nil {
		logger.LogError("Error getting file info: %v", err)
		return false, err
	}
	fileSizeKb := fmt.Sprintf("%d", fileInfo.Size()/1024)
	// Build JSON data
	data := map[string]string{
		"language":   lang,
		"content":    req.Content,
		"FileSizeKb": fileSizeKb,
		"md5":        req.Md5,
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		logger.LogError("Error marshaling JSON: %v", err)
		return false, err
	}
	body := bytes.NewBuffer(jsonBytes)

	// Create client
	client := &http.Client{}

	httpHeaders := map[string]string{
		"Content-Type": "application/json",
	}
	httpReq, err := utils.NewHTTPRequest("POST", CRUD_HOST, body, httpHeaders)
	if err != nil {
		logger.LogError("Error creating HTTP request: %v", err)
		return false, err
	}

	// Headers
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer "+config.R2_DB_TOKEN)

	// Fetch Request
	resp, err := utils.HTTPRequest(client, httpReq)
	if err != nil {
		logger.LogError("Failure : ", err)
		return false, err
	}
	defer resp.Body.Close()

	// Read Response Body
	respBody, _ := io.ReadAll(resp.Body)

	// Display Results
	logger.LogDebug("response Status : ", resp.Status)
	logger.LogDebug("response Headers : ", resp.Header)
	logger.LogDebug("response Body : ", string(respBody))
	if resp.StatusCode != 200 {
		fmt.Println("Appending record Error!")
		os.Exit(1)
	}
	return true, nil
}
