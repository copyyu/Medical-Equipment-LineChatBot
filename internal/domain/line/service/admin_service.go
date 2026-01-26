package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/domain/line/repository"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAdminNotFound      = errors.New("admin not found")
	ErrAdminInactive      = errors.New("admin is inactive")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrUsernameExists     = errors.New("username already exists")
	ErrEmailExists        = errors.New("email already exists")
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

	// Update last login
	if err := s.adminRepo.UpdateLastLogin(ctx, admin.ID); err != nil {
		return nil, "", err
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
