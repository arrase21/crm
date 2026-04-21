package service

import (
	"context"
	"errors"
	"testing"

	"github.com/arrase21/crm/internal/domain"
)

// ============================================
// MOCK IMPLEMENTATION - Simula el Repository
// ============================================

type MockPositionRepo struct {
	positions    map[uint]*domain.Position
	nextID       uint
	nameIndex    map[string]*domain.Position
	CreateErr    error
	GetByIDErr   error
	GetByNameErr error
	UpdateErr    error
	DeleteErr    error
	ListErr      error
}

func NewMockPositionRepo() *MockPositionRepo {
	return &MockPositionRepo{
		positions: make(map[uint]*domain.Position),
		nextID:    1,
		nameIndex: make(map[string]*domain.Position),
	}
}

func (m *MockPositionRepo) Create(ctx context.Context, pstn *domain.Position) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	if pstn == nil {
		return errors.New("position cannot be nil")
	}
	if _, exists := m.nameIndex[pstn.Name]; exists {
		return domain.ErrPositionNameExists
	}
	pstn.ID = m.nextID
	m.nextID++
	m.positions[pstn.ID] = pstn
	m.nameIndex[pstn.Name] = pstn
	return nil
}

func (m *MockPositionRepo) GetByID(ctx context.Context, id uint) (*domain.Position, error) {
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	if id == 0 {
		return nil, errors.New("invalid position id")
	}
	pstn, exists := m.positions[id]
	if !exists {
		return nil, domain.ErrPositionNotFound
	}
	return pstn, nil
}

func (m *MockPositionRepo) GetByName(ctx context.Context, name string) (*domain.Position, error) {
	if m.GetByNameErr != nil {
		return nil, m.GetByNameErr
	}
	if name == "" {
		return nil, errors.New("position name cannot be empty")
	}
	pstn, exists := m.nameIndex[name]
	if !exists {
		return nil, domain.ErrPositionNotFound
	}
	return pstn, nil
}

func (m *MockPositionRepo) List(ctx context.Context, page, limit int) ([]domain.Position, int64, error) {
	if m.ListErr != nil {
		return nil, 0, m.ListErr
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var result []domain.Position
	var total int64
	for _, p := range m.positions {
		result = append(result, *p)
		total++
	}

	if offset >= len(result) {
		return []domain.Position{}, total, nil
	}
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end], total, nil
}

func (m *MockPositionRepo) Update(ctx context.Context, pstn *domain.Position) error {
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	if pstn == nil || pstn.ID == 0 {
		return errors.New("position cannot be nil or have zero id")
	}
	existing, exists := m.positions[pstn.ID]
	if !exists {
		return domain.ErrPositionNotFound
	}
	// Check name duplicate (excluding current position)
	if existing.Name != pstn.Name {
		if _, nameExists := m.nameIndex[pstn.Name]; nameExists {
			return domain.ErrPositionNameExists
		}
		delete(m.nameIndex, existing.Name)
		m.nameIndex[pstn.Name] = pstn
	}
	m.positions[pstn.ID] = pstn
	return nil
}

func (m *MockPositionRepo) Delete(ctx context.Context, id uint) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	if id == 0 {
		return errors.New("invalid position id")
	}
	pstn, exists := m.positions[id]
	if !exists {
		return domain.ErrPositionNotFound
	}
	delete(m.positions, id)
	delete(m.nameIndex, pstn.Name)
	return nil
}

// ============================================
// TESTS DEL SERVICE
// ============================================

func TestPositionService_Create_Success(t *testing.T) {
	mock := NewMockPositionRepo()
	svc := NewPositionService(mock)
	ctx := context.Background()

	position := &domain.Position{
		Name:        "Gerente",
		Description: "Gerente de área",
	}

	err := svc.Create(ctx, position)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if position.ID == 0 {
		t.Error("expected position ID to be set")
	}
}

func TestPositionService_Create_DuplicateName(t *testing.T) {
	mock := NewMockPositionRepo()
	svc := NewPositionService(mock)
	ctx := context.Background()

	position1 := &domain.Position{
		Name:        "Gerente",
		Description: "Gerente de área",
	}
	position2 := &domain.Position{
		Name:        "Gerente",
		Description: "Otro gerente",
	}

	_ = svc.Create(ctx, position1)
	err := svc.Create(ctx, position2)

	if !errors.Is(err, domain.ErrPositionNameExists) {
		t.Fatalf("expected ErrPositionNameExists, got %v", err)
	}
}

