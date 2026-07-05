package usecase

import (
	"context"
	"testing"

	"medical-webhook/internal/domain/line/entity"

	"github.com/google/uuid"
)

// stubAdminService returns a fixed admin from ValidateSession/GetProfile so we
// can assert the usecase maps every field (notably Role, which RBAC depends on)
// into the returned AdminDetail.
type stubAdminService struct{ admin *entity.Admin }

func (s *stubAdminService) Register(context.Context, string, string, string, string) (*entity.Admin, error) {
	return s.admin, nil
}
func (s *stubAdminService) Login(context.Context, string, string, string) (*entity.Admin, string, error) {
	return s.admin, "tok", nil
}
func (s *stubAdminService) Logout(context.Context, string) error { return nil }
func (s *stubAdminService) ValidateSession(context.Context, string) (*entity.Admin, error) {
	return s.admin, nil
}
func (s *stubAdminService) GetProfile(context.Context, uuid.UUID) (*entity.Admin, error) {
	return s.admin, nil
}
func (s *stubAdminService) UpdateProfile(context.Context, uuid.UUID, string, string) error {
	return nil
}
func (s *stubAdminService) ChangePassword(context.Context, uuid.UUID, string, string) error {
	return nil
}
func (s *stubAdminService) GetAllAdmins(context.Context, int, int) ([]*entity.Admin, error) {
	return []*entity.Admin{s.admin}, nil
}
func (s *stubAdminService) EnsureInitialSuperAdmin(context.Context, string, string, string, string) error {
	return nil
}

func newRoleAdmin() *entity.Admin {
	return &entity.Admin{ID: uuid.New(), Username: "boss", Email: "b@x.com", FullName: "Boss", Role: string(entity.RoleSuperAdmin)}
}

// ValidateToken feeds AuthMiddleware, which stores admin.Role for RequireRole.
// If Role is dropped here, every role-gated route rejects everyone.
func TestValidateToken_PopulatesRole(t *testing.T) {
	uc := NewAdminUsecase(&stubAdminService{admin: newRoleAdmin()})

	detail, err := uc.ValidateToken(context.Background(), "tok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if detail.Role != string(entity.RoleSuperAdmin) {
		t.Errorf("ValidateToken Role = %q, want super_admin", detail.Role)
	}
}

func TestGetProfile_PopulatesRole(t *testing.T) {
	admin := newRoleAdmin()
	uc := NewAdminUsecase(&stubAdminService{admin: admin})

	detail, err := uc.GetProfile(context.Background(), admin.ID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if detail.Role != string(entity.RoleSuperAdmin) {
		t.Errorf("GetProfile Role = %q, want super_admin", detail.Role)
	}
}
