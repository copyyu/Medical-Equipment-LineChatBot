package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/domain/line/repository"
	apperrors "medical-webhook/internal/utils/errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Admin sentinel errors alias the canonical set in utils/errors so the HTTP
// error mapper (MapErrorToResponse) matches them with errors.Is and returns the
// right status (401/409/…) instead of a generic 500.
var (
	ErrInvalidCredentials = apperrors.ErrInvalidCredentials
	ErrAdminNotFound      = apperrors.ErrAdminNotFound
	ErrAdminInactive      = apperrors.ErrAdminInactive
	ErrInvalidToken       = apperrors.ErrInvalidToken
	ErrUsernameExists     = apperrors.ErrUsernameExists
	ErrEmailExists        = apperrors.ErrEmailExists
)

type AdminService interface {
	Register(ctx context.Context, username, email, password, fullName string) (*entity.Admin, error)
	Login(ctx context.Context, username, password, ipAddress string) (*entity.Admin, string, error)
	Logout(ctx context.Context, token string) error
	ValidateSession(ctx context.Context, token string) (*entity.Admin, error)
	GetProfile(ctx context.Context, adminID uuid.UUID) (*entity.Admin, error)
	UpdateProfile(ctx context.Context, adminID uuid.UUID, fullName, email string) error
	ChangePassword(ctx context.Context, adminID uuid.UUID, oldPassword, newPassword string) error
	GetAllAdmins(ctx context.Context, limit, offset int) ([]*entity.Admin, error)
}

type adminService struct {
	adminRepo   repository.AdminRepository
	sessionRepo repository.AdminSessionRepository
}

func NewAdminService(adminRepo repository.AdminRepository, sessionRepo repository.AdminSessionRepository) AdminService {
	return &adminService{
		adminRepo:   adminRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *adminService) Register(ctx context.Context, username, email, password, fullName string) (*entity.Admin, error) {
	return s.createAdmin(ctx, username, email, password, fullName, entity.RoleAdmin)
}

// createAdmin validates uniqueness, hashes the password, and persists a new
// admin with the given role. Shared by Register (regular admins) and the
// initial super-admin bootstrap.
func (s *adminService) createAdmin(ctx context.Context, username, email, password, fullName string, role entity.AdminRole) (*entity.Admin, error) {
	// Check if username exists
	existingAdmin, err := s.adminRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existingAdmin != nil {
		return nil, ErrUsernameExists
	}

	// Check if email exists
	existingAdmin, err = s.adminRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingAdmin != nil {
		return nil, ErrEmailExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	admin := &entity.Admin{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		FullName:     fullName,
		Role:         string(role),
	}

	if err := s.adminRepo.Create(ctx, admin); err != nil {
		return nil, err
	}

	return admin, nil
}

func (s *adminService) Login(ctx context.Context, username, password, ipAddress string) (*entity.Admin, string, error) {
	// Get admin by username
	admin, err := s.adminRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, "", err
	}
	if admin == nil {
		return nil, "", ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Generate session token
	token, err := generateToken()
	if err != nil {
		return nil, "", err
	}

	// Create session
	session := &entity.AdminSession{
		AdminID:   admin.ID,
		Token:     token,
		IPAddress: ipAddress,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, "", err
	}

	// Update last login (best-effort): the session is already created and the
	// token is valid, so a failure here must not fail an otherwise-successful
	// login — just log it.
	if err := s.adminRepo.UpdateLastLogin(ctx, admin.ID); err != nil {
		log.Printf("⚠️ Failed to update last-login for admin %s: %v", admin.Username, err)
	}

	return admin, token, nil
}

func (s *adminService) Logout(ctx context.Context, token string) error {
	return s.sessionRepo.DeleteByToken(ctx, token)
}

func (s *adminService) ValidateSession(ctx context.Context, token string) (*entity.Admin, error) {
	session, err := s.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrInvalidToken
	}

	admin, err := s.adminRepo.GetByID(ctx, session.AdminID)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, ErrAdminNotFound
	}

	return admin, nil
}

func (s *adminService) GetProfile(ctx context.Context, adminID uuid.UUID) (*entity.Admin, error) {
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, ErrAdminNotFound
	}

	return admin, nil
}

func (s *adminService) UpdateProfile(ctx context.Context, adminID uuid.UUID, fullName, email string) error {
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return err
	}
	if admin == nil {
		return ErrAdminNotFound
	}

	admin.FullName = fullName
	admin.Email = email

	return s.adminRepo.Update(ctx, admin)
}

func (s *adminService) ChangePassword(ctx context.Context, adminID uuid.UUID, oldPassword, newPassword string) error {
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return err
	}
	if admin == nil {
		return ErrAdminNotFound
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin.PasswordHash = string(hashedPassword)

	return s.adminRepo.Update(ctx, admin)
}

func (s *adminService) GetAllAdmins(ctx context.Context, limit, offset int) ([]*entity.Admin, error) {
	return s.adminRepo.List(ctx, limit, offset)
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
