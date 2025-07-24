package period_detail

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/riskykurniawan15/payrolls/mocks"
	"github.com/riskykurniawan15/payrolls/models/payslip"
	"github.com/riskykurniawan15/payrolls/models/period_detail"
)

func TestPeriodDetailRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodDetailData := &period_detail.PeriodDetail{
			PeriodsID:           1,
			UserID:              1,
			DailyRate:           100000.00,
			TotalWorking:        22,
			AmountSalary:        2200000.00,
			AmountOvertime:      500000.00,
			AmountReimbursement: 200000.00,
			TakeHomePay:         2900000.00,
			CreatedBy:           1,
		}

		// Setup expectations
		mockRepo.On("Create", mock.Anything, periodDetailData).Return(nil)

		// Execute
		err := mockRepo.Create(context.Background(), periodDetailData)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodDetailData := &period_detail.PeriodDetail{
			PeriodsID: 1,
			UserID:    1,
		}

		// Setup expectations
		mockRepo.On("Create", mock.Anything, periodDetailData).Return(assert.AnError)

		// Execute
		err := mockRepo.Create(context.Background(), periodDetailData)

		// Assert
		assert.Error(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestPeriodDetailRepository_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodDetailID := uint(1)
		expectedPeriodDetail := &period_detail.PeriodDetail{
			ID:                  1,
			PeriodsID:           1,
			UserID:              1,
			DailyRate:           100000.00,
			TotalWorking:        22,
			AmountSalary:        2200000.00,
			AmountOvertime:      500000.00,
			AmountReimbursement: 200000.00,
			TakeHomePay:         2900000.00,
			CreatedBy:           1,
			CreatedAt:           time.Now(),
		}

		// Setup expectations
		mockRepo.On("GetByID", mock.Anything, periodDetailID).Return(expectedPeriodDetail, nil)

		// Execute
		foundPeriodDetail, err := mockRepo.GetByID(context.Background(), periodDetailID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedPeriodDetail.ID, foundPeriodDetail.ID)
		assert.Equal(t, expectedPeriodDetail.PeriodsID, foundPeriodDetail.PeriodsID)
		assert.Equal(t, expectedPeriodDetail.UserID, foundPeriodDetail.UserID)
		assert.Equal(t, expectedPeriodDetail.DailyRate, foundPeriodDetail.DailyRate)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("period detail not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodDetailID := uint(999)

		// Setup expectations
		mockRepo.On("GetByID", mock.Anything, periodDetailID).Return(nil, assert.AnError)

		// Execute
		foundPeriodDetail, err := mockRepo.GetByID(context.Background(), periodDetailID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, foundPeriodDetail)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestPeriodDetailRepository_GetByPeriodAndUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodID := uint(1)
		userID := uint(1)
		expectedPeriodDetail := &period_detail.PeriodDetail{
			ID:                  1,
			PeriodsID:           periodID,
			UserID:              userID,
			DailyRate:           100000.00,
			TotalWorking:        22,
			AmountSalary:        2200000.00,
			AmountOvertime:      500000.00,
			AmountReimbursement: 200000.00,
			TakeHomePay:         2900000.00,
			CreatedBy:           1,
			CreatedAt:           time.Now(),
		}

		// Setup expectations
		mockRepo.On("GetByPeriodAndUser", mock.Anything, periodID, userID).Return(expectedPeriodDetail, nil)

		// Execute
		foundPeriodDetail, err := mockRepo.GetByPeriodAndUser(context.Background(), periodID, userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedPeriodDetail.ID, foundPeriodDetail.ID)
		assert.Equal(t, expectedPeriodDetail.PeriodsID, foundPeriodDetail.PeriodsID)
		assert.Equal(t, expectedPeriodDetail.UserID, foundPeriodDetail.UserID)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("period detail not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodID := uint(999)
		userID := uint(999)

		// Setup expectations
		mockRepo.On("GetByPeriodAndUser", mock.Anything, periodID, userID).Return(nil, assert.AnError)

		// Execute
		foundPeriodDetail, err := mockRepo.GetByPeriodAndUser(context.Background(), periodID, userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, foundPeriodDetail)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestPeriodDetailRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodDetailID := uint(1)
		updates := map[string]interface{}{
			"daily_rate":    120000.00,
			"total_working": 25,
		}

		// Setup expectations
		mockRepo.On("Update", mock.Anything, periodDetailID, updates).Return(nil)

		// Execute
		err := mockRepo.Update(context.Background(), periodDetailID, updates)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("period detail not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodDetailID := uint(999)
		updates := map[string]interface{}{
			"daily_rate": 120000.00,
		}

		// Setup expectations
		mockRepo.On("Update", mock.Anything, periodDetailID, updates).Return(assert.AnError)

		// Execute
		err := mockRepo.Update(context.Background(), periodDetailID, updates)

		// Assert
		assert.Error(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestPeriodDetailRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodDetailID := uint(1)

		// Setup expectations
		mockRepo.On("Delete", mock.Anything, periodDetailID).Return(nil)

		// Execute
		err := mockRepo.Delete(context.Background(), periodDetailID)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("period detail not found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodDetailID := uint(999)

		// Setup expectations
		mockRepo.On("Delete", mock.Anything, periodDetailID).Return(assert.AnError)

		// Execute
		err := mockRepo.Delete(context.Background(), periodDetailID)

		// Assert
		assert.Error(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestPeriodDetailRepository_DeleteByPeriodID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodID := uint(1)

		// Setup expectations
		mockRepo.On("DeleteByPeriodID", mock.Anything, periodID).Return(nil)

		// Execute
		err := mockRepo.DeleteByPeriodID(context.Background(), periodID)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodID := uint(999)

		// Setup expectations
		mockRepo.On("DeleteByPeriodID", mock.Anything, periodID).Return(assert.AnError)

		// Execute
		err := mockRepo.DeleteByPeriodID(context.Background(), periodID)

		// Assert
		assert.Error(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestPeriodDetailRepository_GetUsersByBatch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		lastID := uint(0)
		limit := 10
		expectedUserIDs := []uint{1, 2, 3, 4, 5}

		// Setup expectations
		mockRepo.On("GetUsersByBatch", mock.Anything, lastID, limit).Return(expectedUserIDs, nil)

		// Execute
		userIDs, err := mockRepo.GetUsersByBatch(context.Background(), lastID, limit)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, userIDs, 5)
		assert.Equal(t, expectedUserIDs, userIDs)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("no users found", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		lastID := uint(100)
		limit := 10

		// Setup expectations
		mockRepo.On("GetUsersByBatch", mock.Anything, lastID, limit).Return([]uint{}, nil)

		// Execute
		userIDs, err := mockRepo.GetUsersByBatch(context.Background(), lastID, limit)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, userIDs, 0)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		lastID := uint(0)
		limit := 10

		// Setup expectations
		mockRepo.On("GetUsersByBatch", mock.Anything, lastID, limit).Return(nil, assert.AnError)

		// Execute
		userIDs, err := mockRepo.GetUsersByBatch(context.Background(), lastID, limit)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, userIDs)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestPeriodDetailRepository_CreateBatch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodDetails := []period_detail.PeriodDetail{
			{
				PeriodsID:           1,
				UserID:              1,
				DailyRate:           100000.00,
				TotalWorking:        22,
				AmountSalary:        2200000.00,
				AmountOvertime:      500000.00,
				AmountReimbursement: 200000.00,
				TakeHomePay:         2900000.00,
				CreatedBy:           1,
			},
			{
				PeriodsID:           1,
				UserID:              2,
				DailyRate:           120000.00,
				TotalWorking:        20,
				AmountSalary:        2400000.00,
				AmountOvertime:      300000.00,
				AmountReimbursement: 150000.00,
				TakeHomePay:         2850000.00,
				CreatedBy:           1,
			},
		}

		// Setup expectations
		mockRepo.On("CreateBatch", mock.Anything, periodDetails).Return(nil)

		// Execute
		err := mockRepo.CreateBatch(context.Background(), periodDetails)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty batch", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodDetails := []period_detail.PeriodDetail{}

		// Setup expectations
		mockRepo.On("CreateBatch", mock.Anything, periodDetails).Return(nil)

		// Execute
		err := mockRepo.CreateBatch(context.Background(), periodDetails)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodDetails := []period_detail.PeriodDetail{
			{
				PeriodsID: 1,
				UserID:    1,
			},
		}

		// Setup expectations
		mockRepo.On("CreateBatch", mock.Anything, periodDetails).Return(assert.AnError)

		// Execute
		err := mockRepo.CreateBatch(context.Background(), periodDetails)

		// Assert
		assert.Error(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestPeriodDetailRepository_ListPayslip(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		req := payslip.PayslipListRequest{
			Page:     1,
			Limit:    10,
			SortBy:   "created_at",
			SortDesc: true,
		}
		userID := uint(1)
		expectedResponse := &payslip.PayslipListResponse{
			Data: []payslip.PayslipSummary{
				{
					ID:          1,
					PeriodsID:   1,
					PeriodName:  "January 2024",
					StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
					TakeHomePay: 2900000.00,
					CreatedAt:   time.Now(),
				},
			},
			Pagination: payslip.Pagination{
				Page:       1,
				Limit:      10,
				Total:      1,
				TotalPages: 1,
			},
		}

		// Setup expectations
		mockRepo.On("ListPayslip", mock.Anything, req, userID).Return(expectedResponse, nil)

		// Execute
		response, err := mockRepo.ListPayslip(context.Background(), req, userID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, 1)
		assert.Equal(t, expectedResponse.Pagination.Total, response.Pagination.Total)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		req := payslip.PayslipListRequest{
			Page:  1,
			Limit: 10,
		}
		userID := uint(1)

		// Setup expectations
		mockRepo.On("ListPayslip", mock.Anything, req, userID).Return(nil, assert.AnError)

		// Execute
		response, err := mockRepo.ListPayslip(context.Background(), req, userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, response)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestPeriodDetailRepository_GetPayslipData(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodDetailID := uint(1)
		userID := uint(1)
		expectedData := &payslip.PayslipData{
			EmployeeName:       "John Doe",
			PeriodName:         "January 2024",
			StartDate:          time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:            time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			TotalWorking:       22,
			DailyRate:          100000.00,
			BaseSalary:         2200000.00,
			OvertimeDetails:    []payslip.OvertimeData{},
			TotalOvertime:      500000.00,
			Reimbursements:     []payslip.ReimbursementData{},
			TotalReimbursement: 200000.00,
			TakeHomePay:        2900000.00,
			GeneratedAt:        time.Now(),
		}

		// Setup expectations
		mockRepo.On("GetPayslipData", mock.Anything, periodDetailID, userID).Return(expectedData, nil)

		// Execute
		data, err := mockRepo.GetPayslipData(context.Background(), periodDetailID, userID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Equal(t, expectedData.EmployeeName, data.EmployeeName)
		assert.Equal(t, expectedData.PeriodName, data.PeriodName)
		assert.Equal(t, expectedData.TakeHomePay, data.TakeHomePay)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodDetailID := uint(999)
		userID := uint(1)

		// Setup expectations
		mockRepo.On("GetPayslipData", mock.Anything, periodDetailID, userID).Return(nil, assert.AnError)

		// Execute
		data, err := mockRepo.GetPayslipData(context.Background(), periodDetailID, userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, data)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestPeriodDetailRepository_GetPayslipSummaryData(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodID := uint(1)
		expectedData := &payslip.PayslipSummaryData{
			CompanyName:      "Test Company",
			PeriodName:       "January 2024",
			TotalEmployees:   2,
			TotalTakeHomePay: 5750000.00,
			EmployeeList: []payslip.PayslipSummaryEmployee{
				{
					No:           1,
					EmployeeName: "John Doe",
					TakeHomePay:  2900000.00,
				},
				{
					No:           2,
					EmployeeName: "Jane Smith",
					TakeHomePay:  2850000.00,
				},
			},
		}

		// Setup expectations
		mockRepo.On("GetPayslipSummaryData", mock.Anything, periodID).Return(expectedData, nil)

		// Execute
		data, err := mockRepo.GetPayslipSummaryData(context.Background(), periodID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Equal(t, expectedData.CompanyName, data.CompanyName)
		assert.Equal(t, expectedData.PeriodName, data.PeriodName)
		assert.Equal(t, expectedData.TotalEmployees, data.TotalEmployees)
		assert.Len(t, data.EmployeeList, 2)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test data
		periodID := uint(999)

		// Setup expectations
		mockRepo.On("GetPayslipSummaryData", mock.Anything, periodID).Return(nil, assert.AnError)

		// Execute
		data, err := mockRepo.GetPayslipSummaryData(context.Background(), periodID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, data)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// Test untuk memastikan interface berfungsi dengan benar
func TestPeriodDetailRepository_Interface(t *testing.T) {
	t.Run("interface implementation", func(t *testing.T) {
		// Setup mock
		mockRepo := &mocks.MockIPeriodDetailRepository{}

		// Test bahwa mock mengimplementasikan interface
		var repo IPeriodDetailRepository = mockRepo
		assert.NotNil(t, repo)

		// Test data
		periodDetailData := &period_detail.PeriodDetail{
			PeriodsID:           1,
			UserID:              1,
			DailyRate:           100000.00,
			TotalWorking:        22,
			AmountSalary:        2200000.00,
			AmountOvertime:      500000.00,
			AmountReimbursement: 200000.00,
			TakeHomePay:         2900000.00,
			CreatedBy:           1,
		}

		expectedPeriodDetail := &period_detail.PeriodDetail{
			ID:                  1,
			PeriodsID:           1,
			UserID:              1,
			DailyRate:           100000.00,
			TotalWorking:        22,
			AmountSalary:        2200000.00,
			AmountOvertime:      500000.00,
			AmountReimbursement: 200000.00,
			TakeHomePay:         2900000.00,
			CreatedBy:           1,
			CreatedAt:           time.Now(),
		}

		// Setup expectations for all interface methods
		mockRepo.On("Create", mock.Anything, periodDetailData).Return(nil)
		mockRepo.On("GetByID", mock.Anything, uint(1)).Return(expectedPeriodDetail, nil)
		mockRepo.On("GetByPeriodAndUser", mock.Anything, uint(1), uint(1)).Return(expectedPeriodDetail, nil)
		mockRepo.On("Update", mock.Anything, uint(1), mock.Anything).Return(nil)
		mockRepo.On("Delete", mock.Anything, uint(1)).Return(nil)
		mockRepo.On("DeleteByPeriodID", mock.Anything, uint(1)).Return(nil)
		mockRepo.On("GetUsersByBatch", mock.Anything, uint(0), 10).Return([]uint{1, 2, 3}, nil)
		mockRepo.On("CreateBatch", mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("ListPayslip", mock.Anything, mock.Anything, uint(1)).Return(&payslip.PayslipListResponse{}, nil)
		mockRepo.On("GetPayslipData", mock.Anything, uint(1), uint(1)).Return(&payslip.PayslipData{}, nil)
		mockRepo.On("GetPayslipSummaryData", mock.Anything, uint(1)).Return(&payslip.PayslipSummaryData{}, nil)

		// Test semua method interface
		err := repo.Create(context.Background(), periodDetailData)
		assert.NoError(t, err)

		foundPeriodDetail, err := repo.GetByID(context.Background(), uint(1))
		assert.NoError(t, err)
		assert.Equal(t, expectedPeriodDetail.ID, foundPeriodDetail.ID)

		foundPeriodDetail, err = repo.GetByPeriodAndUser(context.Background(), uint(1), uint(1))
		assert.NoError(t, err)
		assert.Equal(t, expectedPeriodDetail.ID, foundPeriodDetail.ID)

		err = repo.Update(context.Background(), uint(1), map[string]interface{}{"daily_rate": 120000.00})
		assert.NoError(t, err)

		err = repo.Delete(context.Background(), uint(1))
		assert.NoError(t, err)

		err = repo.DeleteByPeriodID(context.Background(), uint(1))
		assert.NoError(t, err)

		userIDs, err := repo.GetUsersByBatch(context.Background(), uint(0), 10)
		assert.NoError(t, err)
		assert.Len(t, userIDs, 3)

		err = repo.CreateBatch(context.Background(), []period_detail.PeriodDetail{})
		assert.NoError(t, err)

		_, err = repo.ListPayslip(context.Background(), payslip.PayslipListRequest{}, uint(1))
		assert.NoError(t, err)

		_, err = repo.GetPayslipData(context.Background(), uint(1), uint(1))
		assert.NoError(t, err)

		_, err = repo.GetPayslipSummaryData(context.Background(), uint(1))
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}
