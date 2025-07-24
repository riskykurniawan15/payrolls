package reimbursement

import (
	"context"
	"fmt"
	"time"

	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
	"github.com/riskykurniawan15/payrolls/models/reimbursement"
	reimbursementRepo "github.com/riskykurniawan15/payrolls/repositories/reimbursement"
	"github.com/riskykurniawan15/payrolls/utils/logger"
)

type (
	IReimbursementService interface {
		Create(ctx context.Context, req reimbursement.CreateReimbursementRequest, userID uint) (*reimbursement.ReimbursementResponse, error)
		GetByID(ctx context.Context, id uint, userID uint) (*reimbursement.ReimbursementResponse, error)
		Update(ctx context.Context, id uint, req reimbursement.UpdateReimbursementRequest, userID uint) (*reimbursement.ReimbursementResponse, error)
		Delete(ctx context.Context, id uint, userID uint) error
		List(ctx context.Context, req reimbursement.ListReimbursementsRequest, userID uint) (*reimbursement.ListReimbursementsResponse, error)
	}

	ReimbursementService struct {
		logger            logger.Logger
		reimbursementRepo reimbursementRepo.IReimbursementRepository
	}
)

func NewReimbursementService(logger logger.Logger, reimbursementRepo reimbursementRepo.IReimbursementRepository) IReimbursementService {
	return &ReimbursementService{
		logger:            logger,
		reimbursementRepo: reimbursementRepo,
	}
}

func (s *ReimbursementService) Create(ctx context.Context, req reimbursement.CreateReimbursementRequest, userID uint) (*reimbursement.ReimbursementResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing create reimbursement request", requestID, map[string]interface{}{
		"user_id":    userID,
		"title":      req.Title,
		"date":       req.Date,
		"amount":     req.Amount,
		"created_by": userID,
	})

	// Parse date or use today
	var reimbursementDate time.Time
	if req.Date != nil {
		reimbursementDate = req.Date.Time
	} else {
		// Use today's date at 00:00:00
		now := time.Now()
		reimbursementDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}

	// Check if date is in the future
	if reimbursementDate.After(time.Now()) {
		s.logger.WarningT("reimbursement date is in the future", requestID, map[string]interface{}{
			"date": reimbursementDate.Format("2006-01-02"),
		})
		return nil, fmt.Errorf("reimbursement date cannot be in the future")
	}

	s.logger.InfoT("reimbursement validation passed", requestID, map[string]interface{}{
		"user_id": userID,
		"title":   req.Title,
		"date":    reimbursementDate.Format("2006-01-02"),
		"amount":  req.Amount,
	})

	// Create reimbursement
	reimbursementRecord := &reimbursement.Reimbursement{
		UserID:      userID,
		Title:       req.Title,
		Date:        reimbursementDate,
		Amount:      req.Amount,
		Description: req.Description,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedBy:   nil,
		UpdatedAt:   nil,
	}

	s.logger.InfoT("creating reimbursement in database", requestID, map[string]interface{}{
		"user_id":    userID,
		"title":      req.Title,
		"date":       reimbursementDate.Format("2006-01-02"),
		"amount":     req.Amount,
		"created_by": userID,
	})

	if err := s.reimbursementRepo.Create(ctx, reimbursementRecord); err != nil {
		s.logger.ErrorT("failed to create reimbursement", requestID, map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
			"title":   req.Title,
		})
		return nil, fmt.Errorf("failed to create reimbursement: %w", err)
	}

	s.logger.InfoT("reimbursement created successfully", requestID, map[string]interface{}{
		"reimbursement_id": reimbursementRecord.ID,
		"user_id":          userID,
		"title":            req.Title,
		"date":             reimbursementDate.Format("2006-01-02"),
		"amount":           req.Amount,
	})

	// Convert to response
	response := s.toResponse(*reimbursementRecord)
	return &response, nil
}

