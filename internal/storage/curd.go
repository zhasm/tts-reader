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
		logger.VPrintf("Unsupported language: %s\n", req.Lang)
		return false, fmt.Errorf("unsupported language: %s", req.Lang)
	}

	// Get file size in KB
	fileInfo, err := os.Stat(req.Dest)
	if err != nil {
		logger.VPrintf("Error getting file info: %v\n", err)
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
		logger.VPrintf("Error marshaling JSON: %v\n", err)
		return false, err
	}
	body := bytes.NewBuffer(jsonBytes)

	// Create client
	client := &http.Client{}

	httpHeaders := map[string]string{
		"Content-Type": "application/json",
	}
	httpReq, err := utils.NewHTTPRequestWithRetry("POST", CRUD_HOST, body, httpHeaders)
	if err != nil {
		logger.VPrintf("Error creating HTTP request: %v\n", err)
		return false, err
	}

	// Headers
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer "+config.R2_DB_TOKEN)

	// Fetch Request
	resp, err := client.Do(httpReq)
	if err != nil {
		logger.VPrintln("Failure : ", err)
		return false, err
	}
	defer resp.Body.Close()

	// Read Response Body
	respBody, _ := io.ReadAll(resp.Body)

	// Display Results
	logger.VPrintln("response Status : ", resp.Status)
	logger.VPrintln("response Headers : ", resp.Header)
	logger.VPrintln("response Body : ", string(respBody))
	if resp.StatusCode != 200 {
		fmt.Println("Appending record Error!")
		os.Exit(1)
	}
	return true, nil
}
