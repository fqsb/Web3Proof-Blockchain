package service

import (
	"strings"
	"time"

	"web3proof/backend/internal/model"

	"gorm.io/gorm"
)

type ReputationService struct {
	db *gorm.DB
}

func NewReputationService(db *gorm.DB) *ReputationService {
	return &ReputationService{db: db}
}

func grade(total uint) string {
	switch {
	case total >= 800:
		return "A"
	case total >= 600:
		return "B"
	case total >= 400:
		return "C"
	default:
		return "D"
	}
}

func (s *ReputationService) Recalculate(userID uint64) (*model.ReputationScore, error) {
	return s.recalculate(s.db, userID)
}

func (s *ReputationService) RecalculateTx(tx *gorm.DB, userID uint64) (*model.ReputationScore, error) {
	return s.recalculate(tx, userID)
}

func (s *ReputationService) recalculate(db *gorm.DB, userID uint64) (*model.ReputationScore, error) {
	var projectCount int64
	db.Model(&model.EvidenceRecord{}).Where("user_id = ? AND status = ?", userID, "confirmed").Count(&projectCount)
	projectScore := uint(projectCount) * 100
	if projectScore > 500 {
		projectScore = 500
	}

	var certCount int64
	db.Model(&model.SBTCredential{}).Where("user_id = ? AND status = ?", userID, "active").Count(&certCount)
	certScore := uint(certCount) * 100
	if certScore > 300 {
		certScore = 300
	}

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	activityScore := uint(0)
	if user.LastActiveAt != nil && user.LastActiveAt.After(now.AddDate(0, 0, -7)) {
		activityScore += 30
	}
	var recentProjectCount int64
	db.Model(&model.Work{}).
		Where("user_id = ? AND (created_at >= ? OR updated_at >= ?)", userID, now.AddDate(0, 0, -30), now.AddDate(0, 0, -30)).
		Count(&recentProjectCount)
	if recentProjectCount > 0 {
		activityScore += 50
	}
	var recentApplicationCount int64
	db.Model(&model.CertificationApplication{}).
		Where("user_id = ? AND created_at >= ?", userID, now.AddDate(0, 0, -30)).
		Count(&recentApplicationCount)
	if recentApplicationCount > 0 {
		activityScore += 40
	}
	if nonEmpty(user.Nickname) && nonEmpty(user.Bio) && nonEmpty(user.Email) {
		activityScore += 30
	}
	if user.IsDIDRegistered {
		activityScore += 50
	}
	if activityScore > 200 {
		activityScore = 200
	}

	total := projectScore + certScore + activityScore
	score := model.ReputationScore{
		UserID:        userID,
		ProjectScore:  projectScore,
		CertScore:     certScore,
		ActivityScore: activityScore,
		TotalScore:    total,
		Grade:         grade(total),
	}

	var existing model.ReputationScore
	err := db.Where("user_id = ?", userID).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		if err := db.Create(&score).Error; err != nil {
			return nil, err
		}
		return &score, nil
	}
	if err != nil {
		return nil, err
	}
	existing.ProjectScore = score.ProjectScore
	existing.CertScore = score.CertScore
	existing.ActivityScore = score.ActivityScore
	existing.TotalScore = score.TotalScore
	existing.Grade = score.Grade
	if err := db.Save(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

func nonEmpty(value *string) bool {
	return value != nil && strings.TrimSpace(*value) != ""
}

func (s *ReputationService) GetByUserID(userID uint64) (*model.ReputationScore, error) {
	var score model.ReputationScore
	err := s.db.Where("user_id = ?", userID).First(&score).Error
	if err == gorm.ErrRecordNotFound {
		recalculated, recalcErr := s.Recalculate(userID)
		return recalculated, recalcErr
	}
	return &score, err
}
