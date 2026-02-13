package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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

// SearchResult represents a matched code with similarity score
type SearchResult struct {
	Code       string
	Similarity float64
	IsExact    bool
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
			Timeout: 90 * time.Second,
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

// ========== FUZZY MATCHING FUNCTIONS ==========

// NormalizeCode แยก prefix (ตัวอักษร) และ number (ตัวเลข)
// เช่น "A0021" → ("A", 21)
func NormalizeCode(code string) (string, int, error) {
	re := regexp.MustCompile(`^([A-Za-z]*)(\d+)$`)
	matches := re.FindStringSubmatch(code)

	if len(matches) != 3 {
		return "", 0, fmt.Errorf("invalid code format: %s", code)
	}

	prefix := strings.ToUpper(matches[1])
	number, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", 0, fmt.Errorf("invalid number in code: %w", err)
	}

	return prefix, number, nil
}

// ExactMatch เปรียบเทียบโค้ดแบบ normalize (ไม่สนใจ leading zeros)
func ExactMatch(code1, code2 string) bool {
	prefix1, num1, err1 := NormalizeCode(code1)
	prefix2, num2, err2 := NormalizeCode(code2)

	if err1 != nil || err2 != nil {
		// fallback: เปรียบเทียบแบบ case-insensitive
		return strings.EqualFold(code1, code2)
	}

	return prefix1 == prefix2 && num1 == num2
}

// Levenshtein คำนวณ edit distance
func Levenshtein(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,
				matrix[i][j-1]+1,
				matrix[i-1][j-1]+cost,
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// Similarity คำนวณ % ความคล้าย (0-100)
func Similarity(s1, s2 string) float64 {
	distance := Levenshtein(s1, s2)
	maxLen := len(s1)
	if len(s2) > maxLen {
		maxLen = len(s2)
	}

	if maxLen == 0 {
		return 100.0
	}

	return (1.0 - float64(distance)/float64(maxLen)) * 100.0
}

// SearchInDatabase ค้นหา code ที่ใกล้เคียงใน database
// threshold: ค่าต่ำสุดของ similarity ที่ยอมรับ (0-100)
func (c *OCRClient) SearchInDatabase(ocrCode string, dbCodes []string, threshold float64) []SearchResult {
	results := []SearchResult{}

	for _, dbCode := range dbCodes {
		// ลองเช็คแบบ exact match ก่อน
		if ExactMatch(ocrCode, dbCode) {
			results = append(results, SearchResult{
				Code:       dbCode,
				Similarity: 100.0,
				IsExact:    true,
			})
			continue
		}

		// คำนวณ similarity
		sim := Similarity(ocrCode, dbCode)
		if sim >= threshold {
			results = append(results, SearchResult{
				Code:       dbCode,
				Similarity: sim,
				IsExact:    false,
			})
		}
	}

	// เรียงจากมากไปน้อย
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Similarity > results[i].Similarity {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	return results
}

// FindBestMatch หา code ที่ตรงที่สุดจาก database
func (c *OCRClient) FindBestMatch(ocrCode string, dbCodes []string, minThreshold float64) *SearchResult {
	results := c.SearchInDatabase(ocrCode, dbCodes, minThreshold)

	if len(results) == 0 {
		log.Printf("⚠️ No match found for OCR code: %s (threshold: %.2f)", ocrCode, minThreshold)
		return nil
	}

	bestMatch := results[0]
	log.Printf(
		"✅ Best match: %s (similarity: %.2f%%, exact: %v) for OCR: %s",
		bestMatch.Code,
		bestMatch.Similarity,
		bestMatch.IsExact,
		ocrCode,
	)

	return &bestMatch
}
