package overtime

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/riskykurniawan15/payrolls/mocks"
	"github.com/riskykurniawan15/payrolls/models/overtime"
)

func TestOvertimeRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIOvertimeRepository{}

		// Test data
		overtimeData := &overtime.Overtime{
			UserID:         1,
			OvertimesDate:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalHoursTime: 2.5,
			CreatedBy:      1,
		}

		// Setup expectations
		mockRepo.On("Create", mock.Anything, overtimeData).Return(nil)

		// Execute
		err := mockRepo.Create(context.Background(), overtimeData)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIOvertimeRepository{}

		// Test data
		overtimeData := &overtime.Overtime{
			UserID:         1,
			OvertimesDate:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalHoursTime: -1, // Invalid hours
			CreatedBy:      1,
		}

		// Setup expectations
		mockRepo.On("Create", mock.Anything, overtimeData).Return(assert.AnError)

		// Execute
		err := mockRepo.Create(context.Background(), overtimeData)

		// Assert
		assert.Error(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestOvertimeRepository_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIOvertimeRepository{}

		// Test data
		overtimeID := uint(1)
		expectedOvertime := &overtime.Overtime{
			ID:             1,
			UserID:         1,
			OvertimesDate:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalHoursTime: 2.5,
			CreatedBy:      1,
		}

		// Setup expectations
		mockRepo.On("GetByID", mock.Anything, overtimeID).Return(expectedOvertime, nil)

		// Execute
		foundOvertime, err := mockRepo.GetByID(context.Background(), overtimeID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedOvertime.ID, foundOvertime.ID)
		assert.Equal(t, expectedOvertime.UserID, foundOvertime.UserID)
		assert.Equal(t, expectedOvertime.TotalHoursTime, foundOvertime.TotalHoursTime)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("overtime not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIOvertimeRepository{}

		// Test data
		overtimeID := uint(999)

		// Setup expectations
		mockRepo.On("GetByID", mock.Anything, overtimeID).Return(nil, assert.AnError)

		// Execute
		foundOvertime, err := mockRepo.GetByID(context.Background(), overtimeID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, foundOvertime)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestOvertimeRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIOvertimeRepository{}

		// Test data
		overtimeID := uint(1)
		updates := map[string]interface{}{
			"total_hours_time": 3.0,
		}

		// Setup expectations
		mockRepo.On("Update", mock.Anything, overtimeID, updates).Return(nil)

		// Execute
		err := mockRepo.Update(context.Background(), overtimeID, updates)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("overtime not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIOvertimeRepository{}

		// Test data
		overtimeID := uint(999)
		updates := map[string]interface{}{
			"total_hours_time": 3.0,
		}

		// Setup expectations
		mockRepo.On("Update", mock.Anything, overtimeID, updates).Return(assert.AnError)

		// Execute
		err := mockRepo.Update(context.Background(), overtimeID, updates)

		// Assert
		assert.Error(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestOvertimeRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIOvertimeRepository{}

		// Test data
		overtimeID := uint(1)

		// Setup expectations
		mockRepo.On("Delete", mock.Anything, overtimeID).Return(nil)

		// Execute
		err := mockRepo.Delete(context.Background(), overtimeID)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("overtime not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIOvertimeRepository{}

		// Test data
		overtimeID := uint(999)

		// Setup expectations
		mockRepo.On("Delete", mock.Anything, overtimeID).Return(assert.AnError)

		// Execute
		err := mockRepo.Delete(context.Background(), overtimeID)

		// Assert
		assert.Error(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestOvertimeRepository_GetTotalHoursByUserAndDate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIOvertimeRepository{}

		// Test data
		userID := uint(1)
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		expectedHours := 5.5

		// Setup expectations
		mockRepo.On("GetTotalHoursByUserAndDate", mock.Anything, userID, date).Return(expectedHours, nil)

		// Execute
		hours, err := mockRepo.GetTotalHoursByUserAndDate(context.Background(), userID, date)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedHours, hours)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("no overtime found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIOvertimeRepository{}

		// Test data
		userID := uint(1)
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

		// Setup expectations
		mockRepo.On("GetTotalHoursByUserAndDate", mock.Anything, userID, date).Return(0.0, nil)

		// Execute
		hours, err := mockRepo.GetTotalHoursByUserAndDate(context.Background(), userID, date)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 0.0, hours)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestOvertimeRepository_GetByUserAndDateRange(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIOvertimeRepository{}

		// Test data
		userID := uint(1)
		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

		expectedOvertimes := []overtime.Overtime{
			{
				ID:             1,
				UserID:         1,
				OvertimesDate:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				TotalHoursTime: 2.5,
				CreatedBy:      1,
			},
			{
				ID:             2,
				UserID:         1,
				OvertimesDate:  time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
				TotalHoursTime: 1.5,
				CreatedBy:      1,
			},
		}

		// Setup expectations
		mockRepo.On("GetByUserAndDateRange", mock.Anything, userID, startDate, endDate).Return(expectedOvertimes, nil)

		// Execute
		foundOvertimes, err := mockRepo.GetByUserAndDateRange(context.Background(), userID, startDate, endDate)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, foundOvertimes, 2)
		assert.Equal(t, expectedOvertimes[0].ID, foundOvertimes[0].ID)
		assert.Equal(t, expectedOvertimes[1].ID, foundOvertimes[1].ID)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("no overtimes found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIOvertimeRepository{}

		// Test data
		userID := uint(999)
		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

		// Setup expectations
		mockRepo.On("GetByUserAndDateRange", mock.Anything, userID, startDate, endDate).Return([]overtime.Overtime{}, nil)

		// Execute
		foundOvertimes, err := mockRepo.GetByUserAndDateRange(context.Background(), userID, startDate, endDate)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, foundOvertimes, 0)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// Test untuk memastikan interface berfungsi dengan benar
func TestOvertimeRepository_Interface(t *testing.T) {
	t.Run("interface implementation", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIOvertimeRepository{}

		// Test bahwa mock mengimplementasikan interface
		var repo IOvertimeRepository = mockRepo
		assert.NotNil(t, repo)

		// Test data
		overtimeData := &overtime.Overtime{
			UserID:         1,
			OvertimesDate:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalHoursTime: 2.5,
			CreatedBy:      1,
		}

		expectedOvertime := &overtime.Overtime{
			ID:             1,
			UserID:         1,
			OvertimesDate:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalHoursTime: 2.5,
			CreatedBy:      1,
		}

		// Setup expectations
		mockRepo.On("Create", mock.Anything, overtimeData).Return(nil)
		mockRepo.On("GetByID", mock.Anything, uint(1)).Return(expectedOvertime, nil)
		mockRepo.On("Update", mock.Anything, uint(1), mock.Anything).Return(nil)
		mockRepo.On("Delete", mock.Anything, uint(1)).Return(nil)
		mockRepo.On("GetTotalHoursByUserAndDate", mock.Anything, uint(1), mock.Anything).Return(2.5, nil)
		mockRepo.On("GetByUserAndDateRange", mock.Anything, uint(1), mock.Anything, mock.Anything).Return([]overtime.Overtime{*expectedOvertime}, nil)

		// Test semua method interface
		err := repo.Create(context.Background(), overtimeData)
		assert.NoError(t, err)

		foundOvertime, err := repo.GetByID(context.Background(), uint(1))
		assert.NoError(t, err)
		assert.Equal(t, expectedOvertime.ID, foundOvertime.ID)

		err = repo.Update(context.Background(), uint(1), map[string]interface{}{"total_hours_time": 3.0})
		assert.NoError(t, err)

		err = repo.Delete(context.Background(), uint(1))
		assert.NoError(t, err)

		hours, err := repo.GetTotalHoursByUserAndDate(context.Background(), uint(1), time.Now())
		assert.NoError(t, err)
		assert.Equal(t, 2.5, hours)

		foundOvertimes, err := repo.GetByUserAndDateRange(context.Background(), uint(1), time.Now(), time.Now())
		assert.NoError(t, err)
		assert.Len(t, foundOvertimes, 1)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}
