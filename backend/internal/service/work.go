package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"web3proof/backend/internal/config"
	"web3proof/backend/internal/model"
	"web3proof/backend/internal/pkg/eth"
	"web3proof/backend/internal/pkg/storage"

	"gorm.io/gorm"
)

type WorkService struct {
	db    *gorm.DB
	store *storage.LocalStore
	cfg   *config.Config
	eth   *eth.EthClient
	rep   *ReputationService
}

func NewWorkService(db *gorm.DB, store *storage.LocalStore, cfg *config.Config, ethClient *eth.EthClient, rep *ReputationService) *WorkService {
	return &WorkService{db: db, store: store, cfg: cfg, eth: ethClient, rep: rep}
}

type CreateWorkInput struct {
	Title       string
	Description string
	CategoryID  *uint
	ExternalURL string
	Visibility  string
}

type EvidencePrepareResult struct {
	ContractAddress string `json:"contract_address"`
	EvidenceNo      string `json:"evidence_no"`
	EvidenceNoHash  string `json:"evidence_no_hash"`
	FileHash        string `json:"file_hash"`
	MetadataURI     string `json:"metadata_uri"`
	ChainID         int64  `json:"chain_id"`
}

func (s *WorkService) Create(userID uint64, in CreateWorkInput) (*model.Work, error) {
	desc := strings.TrimSpace(in.Description)
	external := strings.TrimSpace(in.ExternalURL)
	visibility := in.Visibility
	if visibility != "public" {
		visibility = "private"
	}
	work := model.Work{
		UserID:     userID,
		Title:      strings.TrimSpace(in.Title),
		CategoryID: in.CategoryID,
		Visibility: visibility,
		Status:     "draft",
	}
	if desc != "" {
		work.Description = &desc
	}
	if external != "" {
		work.ExternalURL = &external
	}
	if err := s.db.Create(&work).Error; err != nil {
		return nil, err
	}
	return &work, nil
}

func (s *WorkService) ListByUser(userID uint64) ([]model.Work, error) {
	var works []model.Work
	err := s.db.Preload("Category").Where("user_id = ?", userID).Order("created_at desc").Find(&works).Error
	return works, err
}

func (s *WorkService) GetByID(userID, workID uint64) (*model.Work, error) {
	var work model.Work
	err := s.db.Preload("Category").Where("id = ? AND user_id = ?", workID, userID).First(&work).Error
	return &work, err
}