func TestPositionService_Create_ValidationError(t *testing.T) {
	mock := NewMockPositionRepo()
	svc := NewPositionService(mock)
	ctx := context.Background()

	position := &domain.Position{
		Name: "", // Required field
	}

	err := svc.Create(ctx, position)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestPositionService_GetByID_Success(t *testing.T) {
	mock := NewMockPositionRepo()
	svc := NewPositionService(mock)
	ctx := context.Background()

	created := &domain.Position{
		Name:        "Gerente",
		Description: "Gerente de área",
	}
	_ = svc.Create(ctx, created)

	result, err := svc.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Name != "Gerente" {
		t.Errorf("expected name 'Gerente', got %s", result.Name)
	}
}

func TestPositionService_GetByID_NotFound(t *testing.T) {
	mock := NewMockPositionRepo()
	svc := NewPositionService(mock)
	ctx := context.Background()

	_, err := svc.GetByID(ctx, 999)
	if !errors.Is(err, domain.ErrPositionNotFound) {
		t.Fatalf("expected ErrPositionNotFound, got %v", err)
	}
}

func TestPositionService_GetByID_InvalidID(t *testing.T) {
	mock := NewMockPositionRepo()
	svc := NewPositionService(mock)
	ctx := context.Background()

	_, err := svc.GetByID(ctx, 0)
	if err == nil {
		t.Fatal("expected error for invalid id")
	}
}

func TestPositionService_GetByName_Success(t *testing.T) {
	mock := NewMockPositionRepo()
	svc := NewPositionService(mock)
	ctx := context.Background()

	position := &domain.Position{
		Name:        "Gerente",
		Description: "Gerente de área",
	}
	_ = svc.Create(ctx, position)

	result, err := svc.GetByName(ctx, "Gerente")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Name != "Gerente" {
		t.Errorf("expected name 'Gerente', got %s", result.Name)
	}
}

func TestPositionService_GetByName_NotFound(t *testing.T) {
	mock := NewMockPositionRepo()
	svc := NewPositionService(mock)
	ctx := context.Background()

	_, err := svc.GetByName(ctx, "NoExiste")
	if !errors.Is(err, domain.ErrPositionNotFound) {
		t.Fatalf("expected ErrPositionNotFound, got %v", err)
	}
}

func TestPositionService_List_Success(t *testing.T) {
	mock := NewMockPositionRepo()
	svc := NewPositionService(mock)
	ctx := context.Background()

	// Create multiple positions
	for i := 1; i <= 5; i++ {
		_ = svc.Create(ctx, &domain.Position{
			Name:        "Position " + string(rune('0'+i)),
			Description: "Description " + string(rune('0'+i)),
		})
	}

	positions, total, err := svc.List(ctx, 1, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(positions) != 5 {
		t.Errorf("expected 5 positions, got %d", len(positions))
	}
}

func TestPositionService_Update_Success(t *testing.T) {
	mock := NewMockPositionRepo()
	svc := NewPositionService(mock)
	ctx := context.Background()

	position := &domain.Position{
		Name:        "Gerente",
		Description: "Gerente de área",
	}
	_ = svc.Create(ctx, position)

	position.Description = "Gerente actualizado"
	err := svc.Update(ctx, position)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPositionService_Update_DuplicateName(t *testing.T) {
	mock := NewMockPositionRepo()
	svc := NewPositionService(mock)
	ctx := context.Background()

	position1 := &domain.Position{
		Name:        "Gerente",
		Description: "Gerente de área",
	}
	position2 := &domain.Position{
		Name:        "Asistente",
		Description: "Asistente",
	}
	_ = svc.Create(ctx, position1)
	_ = svc.Create(ctx, position2)

	position2.Name = "Gerente"
	err := svc.Update(ctx, position2)
	if !errors.Is(err, domain.ErrPositionNameExists) {
		t.Fatalf("expected ErrPositionNameExists, got %v", err)
	}
}

func TestPositionService_Delete_Success(t *testing.T) {
	mock := NewMockPositionRepo()
	svc := NewPositionService(mock)
	ctx := context.Background()

	position := &domain.Position{
		Name:        "Gerente",
		Description: "Gerente de área",
	}
	_ = svc.Create(ctx, position)

	err := svc.Delete(ctx, position.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify deleted
	_, err = svc.GetByID(ctx, position.ID)
	if !errors.Is(err, domain.ErrPositionNotFound) {
		t.Fatalf("expected ErrPositionNotFound after delete, got %v", err)
	}
}

func TestPositionService_Delete_NotFound(t *testing.T) {
	mock := NewMockPositionRepo()
	svc := NewPositionService(mock)
	ctx := context.Background()

	err := svc.Delete(ctx, 999)
	if !errors.Is(err, domain.ErrPositionNotFound) {
		t.Fatalf("expected ErrPositionNotFound, got %v", err)
	}
}
