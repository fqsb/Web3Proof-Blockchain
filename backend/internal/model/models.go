package model

import "time"

type User struct {
	ID              uint64     `gorm:"primaryKey" json:"id"`
	WalletAddress   string     `gorm:"size:42;uniqueIndex;not null" json:"wallet_address"`
	DID             *string    `gorm:"size:128;uniqueIndex" json:"did"`
	Nickname        *string    `gorm:"size:64" json:"nickname"`
	AvatarURL       *string    `gorm:"size:512" json:"avatar_url"`
	Bio             *string    `gorm:"type:text" json:"bio"`
	Email           *string    `gorm:"size:128" json:"email"`
	ActiveRole      string     `gorm:"size:20;default:creator;index" json:"active_role"`
	IsDIDRegistered bool       `gorm:"default:0" json:"is_did_registered"`
	DIDTxHash       *string    `gorm:"size:66" json:"did_tx_hash"`
	Status          int8       `gorm:"default:1" json:"status"`
	LastActiveAt    *time.Time `json:"last_active_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Roles           []UserRole `gorm:"foreignKey:UserID" json:"roles,omitempty"`
}

func (User) TableName() string { return "users" }

type UserRole struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"uniqueIndex:uk_user_role;not null" json:"user_id"`
	RoleCode  string    `gorm:"size:20;uniqueIndex:uk_user_role;not null;index" json:"role_code"`
	Enabled   bool      `gorm:"default:1" json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserRole) TableName() string { return "user_roles" }

type WalletAccount struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	UserID        uint64    `gorm:"index;not null" json:"user_id"`
	WalletAddress string    `gorm:"size:42;uniqueIndex;not null" json:"wallet_address"`
	IsPrimary     bool      `gorm:"default:1" json:"is_primary"`
	CreatedAt     time.Time `json:"created_at"`
}

func (WalletAccount) TableName() string { return "wallet_accounts" }

type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"size:32;uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"size:64;not null" json:"name"`
	Description *string   `gorm:"type:text" json:"description"`
	IsActive    bool      `gorm:"default:1" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Category) TableName() string { return "categories" }

type Work struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	UserID      uint64    `gorm:"index;not null" json:"user_id"`
	Title       string    `gorm:"size:128;not null" json:"title"`
	Description *string   `gorm:"type:text" json:"description"`
	CategoryID  *uint     `gorm:"index" json:"category_id"`
	ExternalURL *string   `gorm:"size:512" json:"external_url"`
	Visibility  string    `gorm:"size:20;default:private;index" json:"visibility"`
	Status      string    `gorm:"size:20;default:draft;index" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Category    Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (Work) TableName() string { return "works" }

type WorkFile struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	WorkID     uint64    `gorm:"index;not null" json:"work_id"`
	UserID     uint64    `gorm:"index;not null" json:"user_id"`
	FileName   string    `gorm:"size:256;not null" json:"file_name"`
	FileType   string    `gorm:"size:80" json:"file_type"`
	StorageKey string    `gorm:"size:512;not null" json:"storage_key"`
	StorageURL *string   `gorm:"size:512" json:"storage_url"`
	FileSize   int64     `json:"file_size"`
	SHA256Hash string    `gorm:"size:66;index;not null" json:"sha256_hash"`
	CreatedAt  time.Time `json:"created_at"`
}

func (WorkFile) TableName() string { return "work_files" }

