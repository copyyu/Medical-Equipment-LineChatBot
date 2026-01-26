package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

// OCRResponse represents FINAL response from AI (Ai.py)
type OCRResponse struct {
	Status         string  `json:"status"`
	Code           string  `json:"code"`
	Source         string  `json:"source"`
	Confidence     float64 `json:"confidence"`
	ProcessingTime string  `json:"processing_time"`
	Error          string  `json:"error,omitempty"`
}

// OCRClient handles communication with OCR API
type OCRClient struct {
	baseURL string
	client  *http.Client
}

// NewOCRClient creates a new OCR client
func NewOCRClient(baseURL string) *OCRClient {
	return &OCRClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// ProcessImage sends an image to OCR API and returns AI final result
func (c *OCRClient) ProcessImage(imageBytes []byte, filename string) (*OCRResponse, error) {
	log.Printf("📤 Sending image to OCR API: %s", c.baseURL)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := part.Write(imageBytes); err != nil {
		return nil, fmt.Errorf("failed to write image data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/ocr", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var aiResp OCRResponse
	if err := json.Unmarshal(respBody, &aiResp); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	if resp.StatusCode != http.StatusOK || aiResp.Status != "success" {
		errMsg := aiResp.Error
		if errMsg == "" {
			errMsg = string(respBody)
		}
		return nil, fmt.Errorf("AI API error %d: %s", resp.StatusCode, errMsg)
	}

	log.Printf(
		"🧠 AI Result: code=%s source=%s confidence=%.2f time=%s",
		aiResp.Code,
		aiResp.Source,
		aiResp.Confidence,
		aiResp.ProcessingTime,
	)

	return &aiResp, nil
}

// GetDetectedCode returns detected code from AI final response
func (c *OCRClient) GetDetectedCode(resp *OCRResponse) string {
	if resp == nil {
		log.Println("⚠️ AI response is nil")
		return ""
	}

	if resp.Code == "" {
		log.Println("⚠️ No code detected")
		return ""
	}

	log.Printf(
		"📝 Detected code: %s (confidence: %.2f, source: %s)",
		resp.Code,
		resp.Confidence,
		resp.Source,
	)

	return resp.Code
}
