package example_user

import (
	"errors"
	"go-template/pkg/auth"
	"go-template/pkg/custom_errors"
	"go-template/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// ====================================================================================
// Mock Dependencies (นักแสดงแทนของเรา)
// ====================================================================================

// MockRepository คือ "นักแสดงแทน" ที่จะเล่นเป็น Repository ของเรา
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(d *Domain) error {
	args := m.Called(d)
	if d != nil {
		d.ID = 1
	}
	return args.Error(0)
}
func (m *MockRepository) GetByEmail(email string) (*Domain, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain), args.Error(1)
}
func (m *MockRepository) GetByID(id uint) (*Domain, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain), args.Error(1)
}
func (m *MockRepository) Update(d *Domain) error {
	return m.Called(d).Error(0)
}
func (m *MockRepository) Delete(id uint) error {
	return m.Called(id).Error(0)
}
func (m *MockRepository) HardDelete(id uint) error {
	return m.Called(id).Error(0)
}
func (m *MockRepository) ListByPage(limit, offset int, sortField, sortDirection, search string) ([]*Domain, int, error) {
	args := m.Called(limit, offset, sortField, sortDirection, search)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*Domain), args.Int(1), args.Error(2)
}
func (m *MockRepository) ListByCursor(cursorData *cursorInfo, limit int, sortField, sortDirection, search string) ([]*Domain, error) {
	args := m.Called(cursorData, limit, sortField, sortDirection, search)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Domain), args.Error(1)
}

// MockAuthService คือ "นักแสดงแทน" สำหรับ AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}
func (m *MockAuthService) ComparePassword(hashedPassword, plainPassword string) error {
	args := m.Called(hashedPassword, plainPassword)
	return args.Error(0)
}
func (m *MockAuthService) GenerateToken(userID uint, role string) (string, error) {
	args := m.Called(userID, role)
	return args.String(0), args.Error(1)
}
func (m *MockAuthService) ValidateToken(tokenString string) (*auth.JWTClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.JWTClaims), args.Error(1)
}

// ====================================================================================
// Test Cases (บทละครของเรา)
// ====================================================================================

// --- CreateUser Tests ---
func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockAuth := new(MockAuthService)
	mockLogger := logger.NewPrettyLogger()

	mockRepo.On("GetByEmail", "test@example.com").Return(nil, gorm.ErrRecordNotFound)
	mockAuth.On("HashPassword", "password123").Return("hashed_password_string", nil)
	mockRepo.On("Create", mock.AnythingOfType("*example_user.Domain")).Return(nil)

	userService := NewExampleUserService(mockRepo, mockAuth, mockLogger)
	userToCreate := &Domain{Name: "Test User", Email: "test@example.com"}

	createdUser, err := userService.CreateUser(userToCreate, "password123")

	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, uint(1), createdUser.ID)
	assert.Equal(t, "hashed_password_string", createdUser.PasswordHash)
	mockRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

func TestCreateUser_Fail_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(MockRepository)
	mockAuth := new(MockAuthService)
	mockLogger := logger.NewPrettyLogger()

	existingUser := &Domain{ID: 99, Email: "exists@example.com"}
	mockRepo.On("GetByEmail", "exists@example.com").Return(existingUser, nil)

	userService := NewExampleUserService(mockRepo, mockAuth, mockLogger)
	userToCreate := &Domain{Name: "Another User", Email: "exists@example.com"}

	createdUser, err := userService.CreateUser(userToCreate, "password123")

	assert.Error(t, err)
	assert.Nil(t, createdUser)
	var appErr *custom_errors.AppError
	if assert.ErrorAs(t, err, &appErr) {
		assert.Equal(t, custom_errors.ErrAlreadyExists, appErr.Code)
	}
	mockRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

// --- GetUserByID Tests ---
func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockAuth := new(MockAuthService)
	mockLogger := logger.NewPrettyLogger()

	expectedUser := &Domain{ID: 1, Name: "Nipon K.", Email: "nipon.k@example.com"}
	mockRepo.On("GetByID", uint(1)).Return(expectedUser, nil)

	userService := NewExampleUserService(mockRepo, mockAuth, mockLogger)
	foundUser, err := userService.GetUserByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, uint(1), foundUser.ID)
	mockRepo.AssertExpectations(t)
}

// --- Login Tests ---
func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockAuth := new(MockAuthService)
	mockLogger := logger.NewPrettyLogger()

	userInDB := &Domain{ID: 1, Role: "user", PasswordHash: "hashed_password"}
	mockRepo.On("GetByEmail", "test@example.com").Return(userInDB, nil)
	mockAuth.On("ComparePassword", "hashed_password", "password123").Return(nil)
	mockAuth.On("GenerateToken", uint(1), "user").Return("generated_jwt_token", nil)

	userService := NewExampleUserService(mockRepo, mockAuth, mockLogger)
	token, err := userService.Login("test@example.com", "password123")

	assert.NoError(t, err)
	assert.Equal(t, "generated_jwt_token", token)
	mockRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

func TestLogin_Fail_WrongPassword(t *testing.T) {
	mockRepo := new(MockRepository)
	mockAuth := new(MockAuthService)
	mockLogger := logger.NewPrettyLogger()

	userInDB := &Domain{ID: 1, Role: "user", PasswordHash: "hashed_password"}
	mockRepo.On("GetByEmail", "test@example.com").Return(userInDB, nil)
	mockAuth.On("ComparePassword", "hashed_password", "wrong_password").Return(errors.New("password mismatch"))

	userService := NewExampleUserService(mockRepo, mockAuth, mockLogger)
	token, err := userService.Login("test@example.com", "wrong_password")

	assert.Error(t, err)
	assert.Empty(t, token)
	var appErr *custom_errors.AppError
	if assert.ErrorAs(t, err, &appErr) {
		assert.Equal(t, custom_errors.ErrUnauthorized, appErr.Code)
	}
	mockRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}
