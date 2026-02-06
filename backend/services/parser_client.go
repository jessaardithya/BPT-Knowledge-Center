package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ParserResponse struct {
	Filename     string          `json:"filename"`
	ContentType  string          `json:"content_type"`
	ElementCount int             `json:"element_count"`
	Data         []ParsedElement `json:"data"`
}

type ParsedElement struct {
	ElementID string                 `json:"element_id"`
	Type      string                 `json:"type"`
	Text      string                 `json:"text"`
	Metadata  map[string]interface{} `json:"metadata"`
	Vector    []float32              `json:"vector"`
}

func SendToParser(filePath string) (*ParserResponse, error) {
	pythonURL := os.Getenv("PARSER_URL")
	if pythonURL == "" {
		return nil, fmt.Errorf("PARSER_URL not set in .env")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("could not create form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("could not copy file content: %v", err)
	}
	writer.Close()

	req, _ := http.NewRequest("POST", pythonURL, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Python service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("python service error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result ParserResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %v", err)
	}

	return &result, nil
}
