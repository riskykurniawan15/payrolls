package attendance

import (
	"context"
	"errors"
	"time"

	"github.com/riskykurniawan15/payrolls/models/attendance"
	"gorm.io/gorm"
)

type (
	IAttendanceRepository interface {
		GetAttendances(ctx context.Context, userID uint, page, limit int, startDate, endDate *time.Time) ([]attendance.Attendance, int64, error)
		GetAttendanceByID(ctx context.Context, id, userID uint) (attendance.Attendance, error)
		GetLatestCheckInByUserID(ctx context.Context, userID uint, date *time.Time) (attendance.Attendance, error)
		GetByUserAndDate(ctx context.Context, userID uint, date time.Time) (*attendance.Attendance, error)
		CreateAttendance(ctx context.Context, attendance attendance.Attendance) (attendance.Attendance, error)
		UpdateAttendance(ctx context.Context, attendance attendance.Attendance) (attendance.Attendance, error)
		GetAttendanceByIDForUpdate(ctx context.Context, id, userID uint) (attendance.Attendance, error)
	}

	AttendanceRepository struct {
		db *gorm.DB
	}
)

func NewAttendanceRepository(db *gorm.DB) IAttendanceRepository {
	return &AttendanceRepository{
		db: db,
	}
}

func (repo *AttendanceRepository) GetAttendances(ctx context.Context, userID uint, page, limit int, startDate, endDate *time.Time) ([]attendance.Attendance, int64, error) {
	var attendances []attendance.Attendance
	var total int64

	offset := (page - 1) * limit

	// Build query
	query := repo.db.WithContext(ctx).Model(&attendance.Attendance{}).Where("user_id = ?", userID)

	// Add date filter if provided
	if startDate != nil && endDate != nil {
		// Filter by check-in date range using date comparison to avoid timezone issues
		query = query.Where("DATE(check_in_date) >= DATE(?) AND DATE(check_in_date) <= DATE(?)", startDate, endDate)
	} else if startDate != nil {
		// Only start date provided
		query = query.Where("DATE(check_in_date) >= DATE(?)", startDate)
	} else if endDate != nil {
		// Only end date provided
		query = query.Where("DATE(check_in_date) <= DATE(?)", endDate)
	}
	// If both are nil, no date filter applied (get all records)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get attendances with pagination
	if err := query.
		Order("check_in_date DESC").
		Offset(offset).
		Limit(limit).
		Find(&attendances).Error; err != nil {
		return nil, 0, err
	}

	return attendances, total, nil
}

func (repo *AttendanceRepository) GetAttendanceByID(ctx context.Context, id, userID uint) (attendance.Attendance, error) {
	var attendance attendance.Attendance

	if err := repo.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&attendance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return attendance, errors.New("attendance not found")
		}
		return attendance, err
	}

	return attendance, nil
}

func (repo *AttendanceRepository) GetByUserAndDate(ctx context.Context, userID uint, date time.Time) (*attendance.Attendance, error) {
	var attendance attendance.Attendance

	// Get attendance for the specific date (from start of day to end of day)
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := repo.db.WithContext(ctx).
		Where("user_id = ? AND check_in_date >= ? AND check_in_date < ?", userID, startOfDay, endOfDay).
		First(&attendance).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}

	return &attendance, nil
}

func (repo *AttendanceRepository) GetLatestCheckInByUserID(ctx context.Context, userID uint, date *time.Time) (attendance.Attendance, error) {
	var attendance attendance.Attendance

	// Build query
	query := repo.db.WithContext(ctx).Where("user_id = ?", userID)

	// Add date filter if provided
	if date != nil {
		// Filter by check-in date (same day)
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		query = query.Where("check_in_date >= ? AND check_in_date < ?", startOfDay, endOfDay)
	}

	if err := query.Order("check_in_date DESC").First(&attendance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return attendance, errors.New("no active check-in found")
		}
		return attendance, err
	}

	return attendance, nil
}

func (repo *AttendanceRepository) CreateAttendance(ctx context.Context, attendance attendance.Attendance) (attendance.Attendance, error) {
	if err := repo.db.WithContext(ctx).Create(&attendance).Error; err != nil {
		return attendance, err
	}

	return attendance, nil
}

func (repo *AttendanceRepository) UpdateAttendance(ctx context.Context, attendance attendance.Attendance) (attendance.Attendance, error) {
	if err := repo.db.WithContext(ctx).Save(&attendance).Error; err != nil {
		return attendance, err
	}

	return attendance, nil
}

func (repo *AttendanceRepository) GetAttendanceByIDForUpdate(ctx context.Context, id, userID uint) (attendance.Attendance, error) {
	var attendance attendance.Attendance

	if err := repo.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&attendance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return attendance, errors.New("attendance not found")
		}
		return attendance, err
	}

	return attendance, nil
}
