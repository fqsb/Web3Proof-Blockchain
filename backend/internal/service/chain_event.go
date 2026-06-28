package service

import (
	"context"
	"encoding/json"

	"web3proof/backend/internal/model"
	"web3proof/backend/internal/pkg/eth"

	"gorm.io/gorm"
)

type ChainEventService struct {
	db  *gorm.DB
	eth *eth.EthClient
}

type ChainEventSyncResult struct {
	Scanned  int `json:"scanned"`
	Inserted int `json:"inserted"`
}

func NewChainEventService(db *gorm.DB, ethClient *eth.EthClient) *ChainEventService {
	return &ChainEventService{db: db, eth: ethClient}
}

func (s *ChainEventService) SyncRecent(ctx context.Context, lookback uint64) (*ChainEventSyncResult, error) {
	events, err := s.eth.RecentContractEvents(ctx, lookback)
	if err != nil {
		return nil, err
	}
	result := &ChainEventSyncResult{Scanned: len(events)}
	for _, event := range events {
		payload, err := json.Marshal(map[string]interface{}{
			"address": event.Address,
			"topics":  event.Topics,
			"data":    event.Data,
		})
		if err != nil {
			return nil, err
		}
		record := model.ChainEvent{
			ContractName: event.ContractName,
			EventName:    event.EventName,
			TxHash:       event.TxHash,
			BlockNumber:  event.BlockNumber,
			LogIndex:     event.LogIndex,
			Payload:      string(payload),
			Processed:    false,
		}
		tx := s.db.Where("tx_hash = ? AND log_index = ?", event.TxHash, event.LogIndex).FirstOrCreate(&record)
		if tx.Error != nil {
			return nil, tx.Error
		}
		if tx.RowsAffected > 0 {
			result.Inserted++
		}
	}
	return result, nil
}

func (s *ChainEventService) ListRecent(limit int) ([]model.ChainEvent, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	var events []model.ChainEvent
	err := s.db.Order("block_number desc, log_index desc").Limit(limit).Find(&events).Error
	return events, err
}
