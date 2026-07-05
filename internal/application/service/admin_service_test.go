package service

import (
	"context"
	"errors"
	"testing"

	"medical-webhook/internal/domain/line/entity"

	"github.com/google/uuid"
)

// fakeAdminRepo is a minimal in-memory AdminRepository for exercising the
// super-admin bootstrap logic.
type fakeAdminRepo struct {
	admins []*entity.Admin
}

func (f *fakeAdminRepo) Create(_ context.Context, a *entity.Admin) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	f.admins = append(f.admins, a)
	return nil
}
func (f *fakeAdminRepo) GetByID(_ context.Context, id uuid.UUID) (*entity.Admin, error) {
	for _, a := range f.admins {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, nil
}
func (f *fakeAdminRepo) GetByUsername(_ context.Context, username string) (*entity.Admin, error) {
	for _, a := range f.admins {
		if a.Username == username {
			return a, nil
		}
	}
	return nil, nil
}
func (f *fakeAdminRepo) GetByEmail(_ context.Context, email string) (*entity.Admin, error) {
	for _, a := range f.admins {
		if a.Email == email {
			return a, nil
		}
	}
	return nil, nil
}
func (f *fakeAdminRepo) Update(_ context.Context, a *entity.Admin) error      { return nil }
func (f *fakeAdminRepo) UpdateLastLogin(_ context.Context, _ uuid.UUID) error { return nil }
func (f *fakeAdminRepo) Delete(_ context.Context, _ uuid.UUID) error          { return nil }
func (f *fakeAdminRepo) List(_ context.Context, limit, offset int) ([]*entity.Admin, error) {
	return f.admins, nil
}

func newAdminTestService(repo *fakeAdminRepo) AdminService {
	return NewAdminService(repo, nil)
}

func TestEnsureInitialSuperAdmin_CreatesWhenNone(t *testing.T) {
	repo := &fakeAdminRepo{}
	svc := newAdminTestService(repo)

	if err := svc.EnsureInitialSuperAdmin(context.Background(), "root", "root@x.com", "pw", "Root"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.admins) != 1 {
		t.Fatalf("expected 1 admin created, got %d", len(repo.admins))
	}
	if repo.admins[0].Role != string(entity.RoleSuperAdmin) {
		t.Errorf("created admin role = %q, want super_admin", repo.admins[0].Role)
	}
}

func TestEnsureInitialSuperAdmin_PromotesExisting(t *testing.T) {
	repo := &fakeAdminRepo{admins: []*entity.Admin{
		{ID: uuid.New(), Username: "root", Email: "root@x.com", Role: string(entity.RoleAdmin)},
	}}
	svc := newAdminTestService(repo)

	if err := svc.EnsureInitialSuperAdmin(context.Background(), "root", "root@x.com", "pw", "Root"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.admins) != 1 {
		t.Fatalf("expected no new admin, got %d", len(repo.admins))
	}
	if repo.admins[0].Role != string(entity.RoleSuperAdmin) {
		t.Errorf("existing admin role = %q, want promoted to super_admin", repo.admins[0].Role)
	}
}

func TestEnsureInitialSuperAdmin_NoopWhenSuperAdminExists(t *testing.T) {
	repo := &fakeAdminRepo{admins: []*entity.Admin{
		{ID: uuid.New(), Username: "boss", Role: string(entity.RoleSuperAdmin)},
	}}
	svc := newAdminTestService(repo)

	// Even with no config, an existing super admin means nothing to do.
	if err := svc.EnsureInitialSuperAdmin(context.Background(), "", "", "", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.admins) != 1 {
		t.Errorf("expected admins unchanged, got %d", len(repo.admins))
	}
}

func TestEnsureInitialSuperAdmin_ErrWhenUnconfigured(t *testing.T) {
	repo := &fakeAdminRepo{admins: []*entity.Admin{
		{ID: uuid.New(), Username: "reg", Role: string(entity.RoleAdmin)},
	}}
	svc := newAdminTestService(repo)

	err := svc.EnsureInitialSuperAdmin(context.Background(), "", "", "", "")
	if !errors.Is(err, ErrNoInitialAdminConfig) {
		t.Errorf("want ErrNoInitialAdminConfig, got %v", err)
	}
}
