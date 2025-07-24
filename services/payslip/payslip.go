package payslip

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/riskykurniawan15/payrolls/config"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
	"github.com/riskykurniawan15/payrolls/models/payslip"
	periodDetailRepo "github.com/riskykurniawan15/payrolls/repositories/period_detail"
	"github.com/riskykurniawan15/payrolls/utils/logger"
)

type (
	IPayslipService interface {
		List(ctx context.Context, req payslip.PayslipListRequest, userID uint) (*payslip.PayslipListResponse, error)
		GeneratePayslip(ctx context.Context, periodDetailID, userID uint, host string) (*payslip.GeneratePayslipResponse, error)
		GetPayslipData(ctx context.Context, token string) (*payslip.PayslipData, error)
		GeneratePayslipSummary(ctx context.Context, periodID uint, host string) (*payslip.GeneratePayslipResponse, error)
		GetPayslipSummaryData(ctx context.Context, token string) (*payslip.PayslipSummaryData, error)
	}

	PayslipService struct {
		logger           logger.Logger
		periodDetailRepo periodDetailRepo.IPeriodDetailRepository
		config           config.Config
	}

	// TokenData for encrypted token
	TokenData struct {
		PeriodDetailID uint      `json:"period_detail_id"`
		UserID         uint      `json:"user_id"`
		ExpiresAt      time.Time `json:"expires_at"`
	}

	// SummaryTokenData for encrypted summary token
	SummaryTokenData struct {
		PeriodID  uint      `json:"period_id"`
		ExpiresAt time.Time `json:"expires_at"`
	}
)

func NewPayslipService(
	logger logger.Logger,
	periodDetailRepo periodDetailRepo.IPeriodDetailRepository,
	config config.Config,
) IPayslipService {
	return &PayslipService{
		logger:           logger,
		periodDetailRepo: periodDetailRepo,
		config:           config,
	}
}

func (s *PayslipService) List(ctx context.Context, req payslip.PayslipListRequest, userID uint) (*payslip.PayslipListResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing payslip list request", requestID, map[string]interface{}{
		"user_id": userID,
		"page":    req.Page,
		"limit":   req.Limit,
	})

	response, err := s.periodDetailRepo.ListPayslip(ctx, req, userID)
	if err != nil {
		s.logger.ErrorT("failed to get payslip list", requestID, map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return nil, fmt.Errorf("failed to get payslip list: %w", err)
	}

	s.logger.InfoT("payslip list retrieved successfully", requestID, map[string]interface{}{
		"user_id": userID,
		"total":   response.Pagination.Total,
	})

	return response, nil
}