type EvidenceRecord struct {
	ID              uint64     `gorm:"primaryKey" json:"id"`
	WorkID          uint64     `gorm:"index;not null" json:"work_id"`
	WorkFileID      uint64     `gorm:"index;not null" json:"work_file_id"`
	UserID          uint64     `gorm:"index;not null" json:"user_id"`
	EvidenceNo      string     `gorm:"size:40;uniqueIndex;not null" json:"evidence_no"`
	EvidenceNoHash  string     `gorm:"size:66;not null" json:"evidence_no_hash"`
	FileHash        string     `gorm:"size:66;index;not null" json:"file_hash"`
	OwnerAddress    string     `gorm:"size:42;index;not null" json:"owner_address"`
	MetadataURI     string     `gorm:"size:512" json:"metadata_uri"`
	ChainEvidenceID *uint64    `json:"chain_evidence_id"`
	ContractAddress string     `gorm:"size:42" json:"contract_address"`
	TxHash          *string    `gorm:"size:66;uniqueIndex" json:"tx_hash"`
	BlockNumber     *uint64    `json:"block_number"`
	Status          string     `gorm:"size:20;default:pending_chain;index" json:"status"`
	ConfirmedAt     *time.Time `json:"confirmed_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (EvidenceRecord) TableName() string { return "evidence_records" }

type Certificate struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	EvidenceID    uint64    `gorm:"index;not null" json:"evidence_id"`
	UserID        uint64    `gorm:"index;not null" json:"user_id"`
	CertificateNo string    `gorm:"size:40;uniqueIndex;not null" json:"certificate_no"`
	PDFStorageKey string    `gorm:"size:512;not null" json:"pdf_storage_key"`
	VerifyURL     string    `gorm:"size:512;not null" json:"verify_url"`
	CreatedAt     time.Time `json:"created_at"`
}

func (Certificate) TableName() string { return "certificates" }

type CertificationApplication struct {
	ID            uint64     `gorm:"primaryKey" json:"id"`
	UserID        uint64     `gorm:"index;not null" json:"user_id"`
	WorkID        uint64     `gorm:"index;not null" json:"work_id"`
	EvidenceID    uint64     `gorm:"index;not null" json:"evidence_id"`
	SkillID       uint64     `gorm:"index;not null;default:1" json:"skill_id"`
	MaterialsCID  string     `gorm:"column:materials_c_id;size:512;not null" json:"materials_cid"`
	MaterialsDesc *string    `gorm:"type:text" json:"materials_desc"`
	Status        string     `gorm:"size:20;default:pending;index" json:"status"`
	ReviewerID    *uint64    `json:"reviewer_id"`
	ReviewNote    *string    `gorm:"type:text" json:"review_note"`
	ReviewedAt    *time.Time `json:"reviewed_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Work          Work       `gorm:"foreignKey:WorkID" json:"work,omitempty"`
	User          User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (CertificationApplication) TableName() string { return "certification_applications" }

type SBTCredential struct {
	ID              uint64    `gorm:"primaryKey" json:"id"`
	UserID          uint64    `gorm:"index;not null" json:"user_id"`
	WorkID          uint64    `gorm:"index;not null" json:"work_id"`
	EvidenceID      uint64    `gorm:"index;not null" json:"evidence_id"`
	ApplicationID   *uint64   `json:"application_id"`
	TokenID         uint64    `gorm:"not null" json:"token_id"`
	ContractAddress string    `gorm:"size:42;not null" json:"contract_address"`
	TxHash          string    `gorm:"size:66;uniqueIndex;not null" json:"tx_hash"`
	TokenURI        string    `gorm:"size:512" json:"token_uri"`
	Status          string    `gorm:"size:20;default:active;index" json:"status"`
	MintedAt        time.Time `json:"minted_at"`
}

func (SBTCredential) TableName() string { return "sbt_credentials" }

type VerifierProfile struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	UserID       uint64    `gorm:"uniqueIndex;not null" json:"user_id"`
	OrgName      string    `gorm:"size:128;not null" json:"org_name"`
	Industry     *string   `gorm:"size:64" json:"industry"`
	ContactEmail *string   `gorm:"size:128" json:"contact_email"`
	Website      *string   `gorm:"size:256" json:"website"`
	Status       string    `gorm:"size:20;default:approved;index" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (VerifierProfile) TableName() string { return "verifier_profiles" }

type VerificationReport struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	ViewerID     *uint64   `gorm:"index" json:"viewer_id"`
	TargetUserID *uint64   `gorm:"index" json:"target_user_id,omitempty"`
	QueryType    string    `gorm:"size:30;not null;default:legacy" json:"query_type"`
	QueryValue   string    `gorm:"size:512;not null;default:''" json:"query_value"`
	Passed       bool      `gorm:"default:0" json:"passed"`
	ReportJSON   string    `gorm:"type:text;not null;default:'{}'" json:"report_json"`
	CreatedAt    time.Time `json:"created_at"`
}

func (VerificationReport) TableName() string { return "verification_reports" }

type ReputationScore struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	UserID        uint64    `gorm:"uniqueIndex;not null" json:"user_id"`
	ProjectScore  uint      `gorm:"default:0" json:"project_score"`
	CertScore     uint      `gorm:"default:0" json:"cert_score"`
	ActivityScore uint      `gorm:"default:0" json:"activity_score"`
	TotalScore    uint      `gorm:"default:0" json:"total_score"`
	Grade         string    `gorm:"size:1;default:D" json:"grade"`
	ChainTxHash   *string   `gorm:"size:66" json:"chain_tx_hash"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (ReputationScore) TableName() string { return "reputation_scores" }

type ChainNetwork struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"size:32;uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"size:64;not null" json:"name"`
	ChainID     int64     `gorm:"not null" json:"chain_id"`
	RPCURL      string    `gorm:"size:512" json:"rpc_url"`
	ExplorerURL string    `gorm:"size:512" json:"explorer_url"`
	IsActive    bool      `gorm:"default:1" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (ChainNetwork) TableName() string { return "chain_networks" }

type ContractConfig struct {
	ID              uint64    `gorm:"primaryKey" json:"id"`
	NetworkCode     string    `gorm:"size:32;index;not null" json:"network_code"`
	ContractName    string    `gorm:"size:64;not null" json:"contract_name"`
	ContractAddress string    `gorm:"size:42;not null" json:"contract_address"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (ContractConfig) TableName() string { return "contract_configs" }

type ChainEvent struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	ContractName string    `gorm:"size:64;not null" json:"contract_name"`
	EventName    string    `gorm:"size:64;not null" json:"event_name"`
	TxHash       string    `gorm:"size:66;not null;index" json:"tx_hash"`
	BlockNumber  uint64    `gorm:"not null" json:"block_number"`
	LogIndex     uint      `gorm:"not null" json:"log_index"`
	Payload      string    `gorm:"type:text" json:"payload"`
	Processed    bool      `gorm:"default:0;index" json:"processed"`
	CreatedAt    time.Time `json:"created_at"`
}

func (ChainEvent) TableName() string { return "chain_events" }

type AuditLog struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    *uint64   `gorm:"index" json:"user_id"`
	Action    string    `gorm:"size:64;not null" json:"action"`
	Resource  string    `gorm:"size:128" json:"resource"`
	Detail    string    `gorm:"type:text" json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

func (AuditLog) TableName() string { return "audit_logs" }
