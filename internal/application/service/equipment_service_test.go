package service

import (
	"context"
	"errors"
	"testing"

	"medical-webhook/internal/domain/line/entity"
	mockRepo "medical-webhook/internal/mocks/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ─── Helper ───────────────────────────────────────────────

// newTestService creates an equipmentService with all mock repos injected.
func newTestService(t *testing.T) (
	*equipmentService,
	*mockRepo.MockEquipmentRepository,
	*mockRepo.MockBrandRepository,
	*mockRepo.MockEquipmentCategoryRepository,
	*mockRepo.MockDepartmentRepository,
	*mockRepo.MockEquipmentModelRepository,
) {
	equipRepo := mockRepo.NewMockEquipmentRepository(t)
	brandRepo := mockRepo.NewMockBrandRepository(t)
	catRepo := mockRepo.NewMockEquipmentCategoryRepository(t)
	deptRepo := mockRepo.NewMockDepartmentRepository(t)
	modelRepo := mockRepo.NewMockEquipmentModelRepository(t)

	svc := &equipmentService{
		equipmentRepo:  equipRepo,
		brandRepo:      brandRepo,
		categoryRepo:   catRepo,
		departmentRepo: deptRepo,
		modelRepo:      modelRepo,
	}

	return svc, equipRepo, brandRepo, catRepo, deptRepo, modelRepo
}

// ─── FindEquipmentByIDCode ────────────────────────────────

func TestFindEquipmentByIDCode_Success(t *testing.T) {
	svc, equipRepo, _, _, _, _ := newTestService(t)
	ctx := context.Background()

	expected := &entity.Equipment{IDCode: "SSH01234"}
	equipRepo.EXPECT().FindByIDCode("SSH01234").Return(expected, nil)

	result, err := svc.FindEquipmentByIDCode(ctx, "SSH01234")

	require.NoError(t, err)
	assert.Equal(t, "SSH01234", result.IDCode)
	equipRepo.AssertExpectations(t)
}

func TestFindEquipmentByIDCode_NotFound(t *testing.T) {
	svc, equipRepo, _, _, _, _ := newTestService(t)
	ctx := context.Background()

	equipRepo.EXPECT().FindByIDCode("NOTEXIST").Return(nil, nil)

	result, err := svc.FindEquipmentByIDCode(ctx, "NOTEXIST")

	require.NoError(t, err)
	assert.Nil(t, result)
	equipRepo.AssertExpectations(t)
}

func TestFindEquipmentByIDCode_RepoError(t *testing.T) {
	svc, equipRepo, _, _, _, _ := newTestService(t)
	ctx := context.Background()

	equipRepo.EXPECT().FindByIDCode("SSH01234").Return(nil, errors.New("db connection timeout"))

	result, err := svc.FindEquipmentByIDCode(ctx, "SSH01234")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "db connection timeout")
	equipRepo.AssertExpectations(t)
}

// ─── CountEquipments ──────────────────────────────────────

func TestCountEquipments_Success(t *testing.T) {
	svc, equipRepo, _, _, _, _ := newTestService(t)
	ctx := context.Background()

	equipRepo.EXPECT().Count(mock.Anything).Return(int64(42), nil)

	count, err := svc.CountEquipments(ctx)

	require.NoError(t, err)
	assert.Equal(t, int64(42), count)
	equipRepo.AssertExpectations(t)
}

func TestCountEquipments_Error(t *testing.T) {
	svc, equipRepo, _, _, _, _ := newTestService(t)
	ctx := context.Background()

	equipRepo.EXPECT().Count(mock.Anything).Return(int64(0), errors.New("query failed"))

	count, err := svc.CountEquipments(ctx)

	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	equipRepo.AssertExpectations(t)
}

// ─── CreateEquipment ──────────────────────────────────────

func TestCreateEquipment_Success(t *testing.T) {
	svc, equipRepo, _, _, _, _ := newTestService(t)
	ctx := context.Background()

	equipment := &entity.Equipment{
		IDCode:       "MED001",
		ModelID:      1,
		DepartmentID: 1,
	}
	equipRepo.EXPECT().Create(mock.Anything, equipment).Return(nil)

	err := svc.CreateEquipment(ctx, equipment)

	require.NoError(t, err)
	equipRepo.AssertExpectations(t)
}

func TestCreateEquipment_ValidationError_EmptyIDCode(t *testing.T) {
	svc, _, _, _, _, _ := newTestService(t)
	ctx := context.Background()

	equipment := &entity.Equipment{
		IDCode:       "", // invalid: empty
		ModelID:      1,
		DepartmentID: 1,
	}

	err := svc.CreateEquipment(ctx, equipment)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "equipment ID code is required")
}

func TestCreateEquipment_ValidationError_ZeroModelID(t *testing.T) {
	svc, _, _, _, _, _ := newTestService(t)
	ctx := context.Background()

	equipment := &entity.Equipment{
		IDCode:       "MED001",
		ModelID:      0, // invalid: zero
		DepartmentID: 1,
	}

	err := svc.CreateEquipment(ctx, equipment)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "equipment model ID is required")
}

func TestCreateEquipment_ValidationError_ZeroDepartmentID(t *testing.T) {
	svc, _, _, _, _, _ := newTestService(t)
	ctx := context.Background()

	equipment := &entity.Equipment{
		IDCode:       "MED001",
		ModelID:      1,
		DepartmentID: 0, // invalid: zero
	}

	err := svc.CreateEquipment(ctx, equipment)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "equipment department ID is required")
}

