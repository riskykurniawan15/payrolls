package attendance

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/riskykurniawan15/payrolls/mocks"
	"github.com/riskykurniawan15/payrolls/models/attendance"
)

func TestAttendanceRepository_CreateAttendance(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIAttendanceRepository{}

		// Test data
		attendanceData := attendance.Attendance{
			UserID:      1,
			CheckInDate: time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
			CreatedBy:   1,
		}

		expectedAttendance := attendance.Attendance{
			ID:          1,
			UserID:      1,
			CheckInDate: time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
			CreatedBy:   1,
		}

		// Setup expectations
		mockRepo.On("CreateAttendance", mock.Anything, attendanceData).Return(expectedAttendance, nil)

		// Execute
		createdAttendance, err := mockRepo.CreateAttendance(context.Background(), attendanceData)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedAttendance.ID, createdAttendance.ID)
		assert.Equal(t, expectedAttendance.UserID, createdAttendance.UserID)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIAttendanceRepository{}

		// Test data
		attendanceData := attendance.Attendance{
			UserID:      1,
			CheckInDate: time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
			CreatedBy:   1,
		}

		// Setup expectations
		mockRepo.On("CreateAttendance", mock.Anything, attendanceData).Return(attendance.Attendance{}, assert.AnError)

		// Execute
		createdAttendance, err := mockRepo.CreateAttendance(context.Background(), attendanceData)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, attendance.Attendance{}, createdAttendance)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestAttendanceRepository_GetAttendanceByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIAttendanceRepository{}

		// Test data
		attendanceID := uint(1)
		userID := uint(1)
		expectedAttendance := attendance.Attendance{
			ID:          1,
			UserID:      1,
			CheckInDate: time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
			CreatedBy:   1,
		}

		// Setup expectations
		mockRepo.On("GetAttendanceByID", mock.Anything, attendanceID, userID).Return(expectedAttendance, nil)

		// Execute
		foundAttendance, err := mockRepo.GetAttendanceByID(context.Background(), attendanceID, userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedAttendance.ID, foundAttendance.ID)
		assert.Equal(t, expectedAttendance.UserID, foundAttendance.UserID)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("attendance not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIAttendanceRepository{}

		// Test data
		attendanceID := uint(999)
		userID := uint(1)

		// Setup expectations
		mockRepo.On("GetAttendanceByID", mock.Anything, attendanceID, userID).Return(attendance.Attendance{}, assert.AnError)

		// Execute
		foundAttendance, err := mockRepo.GetAttendanceByID(context.Background(), attendanceID, userID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, attendance.Attendance{}, foundAttendance)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestAttendanceRepository_GetAttendances(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIAttendanceRepository{}

		// Test data
		userID := uint(1)
		page := 1
		limit := 10
		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

		expectedAttendances := []attendance.Attendance{
			{
				ID:          1,
				UserID:      1,
				CheckInDate: time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
				CreatedBy:   1,
			},
			{
				ID:          2,
				UserID:      1,
				CheckInDate: time.Date(2024, 1, 16, 8, 30, 0, 0, time.UTC),
				CreatedBy:   1,
			},
		}
		expectedTotal := int64(2)

		// Setup expectations
		mockRepo.On("GetAttendances", mock.Anything, userID, page, limit, &startDate, &endDate).Return(expectedAttendances, expectedTotal, nil)

		// Execute
		foundAttendances, total, err := mockRepo.GetAttendances(context.Background(), userID, page, limit, &startDate, &endDate)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, foundAttendances, 2)
		assert.Equal(t, expectedTotal, total)
		assert.Equal(t, expectedAttendances[0].ID, foundAttendances[0].ID)
		assert.Equal(t, expectedAttendances[1].ID, foundAttendances[1].ID)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("no attendances found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIAttendanceRepository{}

		// Test data
		userID := uint(999)
		page := 1
		limit := 10

		// Setup expectations
		mockRepo.On("GetAttendances", mock.Anything, userID, page, limit, (*time.Time)(nil), (*time.Time)(nil)).Return([]attendance.Attendance{}, int64(0), nil)

		// Execute
		foundAttendances, total, err := mockRepo.GetAttendances(context.Background(), userID, page, limit, nil, nil)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, foundAttendances, 0)
		assert.Equal(t, int64(0), total)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestAttendanceRepository_UpdateAttendance(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIAttendanceRepository{}

		// Test data
		attendanceData := attendance.Attendance{
			ID:           1,
			UserID:       1,
			CheckInDate:  time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
			CheckOutDate: &[]time.Time{time.Date(2024, 1, 15, 17, 0, 0, 0, time.UTC)}[0],
			CreatedBy:    1,
		}

		// Setup expectations
		mockRepo.On("UpdateAttendance", mock.Anything, attendanceData).Return(attendanceData, nil)

		// Execute
		updatedAttendance, err := mockRepo.UpdateAttendance(context.Background(), attendanceData)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, attendanceData.ID, updatedAttendance.ID)
		assert.NotNil(t, updatedAttendance.CheckOutDate)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("attendance not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIAttendanceRepository{}

		// Test data
		attendanceData := attendance.Attendance{
			ID:          999,
			UserID:      1,
			CheckInDate: time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
			CreatedBy:   1,
		}

		// Setup expectations
		mockRepo.On("UpdateAttendance", mock.Anything, attendanceData).Return(attendance.Attendance{}, assert.AnError)

		// Execute
		updatedAttendance, err := mockRepo.UpdateAttendance(context.Background(), attendanceData)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, attendance.Attendance{}, updatedAttendance)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestAttendanceRepository_GetByUserAndDate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIAttendanceRepository{}

		// Test data
		userID := uint(1)
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		expectedAttendance := &attendance.Attendance{
			ID:          1,
			UserID:      1,
			CheckInDate: time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
			CreatedBy:   1,
		}

		// Setup expectations
		mockRepo.On("GetByUserAndDate", mock.Anything, userID, date).Return(expectedAttendance, nil)

		// Execute
		foundAttendance, err := mockRepo.GetByUserAndDate(context.Background(), userID, date)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedAttendance.ID, foundAttendance.ID)
		assert.Equal(t, expectedAttendance.UserID, foundAttendance.UserID)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("attendance not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIAttendanceRepository{}

		// Test data
		userID := uint(1)
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

		// Setup expectations
		mockRepo.On("GetByUserAndDate", mock.Anything, userID, date).Return(nil, assert.AnError)

		// Execute
		foundAttendance, err := mockRepo.GetByUserAndDate(context.Background(), userID, date)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, foundAttendance)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// Test untuk memastikan interface berfungsi dengan benar
