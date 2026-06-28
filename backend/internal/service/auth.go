package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"web3proof/backend/internal/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type nonceEntry struct {
	value     string
	expiresAt time.Time
}

type AuthService struct {
	db          *gorm.DB
	redis       *redis.Client
	chainID     int64
	adminWallet string
	authDomain  string
	authURI     string
	memNonces   map[string]nonceEntry
	nonceMu     sync.Mutex
	redisReady  bool
}

type SIWEMessage struct {
	Domain         string
	Address        string
	URI            string
	Version        string
	ChainID        int64
	Nonce          string
	IssuedAt       time.Time
	ExpirationTime time.Time
}

func NewAuthService(db *gorm.DB, redisClient *redis.Client, chainID int64, adminWallet, authDomain, authURI string) *AuthService {
	s := &AuthService{
		db:          db,
		redis:       redisClient,
		chainID:     chainID,
		adminWallet: strings.ToLower(adminWallet),
		authDomain:  authDomain,
		authURI:     authURI,
		memNonces:   map[string]nonceEntry{},
	}
	if redisClient != nil {
		s.redisReady = redisClient.Ping(context.Background()).Err() == nil
	}
	return s
}

func (s *AuthService) setNonce(ctx context.Context, address, nonce string) error {
	if s.redisReady {
		return s.redis.Set(ctx, "nonce:"+address, nonce, 5*time.Minute).Err()
	}
	s.nonceMu.Lock()
	defer s.nonceMu.Unlock()
	s.memNonces["nonce:"+address] = nonceEntry{value: nonce, expiresAt: time.Now().Add(5 * time.Minute)}
	return nil
}

func (s *AuthService) getNonce(ctx context.Context, address string) (string, error) {
	if s.redisReady {
		return s.redis.Get(ctx, "nonce:"+address).Result()
	}
	s.nonceMu.Lock()
	defer s.nonceMu.Unlock()
	entry, ok := s.memNonces["nonce:"+address]
	if !ok || time.Now().After(entry.expiresAt) {
		return "", fmt.Errorf("nonce expired or not found")
	}
	return entry.value, nil
}

func (s *AuthService) delNonce(ctx context.Context, address string) error {
	if s.redisReady {
		return s.redis.Del(ctx, "nonce:"+address).Err()
	}
	s.nonceMu.Lock()
	defer s.nonceMu.Unlock()
	delete(s.memNonces, "nonce:"+address)
	return nil
}

func (s *AuthService) GenerateNonce(ctx context.Context, address string) (string, string, error) {
	address = strings.ToLower(address)
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	nonce := hex.EncodeToString(b)
	if err := s.setNonce(ctx, address, nonce); err != nil {
		return "", "", err
	}
	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(5 * time.Minute)
	message := fmt.Sprintf(`%s wants you to sign in with your Ethereum account:
%s

Sign in to Web3Proof.

URI: %s
Version: 1
Chain ID: %d
Nonce: %s
Issued At: %s
Expiration Time: %s`, s.authDomain, address, s.authURI, s.chainID, nonce, issuedAt.Format(time.RFC3339), expiresAt.Format(time.RFC3339))
	return nonce, message, nil
}

func (s *AuthService) ValidateSIWEMessage(ctx context.Context, address, message string) (*SIWEMessage, error) {
	address = strings.ToLower(address)
	parsed, err := parseSIWEMessage(message)
	if err != nil {
		return nil, err
	}
	if parsed.Domain != s.authDomain {
		return nil, fmt.Errorf("invalid sign-in domain")
	}
	if strings.ToLower(parsed.Address) != address {
		return nil, fmt.Errorf("message address mismatch")
	}
	if parsed.URI != s.authURI {
		return nil, fmt.Errorf("invalid sign-in uri")
	}
	if parsed.Version != "1" {
		return nil, fmt.Errorf("invalid SIWE version")
	}
	if parsed.ChainID != s.chainID {
		return nil, fmt.Errorf("invalid chain id")
	}
	stored, err := s.getNonce(ctx, address)
	if err != nil {
		return nil, err
	}
	if parsed.Nonce != stored {
		return nil, fmt.Errorf("invalid nonce")
	}
	now := time.Now().UTC()
	if parsed.IssuedAt.After(now.Add(2 * time.Minute)) {
		return nil, fmt.Errorf("message issued in the future")
	}
	if !parsed.ExpirationTime.After(now) {
		return nil, fmt.Errorf("message expired")
	}
	return parsed, nil
}

