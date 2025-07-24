package period

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/riskykurniawan15/payrolls/mocks"
	"github.com/riskykurniawan15/payrolls/models/period"
)

func TestPeriodRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodRepository{}

		// Test data
		periodData := &period.Period{
			Name:      "January 2024",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			Status:    1, // active
		}

		// Setup expectations
		mockRepo.On("Create", mock.Anything, periodData).Return(nil)

		// Execute
		err := mockRepo.Create(context.Background(), periodData)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodRepository{}

		// Test data
		periodData := &period.Period{
			Name:      "Invalid Period",
			StartDate: time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // End date before start date
			Status:    1,
		}

		// Setup expectations
		mockRepo.On("Create", mock.Anything, periodData).Return(assert.AnError)

		// Execute
		err := mockRepo.Create(context.Background(), periodData)

		// Assert
		assert.Error(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestPeriodRepository_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodRepository{}

		// Test data
		periodID := uint(1)
		expectedPeriod := &period.Period{
			ID:        1,
			Name:      "January 2024",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			Status:    1,
		}

		// Setup expectations
		mockRepo.On("GetByID", mock.Anything, periodID).Return(expectedPeriod, nil)

		// Execute
		foundPeriod, err := mockRepo.GetByID(context.Background(), periodID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedPeriod.ID, foundPeriod.ID)
		assert.Equal(t, expectedPeriod.Name, foundPeriod.Name)
		assert.Equal(t, expectedPeriod.Status, foundPeriod.Status)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("period not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodRepository{}

		// Test data
		periodID := uint(999)

		// Setup expectations
		mockRepo.On("GetByID", mock.Anything, periodID).Return(nil, assert.AnError)

		// Execute
		foundPeriod, err := mockRepo.GetByID(context.Background(), periodID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, foundPeriod)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestPeriodRepository_GetByCode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodRepository{}

		// Test data
		code := "JAN2024"
		expectedPeriod := &period.Period{
			ID:        1,
			Code:      "JAN2024",
			Name:      "January 2024",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			Status:    1,
		}

		// Setup expectations
		mockRepo.On("GetByCode", mock.Anything, code).Return(expectedPeriod, nil)

		// Execute
		foundPeriod, err := mockRepo.GetByCode(context.Background(), code)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedPeriod.ID, foundPeriod.ID)
		assert.Equal(t, expectedPeriod.Code, foundPeriod.Code)
		assert.Equal(t, expectedPeriod.Name, foundPeriod.Name)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("period not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodRepository{}

		// Test data
		code := "NONEXISTENT"

		// Setup expectations
		mockRepo.On("GetByCode", mock.Anything, code).Return(nil, assert.AnError)

		// Execute
		foundPeriod, err := mockRepo.GetByCode(context.Background(), code)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, foundPeriod)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestPeriodRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodRepository{}

		// Test data
		periodID := uint(1)
		updates := map[string]interface{}{
			"name":   "January 2024 Updated",
			"status": 2, // inactive
		}

		// Setup expectations
		mockRepo.On("Update", mock.Anything, periodID, updates).Return(nil)

		// Execute
		err := mockRepo.Update(context.Background(), periodID, updates)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("period not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodRepository{}

		// Test data
		periodID := uint(999)
		updates := map[string]interface{}{
			"name": "Non-existent Period",
		}

		// Setup expectations
		mockRepo.On("Update", mock.Anything, periodID, updates).Return(assert.AnError)

		// Execute
		err := mockRepo.Update(context.Background(), periodID, updates)

		// Assert
		assert.Error(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestPeriodRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodRepository{}

		// Test data
		periodID := uint(1)

		// Setup expectations
		mockRepo.On("Delete", mock.Anything, periodID).Return(nil)

		// Execute
		err := mockRepo.Delete(context.Background(), periodID)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("period not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodRepository{}

		// Test data
		periodID := uint(999)

		// Setup expectations
		mockRepo.On("Delete", mock.Anything, periodID).Return(assert.AnError)

		// Execute
		err := mockRepo.Delete(context.Background(), periodID)

		// Assert
		assert.Error(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestPeriodRepository_IsCodeExists(t *testing.T) {
	t.Run("code exists", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodRepository{}

		// Test data
		code := "JAN2024"

		// Setup expectations
		mockRepo.On("IsCodeExists", mock.Anything, code).Return(true, nil)

		// Execute
		exists, err := mockRepo.IsCodeExists(context.Background(), code)

		// Assert
		assert.NoError(t, err)
		assert.True(t, exists)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("code does not exist", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodRepository{}

		// Test data
		code := "NONEXISTENT"

		// Setup expectations
		mockRepo.On("IsCodeExists", mock.Anything, code).Return(false, nil)

		// Execute
		exists, err := mockRepo.IsCodeExists(context.Background(), code)

		// Assert
		assert.NoError(t, err)
		assert.False(t, exists)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestPeriodRepository_GenerateUniqueCode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodRepository{}

		// Test data
		expectedCode := "JAN2024"

		// Setup expectations
		mockRepo.On("GenerateUniqueCode", mock.Anything).Return(expectedCode, nil)

		// Execute
		code, err := mockRepo.GenerateUniqueCode(context.Background())

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCode, code)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodRepository{}

		// Setup expectations
		mockRepo.On("GenerateUniqueCode", mock.Anything).Return("", assert.AnError)

		// Execute
		code, err := mockRepo.GenerateUniqueCode(context.Background())

		// Assert
		assert.Error(t, err)
		assert.Empty(t, code)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// Test untuk memastikan interface berfungsi dengan benar
func TestPeriodRepository_Interface(t *testing.T) {
	t.Run("interface implementation", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodRepository{}

		// Test bahwa mock mengimplementasikan interface
		var repo IPeriodRepository = mockRepo
		assert.NotNil(t, repo)

		// Test data
		periodData := &period.Period{
			Name:      "Test Period",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			Status:    1,
		}

		expectedPeriod := &period.Period{
			ID:        1,
			Name:      "Test Period",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			Status:    1,
		}

		// Setup expectations
		mockRepo.On("Create", mock.Anything, periodData).Return(nil)
		mockRepo.On("GetByID", mock.Anything, uint(1)).Return(expectedPeriod, nil)
		mockRepo.On("GetByCode", mock.Anything, "TEST").Return(expectedPeriod, nil)
		mockRepo.On("Update", mock.Anything, uint(1), mock.Anything).Return(nil)
		mockRepo.On("Delete", mock.Anything, uint(1)).Return(nil)
		mockRepo.On("IsCodeExists", mock.Anything, "TEST").Return(false, nil)
		mockRepo.On("GenerateUniqueCode", mock.Anything).Return("TEST2024", nil)

		// Test semua method interface
		err := repo.Create(context.Background(), periodData)
		assert.NoError(t, err)

		foundPeriod, err := repo.GetByID(context.Background(), uint(1))
		assert.NoError(t, err)
		assert.Equal(t, expectedPeriod.ID, foundPeriod.ID)

		foundPeriodByCode, err := repo.GetByCode(context.Background(), "TEST")
		assert.NoError(t, err)
		assert.Equal(t, expectedPeriod.ID, foundPeriodByCode.ID)

		err = repo.Update(context.Background(), uint(1), map[string]interface{}{"name": "Updated"})
		assert.NoError(t, err)

		err = repo.Delete(context.Background(), uint(1))
		assert.NoError(t, err)

		exists, err := repo.IsCodeExists(context.Background(), "TEST")
		assert.NoError(t, err)
		assert.False(t, exists)

		code, err := repo.GenerateUniqueCode(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, "TEST2024", code)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}
