package usecase

import (
	"log"
	"sync"
	"time"
)

// Session timeout constants
const (
	SessionExpiryDuration = 10 * time.Minute
	CleanupInterval       = 5 * time.Minute
)

// SessionMode represents what action user is doing
type SessionMode string

const (
	ModeNone           SessionMode = ""
	ModeReportProblem  SessionMode = "report_problem"   // แจ้งปัญหา/เช็กสถานะ
	ModeTrackStatus    SessionMode = "track_status"     // ติดตามสถานะ
	ModeInputIssueDesc SessionMode = "input_issue_desc" // รอพิมพ์รายละเอียดปัญหา
)

// OCRSession stores pending OCR result for a user
type OCRSession struct {
	Mode         SessionMode // โหมดที่ผู้ใช้เลือกจากเมนู
	SerialNumber string
	ImageURL     string
	Timestamp    time.Time
}

// SessionStore manages user sessions in memory
type SessionStore struct {
	sessions map[string]*OCRSession
	mu       sync.RWMutex
	done     chan struct{} // channel สำหรับ signal shutdown
}

// NewSessionStore creates a new session store
func NewSessionStore() *SessionStore {
	store := &SessionStore{
		sessions: make(map[string]*OCRSession),
		done:     make(chan struct{}),
	}
	// Start cleanup goroutine
	go store.cleanupLoop()
	return store
}

// Close gracefully shuts down the session store cleanup goroutine
func (s *SessionStore) Close() {
	close(s.done)
	log.Println("🛑 Session store closed")
}

// Set stores OCR session for a user
func (s *SessionStore) Set(userID string, session *OCRSession) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session.Timestamp = time.Now()
	s.sessions[userID] = session
	log.Printf("📦 Session stored for user: %s, serial: %s", userID, session.SerialNumber)
}

// Get retrieves OCR session for a user
func (s *SessionStore) Get(userID string) *OCRSession {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, exists := s.sessions[userID]
	if !exists {
		return nil
	}
	return session
}

// Delete removes OCR session for a user
func (s *SessionStore) Delete(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, userID)
	log.Printf("🗑️ Session deleted for user: %s", userID)
}

// cleanupLoop removes expired sessions periodically
func (s *SessionStore) cleanupLoop() {
	ticker := time.NewTicker(CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanup()
		case <-s.done:
			log.Println("🛑 Session cleanup loop stopped")
			return // graceful shutdown
		}
	}
}

func (s *SessionStore) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	expiry := time.Now().Add(-SessionExpiryDuration)
	for userID, session := range s.sessions {
		if session.Timestamp.Before(expiry) {
			delete(s.sessions, userID)
			log.Printf("🧹 Cleaned expired session for user: %s", userID)
		}
	}
}
