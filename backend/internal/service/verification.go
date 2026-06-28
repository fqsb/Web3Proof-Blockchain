package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"web3proof/backend/internal/model"

	"gorm.io/gorm"
)

type VerificationService struct {
	db *gorm.DB
}

func NewVerificationService(db *gorm.DB) *VerificationService {
	return &VerificationService{db: db}
}

func (s *VerificationService) VerifyFile(viewerID *uint64, file multipart.File) (map[string]interface{}, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, err
	}
	hash := "0x" + hex.EncodeToString(hasher.Sum(nil))
	return s.verifyByFileHash(viewerID, hash, "file", hash)
}

func (s *VerificationService) VerifyEvidenceNo(viewerID *uint64, evidenceNo string) (map[string]interface{}, error) {
	var evidence model.EvidenceRecord
	err := s.db.Where("evidence_no = ?", evidenceNo).First(&evidence).Error
	if err != nil {
		report := s.report("evidence_no", evidenceNo, false, map[string]interface{}{"reason": "evidence not found"})
		_ = s.saveReport(viewerID, viewerID, "evidence_no", evidenceNo, false, report)
		return report, nil
	}
	return s.buildEvidenceReport(viewerID, "evidence_no", evidenceNo, &evidence)
}

func (s *VerificationService) VerifyCertificateNo(viewerID *uint64, certNo string) (map[string]interface{}, error) {
	var cert model.Certificate
	if err := s.db.Where("certificate_no = ?", certNo).First(&cert).Error; err != nil {
		report := s.report("certificate_no", certNo, false, map[string]interface{}{"reason": "certificate not found"})
		_ = s.saveReport(viewerID, viewerID, "certificate_no", certNo, false, report)
		return report, nil
	}
	var evidence model.EvidenceRecord
	if err := s.db.First(&evidence, cert.EvidenceID).Error; err != nil {
		return nil, err
	}
	report, err := s.buildEvidenceReport(viewerID, "certificate_no", certNo, &evidence)
	if err == nil {
		report["certificate"] = cert
	}
	return report, err
}

func (s *VerificationService) VerifyWallet(viewerID *uint64, address string) (map[string]interface{}, error) {
	address = strings.ToLower(strings.TrimSpace(address))
	var user model.User
	if err := s.db.Where("wallet_address = ?", address).First(&user).Error; err != nil {
		report := s.report("wallet", address, false, map[string]interface{}{"reason": "wallet not found"})
		_ = s.saveReport(viewerID, viewerID, "wallet", address, false, report)
		return report, nil
	}
	var works []model.Work
	s.db.Where("user_id = ? AND visibility = ?", user.ID, "public").Order("created_at desc").Find(&works)
	var evidences []model.EvidenceRecord
	s.db.Where("user_id = ? AND status = ?", user.ID, "confirmed").Order("created_at desc").Find(&evidences)
	var credentials []model.SBTCredential
	s.db.Where("user_id = ? AND status = ?", user.ID, "active").Find(&credentials)
	targetUserID := user.ID
	report := s.report("wallet", address, true, map[string]interface{}{
		"user":        user,
		"works":       works,
		"evidences":   evidences,
		"credentials": credentials,
	})
	if err := s.saveReport(viewerID, &targetUserID, "wallet", address, true, report); err != nil {
		return nil, err
	}
	return report, nil
}

func (s *VerificationService) verifyByFileHash(viewerID *uint64, fileHash, queryType, queryValue string) (map[string]interface{}, error) {
	var evidence model.EvidenceRecord
	err := s.db.Where("file_hash = ? AND status = ?", fileHash, "confirmed").First(&evidence).Error
	if err != nil {
		report := s.report(queryType, queryValue, false, map[string]interface{}{
			"file_hash": fileHash,
			"reason":    "matching confirmed evidence not found",
		})
		_ = s.saveReport(viewerID, viewerID, queryType, queryValue, false, report)
		return report, nil
	}
	return s.buildEvidenceReport(viewerID, queryType, queryValue, &evidence)
}

func (s *VerificationService) buildEvidenceReport(viewerID *uint64, queryType, queryValue string, evidence *model.EvidenceRecord) (map[string]interface{}, error) {
	var work model.Work
	_ = s.db.First(&work, evidence.WorkID).Error
	var cert model.Certificate
	_ = s.db.Where("evidence_id = ?", evidence.ID).First(&cert).Error
	var credential model.SBTCredential
	_ = s.db.Where("evidence_id = ? AND status = ?", evidence.ID, "active").First(&credential).Error
	report := s.report(queryType, queryValue, true, map[string]interface{}{
		"evidence":    evidence,
		"work":        work,
		"certificate": cert,
		"credential":  credential,
		"conclusion":  "该材料哈希与平台链上存证记录一致。",
	})
	targetUserID := evidence.UserID
	if err := s.saveReport(viewerID, &targetUserID, queryType, queryValue, true, report); err != nil {
		return nil, err
	}
	return report, nil
}

func (s *VerificationService) report(queryType, queryValue string, passed bool, data map[string]interface{}) map[string]interface{} {
	data["query_type"] = queryType
	data["query_value"] = queryValue
	data["passed"] = passed
	data["verified_at"] = time.Now().UTC()
	return data
}

func (s *VerificationService) saveReport(viewerID, targetUserID *uint64, queryType, queryValue string, passed bool, report map[string]interface{}) error {
	if targetUserID == nil {
		return nil
	}
	raw, err := json.Marshal(report)
	if err != nil {
		return err
	}
	if queryValue == "" {
		queryValue = fmt.Sprintf("%s-%d", queryType, time.Now().UnixNano())
	}
	return s.db.Create(&model.VerificationReport{
		ViewerID:     viewerID,
		TargetUserID: targetUserID,
		QueryType:    queryType,
		QueryValue:   queryValue,
		Passed:       passed,
		ReportJSON:   string(raw),
	}).Error
}

func (s *VerificationService) ListReports(viewerID uint64) ([]model.VerificationReport, error) {
	var reports []model.VerificationReport
	err := s.db.Where("viewer_id = ?", viewerID).Order("created_at desc").Limit(50).Find(&reports).Error
	return reports, err
}
