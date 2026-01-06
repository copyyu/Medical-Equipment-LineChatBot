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

// OCRResult represents a single OCR result within a detection
type OCRResult struct {
	Raw        string  `json:"raw"`        // ข้อความดิบที่ OCR อ่านได้
	Normalized string  `json:"normalized"` // ข้อความหลัง normalize
	Confidence float64 `json:"confidence"` // ค่าความมั่นใจของ OCR (0-1)
	IsSSH      bool    `json:"is_ssh"`     // ตรงกับ pattern SSH หรือไม่
}

// OCRDetection represents a detection with bounding box and OCR results
type OCRDetection struct {
	BBox           []float64   `json:"bbox"`            // พิกัด [x1, y1, x2, y2]
	YoloConfidence float64     `json:"yolo_confidence"` // ค่าความมั่นใจจาก YOLO (0-1)
	OCRResults     []OCRResult `json:"ocr_results"`     // ผลลัพธ์ OCR
}

// OCRResponse represents the response from OCR API
type OCRResponse struct {
	Status     string         `json:"status"`     // สถานะ ("success")
	Count      int            `json:"count"`      // จำนวน detections
	Detections []OCRDetection `json:"detections"` // รายการ detection
	Error      string         `json:"error"`      // error message (if any)
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

	// Parse response
	var ocrResp OCRResponse
	if err := json.Unmarshal(respBody, &ocrResp); err != nil {
		return nil, fmt.Errorf("failed to parse OCR response: %w", err)
	}

	// Check for API error
	if resp.StatusCode != http.StatusOK || ocrResp.Error != "" {
		errMsg := ocrResp.Error
		if errMsg == "" {
			errMsg = string(respBody)
		}
		return nil, fmt.Errorf("OCR API error %d: %s", resp.StatusCode, errMsg)
	}

	log.Printf("✅ OCR API response: status=%s, count=%d", ocrResp.Status, ocrResp.Count)
	return &ocrResp, nil
}

// GetDetectedCode returns the best detected code from OCR response
// Prioritizes: 1) is_ssh=true with highest confidence, 2) highest confidence normalized text
func (c *OCRClient) GetDetectedCode(resp *OCRResponse) string {
	if resp == nil || len(resp.Detections) == 0 {
		return ""
	}

	var bestCode string
	var bestConfidence float64
	var foundSSH bool

	// Search through all detections and OCR results
	for _, detection := range resp.Detections {
		for _, result := range detection.OCRResults {
			// Prioritize is_ssh=true results
			if result.IsSSH {
				if !foundSSH || result.Confidence > bestConfidence {
					foundSSH = true
					bestConfidence = result.Confidence
					bestCode = result.Normalized
				}
			} else if !foundSSH && result.Confidence > bestConfidence {
				// If no SSH found yet, use highest confidence
				bestConfidence = result.Confidence
				bestCode = result.Normalized
			}
		}
	}

	if bestCode != "" {
		log.Printf("📝 Detected code: %s (confidence: %.2f, is_ssh: %v)", bestCode, bestConfidence, foundSSH)
	} else {
		log.Println("⚠️ No code detected from OCR")
	}

	return bestCode
}
