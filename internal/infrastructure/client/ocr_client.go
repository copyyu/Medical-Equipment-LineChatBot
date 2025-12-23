package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

// OCRText represents a single OCR detection result
type OCRText struct {
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence"`
}

// OCRDetection represents a detection with bounding box and OCR texts
type OCRDetection struct {
	BBox           []float64 `json:"bbox"`
	YoloConfidence float64   `json:"yolo_confidence"`
	CropPath       string    `json:"crop_path"`
	OCRTexts       []OCRText `json:"ocr_texts"`
}

// OCRResponse represents the response from OCR API
type OCRResponse struct {
	Status     string         `json:"status"`
	Detections []OCRDetection `json:"detections"`
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
		client:  &http.Client{},
	}
}

// ProcessImage sends an image to OCR API and returns detected text
func (c *OCRClient) ProcessImage(imageBytes []byte, filename string) (*OCRResponse, error) {
	log.Printf("📤 Sending image to OCR API: %s", c.baseURL)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file field
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

	// Create request
	req, err := http.NewRequest("POST", c.baseURL+"/ocr", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OCR API error %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var ocrResp OCRResponse
	if err := json.Unmarshal(respBody, &ocrResp); err != nil {
		return nil, fmt.Errorf("failed to parse OCR response: %w", err)
	}

	log.Printf("✅ OCR API response: status=%s, detections=%d", ocrResp.Status, len(ocrResp.Detections))
	return &ocrResp, nil
}

// GetBestText returns the text with highest confidence from OCR response
func (c *OCRClient) GetBestText(resp *OCRResponse) string {
	if resp == nil || len(resp.Detections) == 0 {
		return ""
	}

	var bestText string
	var bestConfidence float64

	for _, detection := range resp.Detections {
		for _, text := range detection.OCRTexts {
			if text.Confidence > bestConfidence {
				bestConfidence = text.Confidence
				bestText = text.Text
			}
		}
	}

	log.Printf("📝 Best OCR text: %s (confidence: %.2f)", bestText, bestConfidence)
	return bestText
}
