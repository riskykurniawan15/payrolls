package audit_trail

import (
	"context"
	"encoding/json"

	"github.com/riskykurniawan15/payrolls/models/audit_trail"
	auditTrailRepo "github.com/riskykurniawan15/payrolls/repositories/audit_trail"
)

type (
	IAuditTrailService interface {
		LogRequest(ctx context.Context, req *audit_trail.CreateAuditTrailRequest) error
	}
	AuditTrailService struct {
		auditTrailRepo auditTrailRepo.IAuditTrailRepository
	}
)

func NewAuditTrailService(auditTrailRepo auditTrailRepo.IAuditTrailRepository) IAuditTrailService {
	return &AuditTrailService{
		auditTrailRepo: auditTrailRepo,
	}
}

func (s *AuditTrailService) LogRequest(ctx context.Context, req *audit_trail.CreateAuditTrailRequest) error {
	// Convert request to model
	auditTrailModel := &audit_trail.AuditTrail{
		IP:             req.IP,
		Method:         req.Method,
		Path:           req.Path,
		UserID:         req.UserID,
		Payload:        req.Payload,
		ResponseCode:   req.ResponseCode,
		ErrorResponse:  req.ErrorResponse,
		ResponseTimeMs: req.ResponseTimeMs,
		UserAgent:      req.UserAgent,
	}

	// Save to database
	return s.auditTrailRepo.Create(ctx, auditTrailModel)
}

// Helper function untuk serialize payload
func SerializePayload(payload interface{}) (*string, error) {
	if payload == nil {
		return nil, nil
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	jsonStr := string(jsonBytes)
	return &jsonStr, nil
}