func (s *PayslipService) GeneratePayslip(ctx context.Context, periodDetailID, userID uint, host string) (*payslip.GeneratePayslipResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing generate payslip request", requestID, map[string]interface{}{
		"period_detail_id": periodDetailID,
		"user_id":          userID,
	})

	// Verify that the period detail belongs to the user
	_, err := s.periodDetailRepo.GetPayslipData(ctx, periodDetailID, userID)
	if err != nil {
		s.logger.ErrorT("failed to get payslip data", requestID, map[string]interface{}{
			"error":            err.Error(),
			"period_detail_id": periodDetailID,
			"user_id":          userID,
		})
		return nil, fmt.Errorf("failed to get payslip data: %w", err)
	}

	// Generate token
	token, expiresAt, err := s.generateToken(periodDetailID, userID)
	if err != nil {
		s.logger.ErrorT("failed to generate token", requestID, map[string]interface{}{
			"error":            err.Error(),
			"period_detail_id": periodDetailID,
		})
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Generate print URL
	printURL := fmt.Sprintf("%s/payslip/print?token=%s", host, token)

	s.logger.InfoT("payslip generated successfully", requestID, map[string]interface{}{
		"period_detail_id": periodDetailID,
		"user_id":          userID,
		"expires_at":       expiresAt,
	})

	return &payslip.GeneratePayslipResponse{
		PrintURL:  printURL,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *PayslipService) GetPayslipData(ctx context.Context, token string) (*payslip.PayslipData, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing get payslip data request", requestID, map[string]interface{}{
		"token": token,
	})

	// Decrypt and validate token
	tokenData, err := s.decryptToken(token)
	if err != nil {
		s.logger.ErrorT("failed to decrypt token", requestID, map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Check if token is expired
	if time.Now().After(tokenData.ExpiresAt) {
		s.logger.WarningT("token expired", requestID, map[string]interface{}{
			"expires_at": tokenData.ExpiresAt,
		})
		return nil, fmt.Errorf("token expired")
	}

	// Get payslip data
	payslipData, err := s.periodDetailRepo.GetPayslipData(ctx, tokenData.PeriodDetailID, tokenData.UserID)
	if err != nil {
		s.logger.ErrorT("failed to get payslip data", requestID, map[string]interface{}{
			"error":            err.Error(),
			"period_detail_id": tokenData.PeriodDetailID,
			"user_id":          tokenData.UserID,
		})
		return nil, fmt.Errorf("failed to get payslip data: %w", err)
	}

	s.logger.InfoT("payslip data retrieved successfully", requestID, map[string]interface{}{
		"period_detail_id": tokenData.PeriodDetailID,
		"user_id":          tokenData.UserID,
	})

	payslipData.CompanyName = s.config.CompanyName

	return payslipData, nil
}

func (s *PayslipService) GeneratePayslipSummary(ctx context.Context, periodID uint, host string) (*payslip.GeneratePayslipResponse, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing generate payslip summary request", requestID, map[string]interface{}{
		"period_id": periodID,
	})

	// Verify that the period exists and get summary data
	_, err := s.periodDetailRepo.GetPayslipSummaryData(ctx, periodID)
	if err != nil {
		s.logger.ErrorT("failed to get payslip summary data", requestID, map[string]interface{}{
			"error":     err.Error(),
			"period_id": periodID,
		})
		return nil, fmt.Errorf("failed to get payslip summary data: %w", err)
	}

	// Generate token
	token, expiresAt, err := s.generateSummaryToken(periodID)
	if err != nil {
		s.logger.ErrorT("failed to generate summary token", requestID, map[string]interface{}{
			"error":     err.Error(),
			"period_id": periodID,
		})
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Generate print URL
	printURL := fmt.Sprintf("%s/payslip/summary/print?token=%s", host, token)

	s.logger.InfoT("payslip summary generated successfully", requestID, map[string]interface{}{
		"period_id":  periodID,
		"expires_at": expiresAt,
	})

	return &payslip.GeneratePayslipResponse{
		PrintURL:  printURL,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *PayslipService) GetPayslipSummaryData(ctx context.Context, token string) (*payslip.PayslipSummaryData, error) {
	requestID := middleware.GetRequestIDFromContext(ctx)
	s.logger.InfoT("processing get payslip summary data request", requestID, map[string]interface{}{
		"token": token,
	})

	// Decrypt and validate token
	tokenData, err := s.decryptSummaryToken(token)
	if err != nil {
		s.logger.ErrorT("failed to decrypt summary token", requestID, map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Check if token is expired
	if time.Now().After(tokenData.ExpiresAt) {
		s.logger.WarningT("summary token expired", requestID, map[string]interface{}{
			"expires_at": tokenData.ExpiresAt,
		})
		return nil, fmt.Errorf("token expired")
	}

	// Get payslip summary data
	payslipSummaryData, err := s.periodDetailRepo.GetPayslipSummaryData(ctx, tokenData.PeriodID)
	if err != nil {
		s.logger.ErrorT("failed to get payslip summary data", requestID, map[string]interface{}{
			"error":     err.Error(),
			"period_id": tokenData.PeriodID,
		})
		return nil, fmt.Errorf("failed to get payslip summary data: %w", err)
	}

	s.logger.InfoT("payslip summary data retrieved successfully", requestID, map[string]interface{}{
		"period_id": tokenData.PeriodID,
	})

	payslipSummaryData.CompanyName = s.config.CompanyName

	return payslipSummaryData, nil
}

func (s *PayslipService) generateToken(periodDetailID, userID uint) (string, time.Time, error) {
	// Set expiration time (5 minutes from now)
	expiresAt := time.Now().Add(5 * time.Minute)

	// Create token data
	tokenData := TokenData{
		PeriodDetailID: periodDetailID,
		UserID:         userID,
		ExpiresAt:      expiresAt,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(tokenData)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to marshal token data: %w", err)
	}

	// Encrypt the data
	encryptedData, err := s.encrypt(jsonData)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to encrypt token: %w", err)
	}

	// Encode to base64
	token := base64.URLEncoding.EncodeToString(encryptedData)

	return token, expiresAt, nil
}

func (s *PayslipService) decryptToken(token string) (*TokenData, error) {
	// Decode from base64
	encryptedData, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token format: %w", err)
	}

	// Decrypt the data
	decryptedData, err := s.decrypt(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token: %w", err)
	}

	// Parse JSON
	var tokenData TokenData
	if err := json.Unmarshal(decryptedData, &tokenData); err != nil {
		return nil, fmt.Errorf("invalid token data: %w", err)
	}

	return &tokenData, nil
}

func (s *PayslipService) generateSummaryToken(periodID uint) (string, time.Time, error) {
	// Set expiration time (5 minutes from now)
	expiresAt := time.Now().Add(5 * time.Minute)

	// Create token data
	tokenData := SummaryTokenData{
		PeriodID:  periodID,
		ExpiresAt: expiresAt,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(tokenData)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to marshal token data: %w", err)
	}

	// Encrypt the data
	encryptedData, err := s.encrypt(jsonData)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to encrypt token: %w", err)
	}

	// Encode to base64
	token := base64.URLEncoding.EncodeToString(encryptedData)

	return token, expiresAt, nil
}

func (s *PayslipService) decryptSummaryToken(token string) (*SummaryTokenData, error) {
	// Decode from base64
	encryptedData, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token format: %w", err)
	}

	// Decrypt the data
	decryptedData, err := s.decrypt(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token: %w", err)
	}

	// Parse JSON
	var tokenData SummaryTokenData
	if err := json.Unmarshal(decryptedData, &tokenData); err != nil {
		return nil, fmt.Errorf("invalid token data: %w", err)
	}

	return &tokenData, nil
}

func (s *PayslipService) encrypt(data []byte) ([]byte, error) {
	// Use JWT secret as encryption key (pad to 32 bytes for AES-256)
	key := make([]byte, 32)
	copy(key, []byte(s.config.JWT.SecretKey))

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (s *PayslipService) decrypt(data []byte) ([]byte, error) {
	// Use JWT secret as encryption key (pad to 32 bytes for AES-256)
	key := make([]byte, 32)
	copy(key, []byte(s.config.JWT.SecretKey))

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