func (s *ReimbursementService) GetByID(ctx context.Context, id uint, userID uint) (*reimbursement.ReimbursementResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing get reimbursement by ID request", requestID, map[string]interface{}{
		"reimbursement_id": id,
		"user_id":          userID,
	})

	reimb, err := s.reimbursementRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorT("failed to get reimbursement by ID", requestID, map[string]interface{}{
			"error":            err.Error(),
			"reimbursement_id": id,
		})
		return nil, fmt.Errorf("reimbursement not found: %w", err)
	}

	// Check if reimbursement belongs to the user
	if reimb.UserID != userID {
		s.logger.WarningT("user trying to access reimbursement of another user", requestID, map[string]interface{}{
			"reimbursement_id":      id,
			"reimbursement_user_id": reimb.UserID,
			"request_user_id":       userID,
		})
		return nil, fmt.Errorf("reimbursement not found")
	}

	s.logger.InfoT("reimbursement retrieved successfully", requestID, map[string]interface{}{
		"reimbursement_id": id,
		"user_id":          reimb.UserID,
		"title":            reimb.Title,
		"date":             reimb.Date.Format("2006-01-02"),
		"amount":           reimb.Amount,
	})

	response := s.toResponse(*reimb)
	return &response, nil
}

func (s *ReimbursementService) Update(ctx context.Context, id uint, req reimbursement.UpdateReimbursementRequest, userID uint) (*reimbursement.ReimbursementResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing update reimbursement request", requestID, map[string]interface{}{
		"reimbursement_id": id,
		"user_id":          userID,
	})

	// Check if reimbursement exists
	existingReimbursement, err := s.reimbursementRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorT("failed to get reimbursement for update", requestID, map[string]interface{}{
			"error":            err.Error(),
			"reimbursement_id": id,
		})
		return nil, fmt.Errorf("reimbursement not found: %w", err)
	}

	// Check if reimbursement belongs to the user
	if existingReimbursement.UserID != userID {
		s.logger.WarningT("user trying to update reimbursement of another user", requestID, map[string]interface{}{
			"reimbursement_id":      id,
			"reimbursement_user_id": existingReimbursement.UserID,
			"request_user_id":       userID,
		})
		return nil, fmt.Errorf("reimbursement not found")
	}

	s.logger.InfoT("reimbursement validation passed", requestID, map[string]interface{}{
		"reimbursement_id": id,
		"user_id":          existingReimbursement.UserID,
		"title":            existingReimbursement.Title,
		"date":             existingReimbursement.Date.Format("2006-01-02"),
		"amount":           existingReimbursement.Amount,
	})

	// Prepare updates
	updates := make(map[string]interface{})
	updates["updated_by"] = userID
	updates["updated_at"] = time.Now()

	// Update title if provided
	if req.Title != nil {
		updates["title"] = *req.Title
	}

	// Update date if provided
	if req.Date != nil {
		date := req.Date.Time

		// Check if date is in the future
		if date.After(time.Now()) {
			return nil, fmt.Errorf("reimbursement date cannot be in the future")
		}

		updates["date"] = date
	}

	// Update total if provided
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}

	// Update description if provided
	if req.Description != nil {
		updates["description"] = *req.Description
	}

	// Apply updates
	if err := s.reimbursementRepo.Update(ctx, id, updates); err != nil {
		return nil, fmt.Errorf("failed to update reimbursement: %w", err)
	}

	// Get updated reimbursement
	updatedReimbursement, err := s.reimbursementRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated reimbursement: %w", err)
	}

	response := s.toResponse(*updatedReimbursement)
	return &response, nil
}

