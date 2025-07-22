package tts

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/zhasm/tts-reader/internal/utils"
	"github.com/zhasm/tts-reader/pkg/config"
	"github.com/zhasm/tts-reader/pkg/logger"
)

type TTSRequest struct {
	Content string
	Lang    string
	Reader  string
	Speed   float64
	Gender  string
	Dest    string // the output path
	Md5     string
}

const (
	USER_AGENT                 = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36"
	X_MICROSOFT_OUTPUTFORMAT   = "riff-24khz-16bit-mono-pcm"
	HTTP_REQEUEST_HOST         = "westus.tts.speech.microsoft.com"
	HTTP_REQEUEST_CONTENT_TYPE = "application/ssml+xml"
	HTTP_REQEUEST_API          = "https://eastasia.tts.speech.microsoft.com/cognitiveservices/v1"
)

func NewTTSRequest(content, lang, reader string, speed float64) TTSRequest {
	gender := "Male" // default gender

	// Create a unique key based on parameters
	// the ending '\n' is on purpose, please do not delete.
	keyData := fmt.Sprintf("%s-%s-%s-%s-%.1f\n", lang, reader, gender, content, speed)
	logger.VPrintf("DEBUG: Key data string: '%s'\n", keyData)
	key := fmt.Sprintf("%x", md5.Sum([]byte(keyData)))
	logger.VPrintf("DEBUG: Generated MD5: %s\n", key)

	// Create destination path
	dest := fmt.Sprintf("%s/%s.mp3", config.TTS_PATH, key)

	return TTSRequest{
		Content: content,
		Lang:    lang,
		Reader:  reader,
		Speed:   speed,
		Gender:  gender,
		Dest:    dest,
		Md5:     key,
	}
}

func ReqTTS(req TTSRequest) (bool, error) {

	// Check if destination file already exists and is valid
	if valid, _ := isAudioFileValid(req.Dest); valid {
		// Get file info for logging
		if fileInfo, err := os.Stat(req.Dest); err == nil {
			logger.VPrintf("File already exists: %s (size: %d bytes)\n", req.Dest, fileInfo.Size())
		}
		return true, nil
	}

	logger.VPrintf("Content: %s\n", req.Content)
	logger.VPrintf("Language: %s\n", req.Lang)
	logger.VPrintf("Reader: %s\n", req.Reader)
	logger.VPrintf("Speed: %f\n", req.Speed)
	logger.VPrintf("API Key set: %t\n", config.TTS_API_KEY != "")

	// cURL (POST https://eastasia.tts.speech.microsoft.com/cognitiveservices/v1)
	ssmlBody := fmt.Sprintf(`
		<speak version="1.0" xml:lang="%s">
		<voice xml:lang="%s" xml:gender="Male" name="%s">
		<prosody rate="%f">%s</prosody>
		</voice>
		</speak>`, req.Lang, req.Lang, req.Reader, req.Speed, req.Content)

	logger.VPrintf("Generated SSML:\n%s\n", ssmlBody)

	body := strings.NewReader(ssmlBody)

	// Create client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	httpHeaders := map[string]string{
		"X-Microsoft-Outputformat":  X_MICROSOFT_OUTPUTFORMAT,
		"Content-Type":              HTTP_REQEUEST_CONTENT_TYPE,
		"Host":                      HTTP_REQEUEST_HOST,
		"Ocp-Apim-Subscription-Key": config.TTS_API_KEY,
		"User-Agent":                USER_AGENT,
	}
	httpReq, err := utils.NewHTTPRequestWithRetry("POST", HTTP_REQEUEST_API, body, httpHeaders)
	if err != nil {
		logger.VPrintf("Error creating request: %v\n", err)
		return false, err
	}

	logger.VPrintf("Request URL: %s\n", httpReq.URL.String())
	logger.VPrintf("Request Method: %s\n", httpReq.Method)
	logger.VPrintf("Request Headers:\n")
	for key, values := range httpReq.Header {
		for _, value := range values {
			if key == "Ocp-Apim-Subscription-Key" {
				logger.VPrintf("  %s: [HIDDEN]\n", key)
			} else {
				logger.VPrintf("  %s: %s\n", key, value)
			}
		}
	}

	// Fetch Request
	logger.VPrintf("Sending request...\n")
	resp, err := client.Do(httpReq)

	if err != nil {
		logger.VPrintf("HTTP request failed: %v\n", err)
		return false, err
	}

	// Check if response is nil
	if resp == nil {
		logger.VPrintf("Error: Response is nil\n")
		return false, fmt.Errorf("response is nil")
	}

	defer resp.Body.Close()

	// Read Response Body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.VPrintf("Error reading response body: %v\n", err)
		return false, err
	}

	if resp.StatusCode != 200 {
		fmt.Println("Requesting TTS Error!")
		os.Exit(1)
	}

	// Display Results
	logger.VPrintf("=== Response Details ===\n")
	logger.VPrintf("Response Status: %s\n", resp.Status)
	logger.VPrintf("Response Headers:\n")
	for key, values := range resp.Header {
		for _, value := range values {
			logger.VPrintf("  %s: %s\n", key, value)
		}
	}
	logger.VPrintf("Response Body Length: %d bytes\n", len(respBody))
	if len(respBody) < 1000 {
		logger.VPrintf("Response Body: %s\n", string(respBody))
	} else {
		logger.VPrintf("Response Body (first 1000 chars): %s...\n", string(respBody[:1000]))
	}

	// Write audio content to destination file
	logger.VPrintf("Writing audio to: %s\n", req.Dest)

	// Ensure the directory exists
	dir := req.Dest[:strings.LastIndex(req.Dest, "/")]
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.VPrintf("Error creating directory: %v\n", err)
		return false, err
	}

	// Write the audio data to file
	if err := os.WriteFile(req.Dest, respBody, 0644); err != nil {
		logger.VPrintf("Error writing file: %v\n", err)
		return false, err
	}

	logger.VPrintf("Successfully wrote %d bytes to %s\n", len(respBody), req.Dest)
	return true, nil
}

// isAudioFileValid checks if the audio file exists and is valid
func isAudioFileValid(file string) (bool, error) {
	if _, err := os.Stat(file); err != nil {
		logger.VPrintf("Warning: Audio file does not exist: %s\n", file)
		logger.VPrintf("Error: %v\n", err)
		return false, err
	}

	// Check if file has minimum size (not empty/corrupted)
	if fileInfo, err := os.Stat(file); err == nil {
		if fileInfo.Size() < 1000 {
			logger.VPrintf("Warning: Audio file appears to be corrupted or empty (size: %d bytes)\n", fileInfo.Size())
			return false, fmt.Errorf("file too small: %d bytes", fileInfo.Size())
		}
	}
	return true, nil
}
