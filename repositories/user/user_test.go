package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/riskykurniawan15/payrolls/mocks"
	"github.com/riskykurniawan15/payrolls/models/user"
)

func TestUserRepository_CreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIUserRepository{}

		// Test data
		userData := user.User{
			Username: "testuser",
			Password: "hashedpassword",
			Role:     "employee",
			Salary:   5000000,
		}

		expectedUser := user.User{
			ID:       1,
			Username: "testuser",
			Password: "hashedpassword",
			Role:     "employee",
			Salary:   5000000,
		}

		// Setup expectations
		mockRepo.On("CreateUser", mock.Anything, userData).Return(expectedUser, nil)

		// Execute
		createdUser, err := mockRepo.CreateUser(context.Background(), userData)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, createdUser.ID)
		assert.Equal(t, expectedUser.Username, createdUser.Username)
		assert.Equal(t, expectedUser.Role, createdUser.Role)
		assert.Equal(t, expectedUser.Salary, createdUser.Salary)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIUserRepository{}

		// Test data
		userData := user.User{
			Username: "duplicateuser",
			Password: "hashedpassword",
			Role:     "employee",
			Salary:   5000000,
		}

		// Setup expectations
		mockRepo.On("CreateUser", mock.Anything, userData).Return(user.User{}, assert.AnError)

		// Execute
		createdUser, err := mockRepo.CreateUser(context.Background(), userData)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, user.User{}, createdUser)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestUserRepository_GetUserByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIUserRepository{}

		// Test data
		userID := uint(1)
		expectedUser := user.User{
			ID:       1,
			Username: "testuser",
			Password: "hashedpassword",
			Role:     "employee",
			Salary:   5000000,
		}

		// Setup expectations
		mockRepo.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)

		// Execute
		foundUser, err := mockRepo.GetUserByID(context.Background(), userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, foundUser.ID)
		assert.Equal(t, expectedUser.Username, foundUser.Username)
		assert.Equal(t, expectedUser.Role, foundUser.Role)
		assert.Equal(t, expectedUser.Salary, foundUser.Salary)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIUserRepository{}

		// Test data
		userID := uint(999)

		// Setup expectations
		mockRepo.On("GetUserByID", mock.Anything, userID).Return(user.User{}, assert.AnError)

		// Execute
		foundUser, err := mockRepo.GetUserByID(context.Background(), userID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, user.User{}, foundUser)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIUserRepository{}

		// Test data
		userID := uint(0)

		// Setup expectations
		mockRepo.On("GetUserByID", mock.Anything, userID).Return(user.User{}, assert.AnError)

		// Execute
		foundUser, err := mockRepo.GetUserByID(context.Background(), userID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, user.User{}, foundUser)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestUserRepository_GetUserByUsername(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIUserRepository{}

		// Test data
		username := "testuser"
		expectedUser := user.User{
			ID:       1,
			Username: "testuser",
			Password: "hashedpassword",
			Role:     "employee",
			Salary:   5000000,
		}

		// Setup expectations
		mockRepo.On("GetUserByUsername", mock.Anything, username).Return(expectedUser, nil)

		// Execute
		foundUser, err := mockRepo.GetUserByUsername(context.Background(), username)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, foundUser.ID)
		assert.Equal(t, expectedUser.Username, foundUser.Username)
		assert.Equal(t, expectedUser.Role, foundUser.Role)
		assert.Equal(t, expectedUser.Salary, foundUser.Salary)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIUserRepository{}

		// Test data
		username := "nonexistentuser"

		// Setup expectations
		mockRepo.On("GetUserByUsername", mock.Anything, username).Return(user.User{}, assert.AnError)

		// Execute
		foundUser, err := mockRepo.GetUserByUsername(context.Background(), username)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, user.User{}, foundUser)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty username", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIUserRepository{}

		// Test data
		username := ""

		// Setup expectations
		mockRepo.On("GetUserByUsername", mock.Anything, username).Return(user.User{}, assert.AnError)

		// Execute
		foundUser, err := mockRepo.GetUserByUsername(context.Background(), username)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, user.User{}, foundUser)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("case insensitive", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIUserRepository{}

		// Test data
		username := "TESTUSER"
		expectedUser := user.User{
			ID:       1,
			Username: "testuser",
			Password: "hashedpassword",
			Role:     "employee",
			Salary:   5000000,
		}

		// Setup expectations
		mockRepo.On("GetUserByUsername", mock.Anything, username).Return(expectedUser, nil)

		// Execute
		foundUser, err := mockRepo.GetUserByUsername(context.Background(), username)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, foundUser.ID)
		assert.Equal(t, expectedUser.Username, foundUser.Username)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// Test untuk memastikan interface berfungsi dengan benar
func TestUserRepository_Interface(t *testing.T) {
	t.Run("interface implementation", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIUserRepository{}

		// Test bahwa mock mengimplementasikan interface
		var repo IUserRepository = mockRepo
		assert.NotNil(t, repo)

		// Test data
		userData := user.User{
			Username: "testuser",
			Password: "hashedpassword",
			Role:     "employee",
			Salary:   5000000,
		}

		expectedUser := user.User{
			ID:       1,
			Username: "testuser",
			Password: "hashedpassword",
			Role:     "employee",
			Salary:   5000000,
		}

		// Setup expectations
		mockRepo.On("CreateUser", mock.Anything, userData).Return(expectedUser, nil)
		mockRepo.On("GetUserByID", mock.Anything, uint(1)).Return(expectedUser, nil)
		mockRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(expectedUser, nil)

		// Test semua method interface
		createdUser, err := repo.CreateUser(context.Background(), userData)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, createdUser.ID)

		foundUser, err := repo.GetUserByID(context.Background(), uint(1))
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, foundUser.ID)

		foundUserByUsername, err := repo.GetUserByUsername(context.Background(), "testuser")
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, foundUserByUsername.ID)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// Test untuk error handling
func TestUserRepository_ErrorHandling(t *testing.T) {
	t.Run("create user error", func(t *testing.T) {
		mockRepo := &mocks.MockIUserRepository{}

		userData := user.User{
			Username: "testuser",
			Password: "hashedpassword",
			Role:     "employee",
			Salary:   5000000,
		}

		mockRepo.On("CreateUser", mock.Anything, userData).Return(user.User{}, assert.AnError)

		createdUser, err := mockRepo.CreateUser(context.Background(), userData)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, createdUser)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get user by ID error", func(t *testing.T) {
		mockRepo := &mocks.MockIUserRepository{}

		mockRepo.On("GetUserByID", mock.Anything, uint(999)).Return(user.User{}, assert.AnError)

		foundUser, err := mockRepo.GetUserByID(context.Background(), uint(999))

		assert.Error(t, err)
		assert.Equal(t, user.User{}, foundUser)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get user by username error", func(t *testing.T) {
		mockRepo := &mocks.MockIUserRepository{}

		mockRepo.On("GetUserByUsername", mock.Anything, "nonexistent").Return(user.User{}, assert.AnError)

		foundUser, err := mockRepo.GetUserByUsername(context.Background(), "nonexistent")

		assert.Error(t, err)
		assert.Equal(t, user.User{}, foundUser)
		mockRepo.AssertExpectations(t)
	})
}
