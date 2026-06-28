package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort string
	GinMode    string
	JWTSecret  string
	CORSOrigin string
	AuthDomain string
	AuthURI    string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBDriver   string
	DBFile     string

	RedisAddr     string
	RedisPassword string

	PinataJWT     string
	StorageRoot   string
	PublicBaseURL string

	SepoliaRPCURL string
	ChainID       int64

	BackendWalletPrivateKey string
	AdminWalletAddress      string

	DIDRegistryAddress      string
	ProjectRegistryAddress  string
	SkillSBTAddress         string
	ReputationAddress       string
	DIDProfileAddress       string
	EvidenceRegistryAddress string
	CredentialSBTAddress    string
}

func Load() *Config {
	loadDotEnv()
	chainID, _ := strconv.ParseInt(getEnv("CHAIN_ID", "11155111"), 10, 64)
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		GinMode:    getEnv("GIN_MODE", "debug"),
		JWTSecret:  getEnv("JWT_SECRET", "dev-secret-change-me"),
		CORSOrigin: getEnv("CORS_ORIGIN", "http://localhost:5173"),
		AuthDomain: getEnv("AUTH_DOMAIN", "localhost:5173"),
		AuthURI:    getEnv("AUTH_URI", "http://localhost:5173"),

		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "web3proof"),
		DBPassword: getEnv("DB_PASSWORD", "web3proof"),
		DBName:     getEnv("DB_NAME", "web3proof"),
		DBDriver:   getEnv("DB_DRIVER", "sqlite"),
		DBFile:     getEnv("DB_FILE", "./web3proof.db"),

		RedisAddr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		PinataJWT:     getEnv("PINATA_JWT", ""),
		StorageRoot:   getEnv("STORAGE_ROOT", "./storage"),
		PublicBaseURL: getEnv("PUBLIC_BASE_URL", ""),

		SepoliaRPCURL: getEnv("SEPOLIA_RPC_URL", "https://rpc.sepolia.org"),
		ChainID:       chainID,

		BackendWalletPrivateKey: getEnv("BACKEND_WALLET_PRIVATE_KEY", ""),
		AdminWalletAddress:      getEnv("ADMIN_WALLET_ADDRESS", ""),

		DIDRegistryAddress:      getEnv("DID_REGISTRY_ADDRESS", ""),
		ProjectRegistryAddress:  getEnv("PROJECT_REGISTRY_ADDRESS", ""),
		SkillSBTAddress:         getEnv("SKILL_SBT_ADDRESS", ""),
		ReputationAddress:       getEnv("REPUTATION_ADDRESS", ""),
		DIDProfileAddress:       getEnv("DID_PROFILE_ADDRESS", getEnv("DID_REGISTRY_ADDRESS", "")),
		EvidenceRegistryAddress: getEnv("EVIDENCE_REGISTRY_ADDRESS", getEnv("PROJECT_REGISTRY_ADDRESS", "")),
		CredentialSBTAddress:    getEnv("CREDENTIAL_SBT_ADDRESS", getEnv("SKILL_SBT_ADDRESS", "")),
	}
}

func loadDotEnv() {
	for _, path := range []string{"deploy/.env", "../deploy/.env"} {
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			key, value, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}
			key = strings.TrimSpace(key)
			value = strings.Trim(strings.TrimSpace(value), `"'`)
			if key != "" && os.Getenv(key) == "" {
				_ = os.Setenv(key, value)
			}
		}
		return
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func (c *Config) DSN() string {
	return c.DBUser + ":" + c.DBPassword + "@tcp(" + c.DBHost + ":" + c.DBPort + ")/" + c.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
}

func (c *Config) Validate() error {
	if c.GinMode == "release" && (c.JWTSecret == "" || c.JWTSecret == "dev-secret-change-me") {
		return fmt.Errorf("JWT_SECRET must be set to a non-default value")
	}
	if c.AuthDomain == "" {
		return fmt.Errorf("AUTH_DOMAIN must be set")
	}
	if c.AuthURI == "" {
		return fmt.Errorf("AUTH_URI must be set")
	}
	if c.DBDriver != "sqlite" && c.DBDriver != "mysql" {
		return fmt.Errorf("unsupported DB_DRIVER: %s", c.DBDriver)
	}
	return nil
}
