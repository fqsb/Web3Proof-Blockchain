package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"web3proof/backend/internal/config"
	"web3proof/backend/internal/model"
	"web3proof/backend/internal/pkg/eth"
	"web3proof/backend/internal/pkg/storage"

	"gorm.io/gorm"
)

type CertificationService struct {
	db    *gorm.DB
	store *storage.LocalStore
	cfg   *config.Config
	eth   *eth.EthClient
	rep   *ReputationService
}

func NewCertificationService(db *gorm.DB, store *storage.LocalStore, cfg *config.Config, ethClient *eth.EthClient, rep *ReputationService) *CertificationService {
	return &CertificationService{db: db, store: store, cfg: cfg, eth: ethClient, rep: rep}
}

type ApplyInput struct {
	WorkID        uint64
	EvidenceID    uint64
	MaterialsDesc string
}

type CredentialMintPrepareResult struct {
	ContractAddress string `json:"contract_address"`
	ToAddress       string `json:"to_address"`
	EvidenceID      uint64 `json:"evidence_id"`
	TokenURI        string `json:"token_uri"`
	ChainID         int64  `json:"chain_id"`
}

func (s *CertificationService) Apply(userID uint64, in ApplyInput) (*model.CertificationApplication, error) {
	var evidence model.EvidenceRecord
	if err := s.db.Where("id = ? AND work_id = ? AND user_id = ? AND status = ?", in.EvidenceID, in.WorkID, userID, "confirmed").First(&evidence).Error; err != nil {
		return nil, errors.New("confirmed evidence required before certification")
	}
	var existing model.CertificationApplication
	if err := s.db.Where("user_id = ? AND evidence_id = ? AND status IN ?", userID, in.EvidenceID, []string{"pending", "approved", "minting", "minted"}).First(&existing).Error; err == nil {
		return nil, errors.New("application already exists")
	}
	desc := strings.TrimSpace(in.MaterialsDesc)
	app := model.CertificationApplication{
		UserID:       userID,
		WorkID:       in.WorkID,
		EvidenceID:   in.EvidenceID,
		SkillID:      1,
		MaterialsCID: fmt.Sprintf("evidence:%d", in.EvidenceID),
		Status:       "pending",
	}
	if desc != "" {
		app.MaterialsDesc = &desc
	}
	if err := s.db.Create(&app).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (s *CertificationService) ListMy(userID uint64) ([]model.CertificationApplication, error) {
	var apps []model.CertificationApplication
	err := s.db.Preload("Work").Where("user_id = ?", userID).Order("created_at desc").Find(&apps).Error
	return apps, err
}

func (s *CertificationService) ListPending() ([]model.CertificationApplication, error) {
	var apps []model.CertificationApplication
	err := s.db.Preload("Work").Preload("User").Where("status IN ?", []string{"pending", "approved", "minting"}).Order("created_at asc").Find(&apps).Error
	return apps, err
}

func (s *CertificationService) Review(appID, reviewerID uint64, status, note string) (*model.CertificationApplication, error) {
	if status != "approved" && status != "rejected" {
		return nil, errors.New("invalid review status")
	}
	var app model.CertificationApplication
	now := time.Now()
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&app, appID).Error; err != nil {
			return err
		}
		if app.Status != "pending" {
			return errors.New("application already reviewed")
		}
		app.Status = status
		app.ReviewerID = &reviewerID
		app.ReviewNote = &note
		app.ReviewedAt = &now
		return tx.Save(&app).Error
	}); err != nil {
		return nil, err
	}
	return &app, nil
}

func (s *CertificationService) PrepareMint(appID uint64) (*CredentialMintPrepareResult, error) {
	if s.cfg.CredentialSBTAddress == "" {
		return nil, errors.New("credential SBT address not configured")
	}
	var app model.CertificationApplication
	if err := s.db.Preload("User").Preload("Work").First(&app, appID).Error; err != nil {
		return nil, err
	}
	if app.Status != "approved" {
		return nil, errors.New("application not approved")
	}
	meta := map[string]interface{}{
		"name":        fmt.Sprintf("Web3Proof Credential #%d", app.EvidenceID),
		"description": "Web3Proof certified digital work credential",
		"work_title":  app.Work.Title,
		"evidence_id": app.EvidenceID,
		"issued_to":   app.User.WalletAddress,
	}
	raw, _ := json.MarshalIndent(meta, "", "  ")
	saved, err := s.store.SaveBytes("sbt", fmt.Sprintf("credential-%d.json", app.ID), raw)
	if err != nil {
		return nil, err
	}
	return &CredentialMintPrepareResult{
		ContractAddress: s.cfg.CredentialSBTAddress,
		ToAddress:       app.User.WalletAddress,
		EvidenceID:      app.EvidenceID,
		TokenURI:        saved.URL,
		ChainID:         s.cfg.ChainID,
	}, nil
}

func (s *CertificationService) ConfirmMint(appID uint64, txHash string, tokenID uint64, tokenURI string) (*model.SBTCredential, error) {
	var app model.CertificationApplication
	if err := s.db.Preload("User").First(&app, appID).Error; err != nil {
		return nil, err
	}
	if app.Status != "approved" {
		return nil, errors.New("application not approved")
	}
	if s.eth == nil || !s.eth.IsReady() {
		return nil, errors.New("eth client not configured")
	}
	txHash = normalizeHash(txHash)
	verifiedTokenID, err := s.eth.VerifyCredentialMintedTx(context.Background(), txHash, eth.CredentialMintedExpectation{
		ToAddress:  app.User.WalletAddress,
		TokenID:    tokenID,
		EvidenceID: app.EvidenceID,
		TokenURI:   tokenURI,
	})
	if err != nil {
		return nil, err
	}
	record := model.SBTCredential{
		UserID:          app.UserID,
		WorkID:          app.WorkID,
		EvidenceID:      app.EvidenceID,
		ApplicationID:   &app.ID,
		TokenID:         verifiedTokenID,
		ContractAddress: s.cfg.CredentialSBTAddress,
		TxHash:          txHash,
		TokenURI:        tokenURI,
		Status:          "active",
		MintedAt:        time.Now(),
	}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.CertificationApplication{}).Where("id = ?", app.ID).Update("status", "minted").Error; err != nil {
			return err
		}
		_, err := s.rep.RecalculateTx(tx, app.UserID)
		return err
	}); err != nil {
		return nil, err
	}
	return &record, nil
}
