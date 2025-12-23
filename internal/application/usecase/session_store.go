package usecase

import (
	"log"
	"sync"
	"time"
)

// OCRSession stores pending OCR result for a user
type OCRSession struct {
	SerialNumber string
	ImageURL     string
	Timestamp    time.Time
}

// SessionStore manages user sessions in memory
type SessionStore struct {
	sessions map[string]*OCRSession
	mu       sync.RWMutex
}

// NewSessionStore creates a new session store
func NewSessionStore() *SessionStore {
	store := &SessionStore{
		sessions: make(map[string]*OCRSession),
	}
	// Start cleanup goroutine
	go store.cleanupLoop()
	return store
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

// cleanupLoop removes expired sessions (older than 10 minutes)
func (s *SessionStore) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanup()
	}
}

func (s *SessionStore) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	expiry := time.Now().Add(-10 * time.Minute)
	for userID, session := range s.sessions {
		if session.Timestamp.Before(expiry) {
			delete(s.sessions, userID)
			log.Printf("🧹 Cleaned expired session for user: %s", userID)
		}
	}
}