func (s *WorkService) UploadFile(userID, workID uint64, file multipart.File, header *multipart.FileHeader) (*model.WorkFile, error) {
	if _, err := s.GetByID(userID, workID); err != nil {
		return nil, err
	}
	saved, err := s.store.SaveUploaded(fmt.Sprintf("works/%d", workID), file, header)
	if err != nil {
		return nil, err
	}
	url := saved.URL
	record := model.WorkFile{
		WorkID:     workID,
		UserID:     userID,
		FileName:   header.Filename,
		FileType:   header.Header.Get("Content-Type"),
		StorageKey: saved.Key,
		StorageURL: &url,
		FileSize:   saved.Size,
		SHA256Hash: saved.SHA256,
	}
	if err := s.db.Create(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *WorkService) ListFiles(workID uint64) ([]model.WorkFile, error) {
	var files []model.WorkFile
	err := s.db.Where("work_id = ?", workID).Order("created_at desc").Find(&files).Error
	return files, err
}

func (s *WorkService) ListEvidenceByUser(userID uint64) ([]model.EvidenceRecord, error) {
	var records []model.EvidenceRecord
	err := s.db.Where("user_id = ?", userID).Order("created_at desc").Find(&records).Error
	return records, err
}

func (s *WorkService) ListEvidenceByWork(workID uint64) ([]model.EvidenceRecord, error) {
	var records []model.EvidenceRecord
	err := s.db.Where("work_id = ?", workID).Order("created_at desc").Find(&records).Error
	return records, err
}

func (s *WorkService) ListCertificatesByUser(userID uint64) ([]model.Certificate, error) {
	var certs []model.Certificate
	err := s.db.Where("user_id = ?", userID).Order("created_at desc").Find(&certs).Error
	return certs, err
}

func (s *WorkService) ListCertificatesByWork(workID uint64) ([]model.Certificate, error) {
	var certs []model.Certificate
	err := s.db.Joins("JOIN evidence_records ON evidence_records.id = certificates.evidence_id").
		Where("evidence_records.work_id = ?", workID).
		Order("certificates.created_at desc").
		Find(&certs).Error
	return certs, err
}

func (s *WorkService) PrepareEvidence(userID, workID uint64) (*EvidencePrepareResult, error) {
	if s.cfg.EvidenceRegistryAddress == "" {
		return nil, errors.New("evidence registry address not configured")
	}
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	var work model.Work
	if err := s.db.Where("id = ? AND user_id = ?", workID, userID).First(&work).Error; err != nil {
		return nil, err
	}
	var file model.WorkFile
	if err := s.db.Where("work_id = ? AND user_id = ?", workID, userID).Order("created_at desc").First(&file).Error; err != nil {
		return nil, errors.New("work file required before chain evidence")
	}
	var existing model.EvidenceRecord
	if err := s.db.Where("work_id = ? AND work_file_id = ? AND status IN ?", workID, file.ID, []string{"pending_chain", "confirmed"}).First(&existing).Error; err == nil {
		return &EvidencePrepareResult{
			ContractAddress: existing.ContractAddress,
			EvidenceNo:      existing.EvidenceNo,
			EvidenceNoHash:  existing.EvidenceNoHash,
			FileHash:        existing.FileHash,
			MetadataURI:     existing.MetadataURI,
			ChainID:         s.cfg.ChainID,
		}, nil
	}
	evidenceNo := fmt.Sprintf("EV-%s-%06d", time.Now().Format("20060102"), time.Now().UnixNano()%1000000)
	meta := map[string]interface{}{
		"evidence_no": evidenceNo,
		"title":       work.Title,
		"file_name":   file.FileName,
		"file_hash":   file.SHA256Hash,
		"owner":       user.WalletAddress,
		"created_at":  time.Now().UTC(),
	}
	raw, _ := json.MarshalIndent(meta, "", "  ")
	savedMeta, err := s.store.SaveBytes(fmt.Sprintf("metadata/%d", workID), evidenceNo+".json", raw)
	if err != nil {
		return nil, err
	}
	evidenceNoHash := storage.SHA256Hex([]byte(evidenceNo))
	record := model.EvidenceRecord{
		WorkID:          work.ID,
		WorkFileID:      file.ID,
		UserID:          userID,
		EvidenceNo:      evidenceNo,
		EvidenceNoHash:  evidenceNoHash,
		FileHash:        file.SHA256Hash,
		OwnerAddress:    user.WalletAddress,
		MetadataURI:     savedMeta.URL,
		ContractAddress: s.cfg.EvidenceRegistryAddress,
		Status:          "pending_chain",
	}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return tx.Model(&model.Work{}).Where("id = ?", work.ID).Update("status", "pending_chain").Error
	}); err != nil {
		return nil, err
	}
	return &EvidencePrepareResult{
		ContractAddress: record.ContractAddress,
		EvidenceNo:      evidenceNo,
		EvidenceNoHash:  record.EvidenceNoHash,
		FileHash:        record.FileHash,
		MetadataURI:     record.MetadataURI,
		ChainID:         s.cfg.ChainID,
	}, nil
}