func TestCreateEquipment_RepoError(t *testing.T) {
	svc, equipRepo, _, _, _, _ := newTestService(t)
	ctx := context.Background()

	equipment := &entity.Equipment{
		IDCode:       "MED001",
		ModelID:      1,
		DepartmentID: 1,
	}
	equipRepo.EXPECT().Create(mock.Anything, equipment).Return(errors.New("duplicate key"))

	err := svc.CreateEquipment(ctx, equipment)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key")
	equipRepo.AssertExpectations(t)
}

// ─── DeleteEquipment ──────────────────────────────────────

func TestDeleteEquipment_Success(t *testing.T) {
	svc, equipRepo, _, _, _, _ := newTestService(t)
	ctx := context.Background()

	equipRepo.EXPECT().Delete(mock.Anything, uint(10)).Return(nil)

	err := svc.DeleteEquipment(ctx, 10)

	require.NoError(t, err)
	equipRepo.AssertExpectations(t)
}

func TestDeleteEquipment_ZeroID(t *testing.T) {
	svc, _, _, _, _, _ := newTestService(t)
	ctx := context.Background()

	err := svc.DeleteEquipment(ctx, 0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "equipment ID is required")
}

func TestDeleteEquipment_RepoError(t *testing.T) {
	svc, equipRepo, _, _, _, _ := newTestService(t)
	ctx := context.Background()

	equipRepo.EXPECT().Delete(mock.Anything, uint(99)).Return(errors.New("record not found"))

	err := svc.DeleteEquipment(ctx, 99)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
	equipRepo.AssertExpectations(t)
}

// ─── UpdateEquipment ──────────────────────────────────────

func TestUpdateEquipment_Success(t *testing.T) {
	svc, equipRepo, _, _, _, _ := newTestService(t)
	ctx := context.Background()

	equipment := &entity.Equipment{
		IDCode:       "MED001",
		ModelID:      1,
		DepartmentID: 1,
	}
	equipRepo.EXPECT().Update(mock.Anything, equipment).Return(nil)

	err := svc.UpdateEquipment(ctx, equipment)

	require.NoError(t, err)
	equipRepo.AssertExpectations(t)
}

func TestUpdateEquipment_ValidationError(t *testing.T) {
	svc, _, _, _, _, _ := newTestService(t)
	ctx := context.Background()

	equipment := &entity.Equipment{
		IDCode: "", // invalid
	}

	err := svc.UpdateEquipment(ctx, equipment)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "equipment ID code is required")
}

// ─── FindOrCreateBrand ────────────────────────────────────

func TestFindOrCreateBrand_ExistingBrand(t *testing.T) {
	svc, _, brandRepo, _, _, _ := newTestService(t)
	ctx := context.Background()

	existing := &entity.Brand{Name: "Philips"}
	existing.ID = 5
	brandRepo.EXPECT().FindByName(mock.Anything, "Philips").Return(existing, nil)

	result, err := svc.FindOrCreateBrand(ctx, "Philips")

	require.NoError(t, err)
	assert.Equal(t, uint(5), result.ID)
	assert.Equal(t, "Philips", result.Name)
	brandRepo.AssertExpectations(t)
}

func TestFindOrCreateBrand_NewBrand(t *testing.T) {
	svc, _, brandRepo, _, _, _ := newTestService(t)
	ctx := context.Background()

	brandRepo.EXPECT().FindByName(mock.Anything, "Siemens").Return(nil, nil)
	brandRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*entity.Brand")).Return(nil)

	result, err := svc.FindOrCreateBrand(ctx, "Siemens")

	require.NoError(t, err)
	assert.Equal(t, "Siemens", result.Name)
	brandRepo.AssertExpectations(t)
}

// ─── FindOrCreateDepartment ──────────────────────────────

func TestFindOrCreateDepartment_Existing(t *testing.T) {
	svc, _, _, _, deptRepo, _ := newTestService(t)
	ctx := context.Background()

	existing := &entity.Department{Name: "ICU"}
	existing.ID = 3
	deptRepo.EXPECT().FindByName(mock.Anything, "ICU").Return(existing, nil)

	result, err := svc.FindOrCreateDepartment(ctx, "ICU")

	require.NoError(t, err)
	assert.Equal(t, uint(3), result.ID)
	deptRepo.AssertExpectations(t)
}

func TestFindOrCreateDepartment_New(t *testing.T) {
	svc, _, _, _, deptRepo, _ := newTestService(t)
	ctx := context.Background()

	deptRepo.EXPECT().FindByName(mock.Anything, "ER").Return(nil, nil)
	deptRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*entity.Department")).Return(nil)

	result, err := svc.FindOrCreateDepartment(ctx, "ER")

	require.NoError(t, err)
	assert.Equal(t, "ER", result.Name)
	deptRepo.AssertExpectations(t)
}

// ─── GetAllCategories ─────────────────────────────────────

func TestGetAllCategories_Success(t *testing.T) {
	svc, _, _, catRepo, _, _ := newTestService(t)
	ctx := context.Background()

	expected := []entity.EquipmentCategory{
		{Name: "Ventilator"},
		{Name: "Monitor"},
	}
	catRepo.EXPECT().FindAll(mock.Anything).Return(expected, nil)

	result, err := svc.GetAllCategories(ctx)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Ventilator", result[0].Name)
	catRepo.AssertExpectations(t)
}

func TestGetAllCategories_Empty(t *testing.T) {
	svc, _, _, catRepo, _, _ := newTestService(t)
	ctx := context.Background()

	catRepo.EXPECT().FindAll(mock.Anything).Return([]entity.EquipmentCategory{}, nil)

	result, err := svc.GetAllCategories(ctx)

	require.NoError(t, err)
	assert.Empty(t, result)
	catRepo.AssertExpectations(t)
}