func (s *AuthService) ConsumeNonce(ctx context.Context, address string) error {
	address = strings.ToLower(address)
	return s.delNonce(ctx, address)
}

func parseSIWEMessage(message string) (*SIWEMessage, error) {
	lines := strings.Split(strings.ReplaceAll(message, "\r\n", "\n"), "\n")
	if len(lines) < 10 {
		return nil, fmt.Errorf("invalid SIWE message")
	}
	domain, ok := strings.CutSuffix(lines[0], " wants you to sign in with your Ethereum account:")
	if !ok || domain == "" {
		return nil, fmt.Errorf("invalid SIWE domain line")
	}
	address := strings.TrimSpace(lines[1])
	if address == "" {
		return nil, fmt.Errorf("missing SIWE address")
	}
	fields := map[string]string{}
	for _, line := range lines[2:] {
		key, value, ok := strings.Cut(line, ": ")
		if ok {
			fields[key] = strings.TrimSpace(value)
		}
	}
	chainID, err := strconv.ParseInt(fields["Chain ID"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid chain id")
	}
	issuedAt, err := time.Parse(time.RFC3339, fields["Issued At"])
	if err != nil {
		return nil, fmt.Errorf("invalid issued at")
	}
	expirationTime, err := time.Parse(time.RFC3339, fields["Expiration Time"])
	if err != nil {
		return nil, fmt.Errorf("invalid expiration time")
	}
	return &SIWEMessage{
		Domain:         domain,
		Address:        address,
		URI:            fields["URI"],
		Version:        fields["Version"],
		ChainID:        chainID,
		Nonce:          fields["Nonce"],
		IssuedAt:       issuedAt,
		ExpirationTime: expirationTime,
	}, nil
}

func (s *AuthService) FindOrCreateUser(address string) (*model.User, error) {
	address = strings.ToLower(address)
	var user model.User
	err := s.db.Where("wallet_address = ?", address).First(&user).Error
	if err == nil {
		now := time.Now()
		user.LastActiveAt = &now
		if s.adminWallet != "" && strings.EqualFold(address, s.adminWallet) {
			user.ActiveRole = "admin"
			_ = s.ensureRoles(user.ID, []string{"creator", "verifier", "auditor", "admin"})
		}
		s.db.Save(&user)
		return &user, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	did := fmt.Sprintf("did:web3proof:%s", address)
	activeRole := "creator"
	roles := []string{"creator"}
	if s.adminWallet != "" && strings.EqualFold(address, s.adminWallet) {
		activeRole = "admin"
		roles = []string{"creator", "verifier", "auditor", "admin"}
	}
	user = model.User{
		WalletAddress: address,
		DID:           &did,
		ActiveRole:    activeRole,
		Status:        1,
	}
	now := time.Now()
	user.LastActiveAt = &now
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.WalletAccount{UserID: user.ID, WalletAddress: address, IsPrimary: true}).Error; err != nil {
			return err
		}
		for _, role := range roles {
			if err := tx.Create(&model.UserRole{UserID: user.ID, RoleCode: role, Enabled: true}).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) ensureRoles(userID uint64, roles []string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, role := range roles {
			record := model.UserRole{UserID: userID, RoleCode: role, Enabled: true}
			if err := tx.Where("user_id = ? AND role_code = ?", userID, role).FirstOrCreate(&record).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
