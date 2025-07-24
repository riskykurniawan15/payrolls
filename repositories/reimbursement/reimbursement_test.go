package reimbursement

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/riskykurniawan15/payrolls/mocks"
	"github.com/riskykurniawan15/payrolls/models/reimbursement"
)

func TestReimbursementRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIReimbursementRepository{}

		// Test data
		reimbursementData := &reimbursement.Reimbursement{
			UserID:    1,
			Title:     "Transportation",
			Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Amount:    50000,
			CreatedBy: 1,
		}

		// Setup expectations
		mockRepo.On("Create", mock.Anything, reimbursementData).Return(nil)

		// Execute
		err := mockRepo.Create(context.Background(), reimbursementData)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIReimbursementRepository{}

		// Test data
		reimbursementData := &reimbursement.Reimbursement{
			UserID:    1,
			Title:     "", // Empty title
			Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Amount:    -1000, // Negative amount
			CreatedBy: 1,
		}

		// Setup expectations
		mockRepo.On("Create", mock.Anything, reimbursementData).Return(assert.AnError)

		// Execute
		err := mockRepo.Create(context.Background(), reimbursementData)

		// Assert
		assert.Error(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestReimbursementRepository_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIReimbursementRepository{}

		// Test data
		reimbursementID := uint(1)
		expectedReimbursement := &reimbursement.Reimbursement{
			ID:        1,
			UserID:    1,
			Title:     "Transportation",
			Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Amount:    50000,
			CreatedBy: 1,
		}

		// Setup expectations
		mockRepo.On("GetByID", mock.Anything, reimbursementID).Return(expectedReimbursement, nil)

		// Execute
		foundReimbursement, err := mockRepo.GetByID(context.Background(), reimbursementID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedReimbursement.ID, foundReimbursement.ID)
		assert.Equal(t, expectedReimbursement.UserID, foundReimbursement.UserID)
		assert.Equal(t, expectedReimbursement.Title, foundReimbursement.Title)
		assert.Equal(t, expectedReimbursement.Amount, foundReimbursement.Amount)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("reimbursement not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIReimbursementRepository{}

		// Test data
		reimbursementID := uint(999)

		// Setup expectations
		mockRepo.On("GetByID", mock.Anything, reimbursementID).Return(nil, assert.AnError)

		// Execute
		foundReimbursement, err := mockRepo.GetByID(context.Background(), reimbursementID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, foundReimbursement)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestReimbursementRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIReimbursementRepository{}

		// Test data
		reimbursementID := uint(1)
		updates := map[string]interface{}{
			"title":  "Transportation Updated",
			"amount": 60000,
		}

		// Setup expectations
		mockRepo.On("Update", mock.Anything, reimbursementID, updates).Return(nil)

		// Execute
		err := mockRepo.Update(context.Background(), reimbursementID, updates)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("reimbursement not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIReimbursementRepository{}

		// Test data
		reimbursementID := uint(999)
		updates := map[string]interface{}{
			"title": "Non-existent Reimbursement",
		}

		// Setup expectations
		mockRepo.On("Update", mock.Anything, reimbursementID, updates).Return(assert.AnError)

		// Execute
		err := mockRepo.Update(context.Background(), reimbursementID, updates)

		// Assert
		assert.Error(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestReimbursementRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIReimbursementRepository{}

		// Test data
		reimbursementID := uint(1)

		// Setup expectations
		mockRepo.On("Delete", mock.Anything, reimbursementID).Return(nil)

		// Execute
		err := mockRepo.Delete(context.Background(), reimbursementID)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("reimbursement not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIReimbursementRepository{}

		// Test data
		reimbursementID := uint(999)

		// Setup expectations
		mockRepo.On("Delete", mock.Anything, reimbursementID).Return(assert.AnError)

		// Execute
		err := mockRepo.Delete(context.Background(), reimbursementID)

		// Assert
		assert.Error(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestReimbursementRepository_GetByUserAndDateRange(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIReimbursementRepository{}

		// Test data
		userID := uint(1)
		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

		expectedReimbursements := []reimbursement.Reimbursement{
			{
				ID:        1,
				UserID:    1,
				Title:     "Transportation",
				Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				Amount:    50000,
				CreatedBy: 1,
			},
			{
				ID:        2,
				UserID:    1,
				Title:     "Meal",
				Date:      time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
				Amount:    25000,
				CreatedBy: 1,
			},
		}

		// Setup expectations
		mockRepo.On("GetByUserAndDateRange", mock.Anything, userID, startDate, endDate).Return(expectedReimbursements, nil)

		// Execute
		foundReimbursements, err := mockRepo.GetByUserAndDateRange(context.Background(), userID, startDate, endDate)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, foundReimbursements, 2)
		assert.Equal(t, expectedReimbursements[0].ID, foundReimbursements[0].ID)
		assert.Equal(t, expectedReimbursements[1].ID, foundReimbursements[1].ID)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("no reimbursements found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIReimbursementRepository{}

		// Test data
		userID := uint(999)
		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

		// Setup expectations
		mockRepo.On("GetByUserAndDateRange", mock.Anything, userID, startDate, endDate).Return([]reimbursement.Reimbursement{}, nil)

		// Execute
		foundReimbursements, err := mockRepo.GetByUserAndDateRange(context.Background(), userID, startDate, endDate)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, foundReimbursements, 0)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// Test untuk memastikan interface berfungsi dengan benar
func TestReimbursementRepository_Interface(t *testing.T) {
	t.Run("interface implementation", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIReimbursementRepository{}

		// Test bahwa mock mengimplementasikan interface
		var repo IReimbursementRepository = mockRepo
		assert.NotNil(t, repo)

		// Test data
		reimbursementData := &reimbursement.Reimbursement{
			UserID:    1,
			Title:     "Transportation",
			Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Amount:    50000,
			CreatedBy: 1,
		}

		expectedReimbursement := &reimbursement.Reimbursement{
			ID:        1,
			UserID:    1,
			Title:     "Transportation",
			Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Amount:    50000,
			CreatedBy: 1,
		}

		// Setup expectations
		mockRepo.On("Create", mock.Anything, reimbursementData).Return(nil)
		mockRepo.On("GetByID", mock.Anything, uint(1)).Return(expectedReimbursement, nil)
		mockRepo.On("Update", mock.Anything, uint(1), mock.Anything).Return(nil)
		mockRepo.On("Delete", mock.Anything, uint(1)).Return(nil)
		mockRepo.On("GetByUserAndDateRange", mock.Anything, uint(1), mock.Anything, mock.Anything).Return([]reimbursement.Reimbursement{*expectedReimbursement}, nil)

		// Test semua method interface
		err := repo.Create(context.Background(), reimbursementData)
		assert.NoError(t, err)

		foundReimbursement, err := repo.GetByID(context.Background(), uint(1))
		assert.NoError(t, err)
		assert.Equal(t, expectedReimbursement.ID, foundReimbursement.ID)

		err = repo.Update(context.Background(), uint(1), map[string]interface{}{"title": "Updated"})
		assert.NoError(t, err)

		err = repo.Delete(context.Background(), uint(1))
		assert.NoError(t, err)

		foundReimbursements, err := repo.GetByUserAndDateRange(context.Background(), uint(1), time.Now(), time.Now())
		assert.NoError(t, err)
		assert.Len(t, foundReimbursements, 1)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}