func (s *ReimbursementService) Delete(ctx context.Context, id uint, userID uint) error {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing delete reimbursement request", requestID, map[string]interface{}{
		"reimbursement_id": id,
		"user_id":          userID,
	})

	// Check if reimbursement exists
	existingReimbursement, err := s.reimbursementRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorT("failed to get reimbursement for delete", requestID, map[string]interface{}{
			"error":            err.Error(),
			"reimbursement_id": id,
		})
		return fmt.Errorf("reimbursement not found: %w", err)
	}

	// Check if reimbursement belongs to the user
	if existingReimbursement.UserID != userID {
		s.logger.WarningT("user trying to delete reimbursement of another user", requestID, map[string]interface{}{
			"reimbursement_id":      id,
			"reimbursement_user_id": existingReimbursement.UserID,
			"request_user_id":       userID,
		})
		return fmt.Errorf("reimbursement not found")
	}

	s.logger.InfoT("reimbursement validation passed for delete", requestID, map[string]interface{}{
		"reimbursement_id": id,
		"user_id":          existingReimbursement.UserID,
		"title":            existingReimbursement.Title,
		"date":             existingReimbursement.Date.Format("2006-01-02"),
		"amount":           existingReimbursement.Amount,
	})

	// Delete reimbursement
	if err := s.reimbursementRepo.Delete(ctx, id); err != nil {
		s.logger.ErrorT("failed to delete reimbursement", requestID, map[string]interface{}{
			"error":            err.Error(),
			"reimbursement_id": id,
		})
		return fmt.Errorf("failed to delete reimbursement: %w", err)
	}

	s.logger.InfoT("reimbursement deleted successfully", requestID, map[string]interface{}{
		"reimbursement_id": id,
		"user_id":          existingReimbursement.UserID,
		"title":            existingReimbursement.Title,
		"date":             existingReimbursement.Date.Format("2006-01-02"),
		"amount":           existingReimbursement.Amount,
	})

	return nil
}

func (s *ReimbursementService) List(ctx context.Context, req reimbursement.ListReimbursementsRequest, userID uint) (*reimbursement.ListReimbursementsResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing list reimbursements request", requestID, map[string]interface{}{
		"page":       req.Page,
		"limit":      req.Limit,
		"user_id":    userID,
		"start_date": req.StartDate,
		"end_date":   req.EndDate,
		"sort_by":    req.SortBy,
		"sort_desc":  req.SortDesc,
	})

	// Set default values if not provided
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	s.logger.InfoT("list reimbursements with filters", requestID, map[string]interface{}{
		"page":       req.Page,
		"limit":      req.Limit,
		"user_id":    userID,
		"start_date": req.StartDate,
		"end_date":   req.EndDate,
		"sort_by":    req.SortBy,
		"sort_desc":  req.SortDesc,
	})

	response, err := s.reimbursementRepo.List(ctx, req, userID)
	if err != nil {
		s.logger.ErrorT("failed to list reimbursements", requestID, map[string]interface{}{
			"error": err.Error(),
			"page":  req.Page,
			"limit": req.Limit,
		})
		return nil, fmt.Errorf("failed to list reimbursements: %w", err)
	}

	s.logger.InfoT("reimbursements listed successfully", requestID, map[string]interface{}{
		"total_count":  response.Pagination.Total,
		"total_pages":  response.Pagination.TotalPages,
		"current_page": response.Pagination.Page,
		"limit":        response.Pagination.Limit,
		"data_count":   len(response.Data),
	})

	return response, nil
}

// Helper function to convert Reimbursement to ReimbursementResponse
func (s *ReimbursementService) toResponse(reimb reimbursement.Reimbursement) reimbursement.ReimbursementResponse {
	return reimbursement.ReimbursementResponse{
		ID:          reimb.ID,
		UserID:      reimb.UserID,
		Title:       reimb.Title,
		Date:        reimb.Date,
		Amount:      reimb.Amount,
		Description: reimb.Description,
		CreatedBy:   reimb.CreatedBy,
		CreatedAt:   reimb.CreatedAt,
		UpdatedBy:   reimb.UpdatedBy,
		UpdatedAt:   reimb.UpdatedAt,
	}
}