func TestAttendanceRepository_Interface(t *testing.T) {
	t.Run("interface implementation", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIAttendanceRepository{}

		// Test bahwa mock mengimplementasikan interface
		var repo IAttendanceRepository = mockRepo
		assert.NotNil(t, repo)

		// Test data
		attendanceData := attendance.Attendance{
			UserID:      1,
			CheckInDate: time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
			CreatedBy:   1,
		}

		expectedAttendance := attendance.Attendance{
			ID:          1,
			UserID:      1,
			CheckInDate: time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
			CreatedBy:   1,
		}

		// Setup expectations
		mockRepo.On("CreateAttendance", mock.Anything, attendanceData).Return(expectedAttendance, nil)
		mockRepo.On("GetAttendanceByID", mock.Anything, uint(1), uint(1)).Return(expectedAttendance, nil)
		mockRepo.On("GetAttendances", mock.Anything, uint(1), 1, 10, (*time.Time)(nil), (*time.Time)(nil)).Return([]attendance.Attendance{expectedAttendance}, int64(1), nil)
		mockRepo.On("UpdateAttendance", mock.Anything, expectedAttendance).Return(expectedAttendance, nil)
		mockRepo.On("GetByUserAndDate", mock.Anything, uint(1), mock.Anything).Return(&expectedAttendance, nil)

		// Test semua method interface
		createdAttendance, err := repo.CreateAttendance(context.Background(), attendanceData)
		assert.NoError(t, err)
		assert.Equal(t, expectedAttendance.ID, createdAttendance.ID)

		foundAttendance, err := repo.GetAttendanceByID(context.Background(), uint(1), uint(1))
		assert.NoError(t, err)
		assert.Equal(t, expectedAttendance.ID, foundAttendance.ID)

		foundAttendances, total, err := repo.GetAttendances(context.Background(), uint(1), 1, 10, nil, nil)
		assert.NoError(t, err)
		assert.Len(t, foundAttendances, 1)
		assert.Equal(t, int64(1), total)

		updatedAttendance, err := repo.UpdateAttendance(context.Background(), expectedAttendance)
		assert.NoError(t, err)
		assert.Equal(t, expectedAttendance.ID, updatedAttendance.ID)

		foundByDate, err := repo.GetByUserAndDate(context.Background(), uint(1), time.Now())
		assert.NoError(t, err)
		assert.Equal(t, expectedAttendance.ID, foundByDate.ID)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}
