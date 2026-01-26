package usecase

import (
	"context"
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/domain/line/service"

	"github.com/google/uuid"
)

type AdminUsecase interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AdminDetail, error)
	Login(ctx context.Context, req *dto.LoginRequest, ipAddress string) (*dto.LoginResponse, error)
	Logout(ctx context.Context, token string) error
	GetProfile(ctx context.Context, adminID string) (*dto.AdminDetail, error)
	UpdateProfile(ctx context.Context, adminID string, req *dto.UpdateProfileRequest) error
	ChangePassword(ctx context.Context, adminID string, req *dto.ChangePasswordRequest) error
	ValidateToken(ctx context.Context, token string) (*dto.AdminDetail, error)
}

type adminUsecase struct {
	adminService service.AdminService
}

func NewAdminUsecase(adminService service.AdminService) AdminUsecase {
	return &adminUsecase{
		adminService: adminService,
	}
}

func (u *adminUsecase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AdminDetail, error) {
	admin, err := u.adminService.Register(ctx, req.Username, req.Email, req.Password, req.FullName)
	if err != nil {
		return nil, err
	}

	return &dto.AdminDetail{
		ID:       admin.ID.String(),
		Username: admin.Username,
		Email:    admin.Email,
		FullName: admin.FullName,
		Role:     admin.Role,
	}, nil
}

func (u *adminUsecase) Login(ctx context.Context, req *dto.LoginRequest, ipAddress string) (*dto.LoginResponse, error) {
	admin, token, err := u.adminService.Login(ctx, req.Username, req.Password, ipAddress)
	if err != nil {
		return nil, err
	}

	var lastLogin *string
	if admin.LastLoginAt != nil {
		loginStr := admin.LastLoginAt.Format("2006-01-02 15:04:05")
		lastLogin = &loginStr
	}

	return &dto.LoginResponse{
		Token: token,
		Admin: dto.AdminDetail{
			ID:          admin.ID.String(),
			Username:    admin.Username,
			Email:       admin.Email,
			FullName:    admin.FullName,
			Role:        admin.Role,
			LastLoginAt: lastLogin,
		},
	}, nil
}

func (u *adminUsecase) Logout(ctx context.Context, token string) error {
	return u.adminService.Logout(ctx, token)
}

func (u *adminUsecase) GetProfile(ctx context.Context, adminID string) (*dto.AdminDetail, error) {
	id, err := uuid.Parse(adminID)
	if err != nil {
		return nil, err
	}

	admin, err := u.adminService.GetProfile(ctx, id)
	if err != nil {
		return nil, err
	}

	var lastLogin *string
	if admin.LastLoginAt != nil {
		loginStr := admin.LastLoginAt.Format("2006-01-02 15:04:05")
		lastLogin = &loginStr
	}

	return &dto.AdminDetail{
		ID:          admin.ID.String(),
		Username:    admin.Username,
		Email:       admin.Email,
		FullName:    admin.FullName,
		LastLoginAt: lastLogin,
	}, nil
}

func (u *adminUsecase) UpdateProfile(ctx context.Context, adminID string, req *dto.UpdateProfileRequest) error {
	id, err := uuid.Parse(adminID)
	if err != nil {
		return err
	}

	return u.adminService.UpdateProfile(ctx, id, req.FullName, req.Email)
}

func (u *adminUsecase) ChangePassword(ctx context.Context, adminID string, req *dto.ChangePasswordRequest) error {
	id, err := uuid.Parse(adminID)
	if err != nil {
		return err
	}

	return u.adminService.ChangePassword(ctx, id, req.OldPassword, req.NewPassword)
}

func (u *adminUsecase) ValidateToken(ctx context.Context, token string) (*dto.AdminDetail, error) {
	admin, err := u.adminService.ValidateSession(ctx, token)
	if err != nil {
		return nil, err
	}

	var lastLogin *string
	if admin.LastLoginAt != nil {
		loginStr := admin.LastLoginAt.Format("2006-01-02 15:04:05")
		lastLogin = &loginStr
	}

	return &dto.AdminDetail{
		ID:          admin.ID.String(),
		Username:    admin.Username,
		Email:       admin.Email,
		FullName:    admin.FullName,
		LastLoginAt: lastLogin,
	}, nil
}