func (s *WorkService) ConfirmEvidence(userID, workID uint64, txHash string, chainEvidenceID uint64) (*model.EvidenceRecord, error) {
	var record model.EvidenceRecord
	if err := s.db.Where("work_id = ? AND user_id = ? AND status = ?", workID, userID, "pending_chain").Order("created_at desc").First(&record).Error; err != nil {
		return nil, err
	}
	txHash = normalizeHash(txHash)
	if s.eth == nil || !s.eth.IsReady() {
		return nil, errors.New("eth client not configured")
	}
	verified, blockNumber, err := s.eth.VerifyEvidenceCreatedTx(context.Background(), txHash, eth.EvidenceCreatedExpectation{
		Owner:          record.OwnerAddress,
		EvidenceNoHash: record.EvidenceNoHash,
		FileHash:       record.FileHash,
		MetadataURI:    record.MetadataURI,
	})
	if err != nil {
		return nil, err
	}
	if chainEvidenceID != 0 && chainEvidenceID != verified {
		return nil, errors.New("chain evidence id mismatch")
	}
	now := time.Now()
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		record.TxHash = &txHash
		record.ChainEvidenceID = &verified
		record.BlockNumber = &blockNumber
		record.Status = "confirmed"
		record.ConfirmedAt = &now
		if err := tx.Save(&record).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Work{}).Where("id = ?", workID).Updates(map[string]interface{}{
			"status":     "confirmed",
			"visibility": "public",
		}).Error; err != nil {
			return err
		}
		_, err := s.rep.RecalculateTx(tx, userID)
		return err
	}); err != nil {
		return nil, err
	}
	_, _ = s.GenerateCertificate(record.ID)
	return &record, nil
}

func (s *WorkService) GenerateCertificate(evidenceID uint64) (*model.Certificate, error) {
	var existing model.Certificate
	if err := s.db.Where("evidence_id = ?", evidenceID).First(&existing).Error; err == nil {
		return &existing, nil
	}
	var evidence model.EvidenceRecord
	if err := s.db.First(&evidence, evidenceID).Error; err != nil {
		return nil, err
	}
	certNo := fmt.Sprintf("CERT-%s-%06d", time.Now().Format("20060102"), evidence.ID)
	verifyURL := fmt.Sprintf("/verify?certificate_no=%s", certNo)
	pdf := simplePDF(fmt.Sprintf("Web3Proof Evidence Certificate\nCertificate No: %s\nEvidence No: %s\nFile Hash: %s\nOwner: %s\nTx Hash: %s\n", certNo, evidence.EvidenceNo, evidence.FileHash, evidence.OwnerAddress, stringValue(evidence.TxHash)))
	saved, err := s.store.SaveBytes("certificates", certNo+".pdf", pdf)
	if err != nil {
		return nil, err
	}
	cert := model.Certificate{
		EvidenceID:    evidence.ID,
		UserID:        evidence.UserID,
		CertificateNo: certNo,
		PDFStorageKey: saved.Key,
		VerifyURL:     verifyURL,
	}
	if err := s.db.Create(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

func simplePDF(text string) []byte {
	escaped := strings.NewReplacer(`\`, `\\`, `(`, `\(`, `)`, `\)`, "\n", `\n`).Replace(text)
	stream := fmt.Sprintf("BT /F1 12 Tf 72 760 Td (%s) Tj ET", escaped)
	return []byte(fmt.Sprintf("%%PDF-1.4\n1 0 obj << /Type /Catalog /Pages 2 0 R >> endobj\n2 0 obj << /Type /Pages /Kids [3 0 R] /Count 1 >> endobj\n3 0 obj << /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] /Resources << /Font << /F1 4 0 R >> >> /Contents 5 0 R >> endobj\n4 0 obj << /Type /Font /Subtype /Type1 /BaseFont /Helvetica >> endobj\n5 0 obj << /Length %d >> stream\n%s\nendstream endobj\ntrailer << /Root 1 0 R >>\n%%%%EOF", len(stream), stream))
}

func normalizeHash(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value != "" && !strings.HasPrefix(value, "0x") {
		value = "0x" + value
	}
	return value
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
