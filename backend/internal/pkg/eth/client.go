package eth

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"web3proof/backend/internal/config"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const reputationABI = `[{"inputs":[{"internalType":"address","name":"user","type":"address"},{"internalType":"uint256","name":"projectScore","type":"uint256"},{"internalType":"uint256","name":"certScore","type":"uint256"},{"internalType":"uint256","name":"activityScore","type":"uint256"}],"name":"updateScore","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
const evidenceRegistryABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"evidenceId","type":"uint256"},{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"bytes32","name":"fileHash","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"evidenceNoHash","type":"bytes32"},{"indexed":false,"internalType":"string","name":"metadataURI","type":"string"}],"name":"EvidenceCreated","type":"event"},{"inputs":[{"internalType":"address","name":"owner","type":"address"}],"name":"getEvidenceIdsByOwner","outputs":[{"internalType":"uint256[]","name":"","type":"uint256[]"}],"stateMutability":"view","type":"function"}]`
const didProfileABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"string","name":"did","type":"string"},{"indexed":false,"internalType":"string","name":"metadataURI","type":"string"}],"name":"ProfileRegistered","type":"event"},{"inputs":[{"internalType":"address","name":"owner","type":"address"}],"name":"hasProfile","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"}]`
const credentialSBTABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":true,"internalType":"uint256","name":"tokenId","type":"uint256"},{"indexed":true,"internalType":"uint256","name":"evidenceId","type":"uint256"},{"indexed":false,"internalType":"string","name":"tokenURIValue","type":"string"}],"name":"CredentialMinted","type":"event"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"ownerOf","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"locked","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"credentialMetadata","outputs":[{"internalType":"uint256","name":"evidenceId","type":"uint256"},{"internalType":"string","name":"tokenURIValue","type":"string"},{"internalType":"uint64","name":"issuedAt","type":"uint64"},{"internalType":"bool","name":"revoked","type":"bool"}],"stateMutability":"view","type":"function"}]`
const projectRegistryABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"projectId","type":"uint256"},{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"string","name":"ipfsCID","type":"string"},{"indexed":false,"internalType":"bytes32","name":"contentHash","type":"bytes32"}],"name":"ProjectAdded","type":"event"}]`
const didRegistryABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"string","name":"did","type":"string"},{"indexed":false,"internalType":"string","name":"metadataCID","type":"string"}],"name":"DIDRegistered","type":"event"},{"inputs":[{"internalType":"address","name":"owner","type":"address"}],"name":"hasProfile","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"}]`
const projectRegistryReadABI = `[{"inputs":[{"internalType":"address","name":"owner","type":"address"}],"name":"getProjectsByOwner","outputs":[{"internalType":"uint256[]","name":"","type":"uint256[]"}],"stateMutability":"view","type":"function"}]`
const reputationReadABI = `[{"inputs":[{"internalType":"address","name":"user","type":"address"}],"name":"getScore","outputs":[{"components":[{"internalType":"uint256","name":"total","type":"uint256"},{"internalType":"uint256","name":"projectScore","type":"uint256"},{"internalType":"uint256","name":"certScore","type":"uint256"},{"internalType":"uint256","name":"activityScore","type":"uint256"},{"internalType":"uint64","name":"updatedAt","type":"uint64"}],"internalType":"struct Reputation.Score","name":"","type":"tuple"}],"stateMutability":"view","type":"function"}]`
const skillSBTABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":true,"internalType":"uint256","name":"tokenId","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"skillType","type":"uint256"}],"name":"SkillMinted","type":"event"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"ownerOf","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"locked","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"skillMetadata","outputs":[{"internalType":"uint256","name":"skillType","type":"uint256"},{"internalType":"string","name":"metadataCID","type":"string"},{"internalType":"uint64","name":"issuedAt","type":"uint64"}],"stateMutability":"view","type":"function"}]`

type ProjectAddedExpectation struct {
	Owner       string
	IPFSCID     string
	ContentHash string
}

type EvidenceCreatedExpectation struct {
	Owner          string
	EvidenceNoHash string
	FileHash       string
	MetadataURI    string
}

type ContractEventLog struct {
	ContractName string
	EventName    string
	Address      string
	TxHash       string
	BlockNumber  uint64
	LogIndex     uint
	Topics       []string
	Data         string
}
type SkillMintedExpectation struct {
	ToAddress   string
	TokenID     uint64
	SkillType   uint64
	MetadataCID string
}

type CredentialMintedExpectation struct {
	ToAddress  string
	TokenID    uint64
	EvidenceID uint64
	TokenURI   string
}

type DIDRegisteredExpectation struct {
	Owner string
	DID   string
}

type OnchainScore struct {
	Total         uint64
	ProjectScore  uint64
	CertScore     uint64
	ActivityScore uint64
	UpdatedAt     uint64
}

type OnchainSBTVerification struct {
	Valid       bool
	Owner       string
	Locked      bool
	SkillType   uint64
	MetadataCID string
}

type EthClient struct {
	cfg    *config.Config
	client *ethclient.Client
	key    *ecdsa.PrivateKey
}

func NewEthClient(cfg *config.Config) (*EthClient, error) {
	if cfg.SepoliaRPCURL == "" || (cfg.DIDProfileAddress == "" && cfg.EvidenceRegistryAddress == "" && cfg.CredentialSBTAddress == "" && cfg.ReputationAddress == "" && cfg.DIDRegistryAddress == "" && cfg.ProjectRegistryAddress == "" && cfg.SkillSBTAddress == "") {
		return &EthClient{cfg: cfg}, nil
	}
	client, err := ethclient.Dial(cfg.SepoliaRPCURL)
	if err != nil {
		return nil, err
	}
	ec := &EthClient{cfg: cfg, client: client}
	if cfg.BackendWalletPrivateKey != "" {
		keyHex := strings.TrimPrefix(cfg.BackendWalletPrivateKey, "0x")
		key, err := crypto.HexToECDSA(keyHex)
		if err != nil {
			return nil, err
		}
		ec.key = key
	}
	return ec, nil
}

func (e *EthClient) IsReady() bool {
	return e != nil && e.client != nil
}

func (e *EthClient) callContract(ctx context.Context, contractAddress, abiJSON, method string, args ...interface{}) ([]interface{}, error) {
	if e.client == nil {
		return nil, fmt.Errorf("rpc not configured")
	}
	if contractAddress == "" {
		return nil, fmt.Errorf("contract address not configured")
	}
	parsed, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, err
	}
	data, err := parsed.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	address := common.HexToAddress(contractAddress)
	result, err := e.client.CallContract(ctx, ethereum.CallMsg{
		To:   &address,
		Data: data,
	}, nil)
	if err != nil {
		return nil, err
	}
	return parsed.Unpack(method, result)
}

func (e *EthClient) HasDIDProfile(ctx context.Context, wallet string) (bool, error) {
	address := e.cfg.DIDProfileAddress
	abiJSON := didProfileABI
	if address == "" {
		address = e.cfg.DIDRegistryAddress
		abiJSON = didRegistryABI
	}
	values, err := e.callContract(ctx, address, abiJSON, "hasProfile", common.HexToAddress(wallet))
	if err != nil {
		return false, err
	}
	if len(values) == 0 {
		return false, fmt.Errorf("empty hasProfile response")
	}
	exists, ok := values[0].(bool)
	if !ok {
		return false, fmt.Errorf("invalid hasProfile response")
	}
	return exists, nil
}

func (e *EthClient) GetProjectIDsByOwner(ctx context.Context, wallet string) ([]uint64, error) {
	if e.cfg.EvidenceRegistryAddress != "" {
		values, err := e.callContract(ctx, e.cfg.EvidenceRegistryAddress, evidenceRegistryABI, "getEvidenceIdsByOwner", common.HexToAddress(wallet))
		if err != nil {
			return nil, err
		}
		return bigIntSlice(values)
	}
	values, err := e.callContract(ctx, e.cfg.ProjectRegistryAddress, projectRegistryReadABI, "getProjectsByOwner", common.HexToAddress(wallet))
	if err != nil {
		return nil, err
	}
	return bigIntSlice(values)
}

func bigIntSlice(values []interface{}) ([]uint64, error) {
	if len(values) == 0 {
		return nil, fmt.Errorf("empty id list response")
	}
	ids := make([]uint64, 0)
	switch raw := values[0].(type) {
	case []*big.Int:
		for _, id := range raw {
			ids = append(ids, id.Uint64())
		}
	default:
		v := reflect.ValueOf(raw)
		if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
			return nil, fmt.Errorf("invalid getProjectsByOwner response")
		}
		for i := 0; i < v.Len(); i++ {
			id, ok := v.Index(i).Interface().(*big.Int)
			if !ok {
				return nil, fmt.Errorf("invalid project id response")
			}
			ids = append(ids, id.Uint64())
		}
	}
	return ids, nil
}

func (e *EthClient) VerifyEvidenceCreatedTx(ctx context.Context, txHash string, expected EvidenceCreatedExpectation) (uint64, uint64, error) {
	if e.client == nil {
		return 0, 0, fmt.Errorf("rpc not configured")
	}
	if e.cfg.EvidenceRegistryAddress == "" {
		return 0, 0, fmt.Errorf("evidence registry address not configured")
	}
	receipt, err := e.client.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return 0, 0, err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return 0, 0, fmt.Errorf("transaction failed")
	}
	registryAddress := common.HexToAddress(e.cfg.EvidenceRegistryAddress)
	expectedOwner := common.HexToAddress(expected.Owner)
	expectedFileHash := common.HexToHash(expected.FileHash)
	expectedEvidenceNoHash := common.HexToHash(expected.EvidenceNoHash)
	parsed, err := abi.JSON(strings.NewReader(evidenceRegistryABI))
	if err != nil {
		return 0, 0, err
	}
	event := parsed.Events["EvidenceCreated"]
	sig := crypto.Keccak256Hash([]byte("EvidenceCreated(uint256,address,bytes32,bytes32,string)"))
	for _, log := range receipt.Logs {
		if log.Address != registryAddress || len(log.Topics) < 4 || log.Topics[0] != sig {
			continue
		}
		owner := common.BytesToAddress(log.Topics[2].Bytes())
		fileHash := log.Topics[3]
		if owner != expectedOwner {
			return 0, 0, fmt.Errorf("evidence owner mismatch")
		}
		if fileHash != expectedFileHash {
			return 0, 0, fmt.Errorf("evidence file hash mismatch")
		}
		values := map[string]interface{}{}
		if err := event.Inputs.NonIndexed().UnpackIntoMap(values, log.Data); err != nil {
			return 0, 0, err
		}
		evidenceNoHash, ok := values["evidenceNoHash"].([32]byte)
		if !ok || common.BytesToHash(evidenceNoHash[:]) != expectedEvidenceNoHash {
			return 0, 0, fmt.Errorf("evidence number hash mismatch")
		}
		metadataURI, ok := values["metadataURI"].(string)
		if !ok || metadataURI != expected.MetadataURI {
			return 0, 0, fmt.Errorf("evidence metadata uri mismatch")
		}
		return log.Topics[1].Big().Uint64(), receipt.BlockNumber.Uint64(), nil
	}
	return 0, 0, fmt.Errorf("EvidenceCreated event not found")
}

func (e *EthClient) VerifyCredentialMintedTx(ctx context.Context, txHash string, expected CredentialMintedExpectation) (uint64, error) {
	if e.client == nil {
		return 0, fmt.Errorf("rpc not configured")
	}
	if e.cfg.CredentialSBTAddress == "" {
		return 0, fmt.Errorf("credential SBT address not configured")
	}
	receipt, err := e.client.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return 0, err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return 0, fmt.Errorf("transaction failed")
	}
	sbtAddress := common.HexToAddress(e.cfg.CredentialSBTAddress)
	expectedTo := common.HexToAddress(expected.ToAddress)
	parsed, err := abi.JSON(strings.NewReader(credentialSBTABI))
	if err != nil {
		return 0, err
	}
	event := parsed.Events["CredentialMinted"]
	sig := crypto.Keccak256Hash([]byte("CredentialMinted(address,uint256,uint256,string)"))
	for _, log := range receipt.Logs {
		if log.Address != sbtAddress || len(log.Topics) < 4 || log.Topics[0] != sig {
			continue
		}
		to := common.BytesToAddress(log.Topics[1].Bytes())
		tokenID := log.Topics[2].Big().Uint64()
		evidenceID := log.Topics[3].Big().Uint64()
		if to != expectedTo {
			return 0, fmt.Errorf("credential recipient mismatch")
		}
		if expected.TokenID != 0 && tokenID != expected.TokenID {
			return 0, fmt.Errorf("credential token id mismatch")
		}
		if evidenceID != expected.EvidenceID {
			return 0, fmt.Errorf("credential evidence id mismatch")
		}
		values := map[string]interface{}{}
		if err := event.Inputs.NonIndexed().UnpackIntoMap(values, log.Data); err != nil {
			return 0, err
		}
		tokenURI, ok := values["tokenURIValue"].(string)
		if !ok || tokenURI != expected.TokenURI {
			return 0, fmt.Errorf("credential token uri mismatch")
		}
		return tokenID, nil
	}
	return 0, fmt.Errorf("CredentialMinted event not found")
}

func (e *EthClient) GetReputationScore(ctx context.Context, wallet string) (*OnchainScore, error) {
	values, err := e.callContract(ctx, e.cfg.ReputationAddress, reputationReadABI, "getScore", common.HexToAddress(wallet))
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("empty getScore response")
	}
	score := values[0]
	return &OnchainScore{
		Total:         uint64Field(score, "Total"),
		ProjectScore:  uint64Field(score, "ProjectScore"),
		CertScore:     uint64Field(score, "CertScore"),
		ActivityScore: uint64Field(score, "ActivityScore"),
		UpdatedAt:     uint64Field(score, "UpdatedAt"),
	}, nil
}

func (e *EthClient) VerifySBTOnchain(ctx context.Context, wallet string, tokenID, expectedSkillType uint64, expectedMetadataCID string) (*OnchainSBTVerification, error) {
	if e.client == nil {
		return nil, fmt.Errorf("rpc not configured")
	}
	if e.cfg.SkillSBTAddress == "" {
		return nil, fmt.Errorf("SkillSBT address not configured")
	}
	token := new(big.Int).SetUint64(tokenID)
	ownerValues, err := e.callContract(ctx, e.cfg.SkillSBTAddress, skillSBTABI, "ownerOf", token)
	if err != nil {
		return nil, err
	}
	owner, ok := ownerValues[0].(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid ownerOf response")
	}
	lockedValues, err := e.callContract(ctx, e.cfg.SkillSBTAddress, skillSBTABI, "locked", token)
	if err != nil {
		return nil, err
	}
	locked, ok := lockedValues[0].(bool)
	if !ok {
		return nil, fmt.Errorf("invalid locked response")
	}
	metadataValues, err := e.callContract(ctx, e.cfg.SkillSBTAddress, skillSBTABI, "skillMetadata", token)
	if err != nil {
		return nil, err
	}
	if len(metadataValues) < 2 {
		return nil, fmt.Errorf("invalid skillMetadata response")
	}
	skillType, ok := metadataValues[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid skill type response")
	}
	metadataCID, ok := metadataValues[1].(string)
	if !ok {
		return nil, fmt.Errorf("invalid metadata cid response")
	}
	valid := strings.EqualFold(owner.Hex(), wallet) && locked && skillType.Uint64() == expectedSkillType
	if expectedMetadataCID != "" {
		valid = valid && metadataCID == expectedMetadataCID
	}
	return &OnchainSBTVerification{
		Valid:       valid,
		Owner:       strings.ToLower(owner.Hex()),
		Locked:      locked,
		SkillType:   skillType.Uint64(),
		MetadataCID: metadataCID,
	}, nil
}

func (e *EthClient) SyncReputation(ctx context.Context, wallet string, projectScore, certScore, activityScore uint) (string, error) {
	if e.client == nil || e.key == nil {
		return "", nil
	}
	parsed, err := abi.JSON(strings.NewReader(reputationABI))
	if err != nil {
		return "", err
	}
	data, err := parsed.Pack("updateScore",
		common.HexToAddress(wallet),
		big.NewInt(int64(projectScore)),
		big.NewInt(int64(certScore)),
		big.NewInt(int64(activityScore)),
	)
	if err != nil {
		return "", err
	}
	chainID := big.NewInt(e.cfg.ChainID)
	from := crypto.PubkeyToAddress(e.key.PublicKey)
	nonce, err := e.client.PendingNonceAt(ctx, from)
	if err != nil {
		return "", err
	}
	gasPrice, err := e.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}
	tx := types.NewTransaction(nonce, common.HexToAddress(e.cfg.ReputationAddress), big.NewInt(0), 300000, gasPrice, data)
	signed, err := types.SignTx(tx, types.NewEIP155Signer(chainID), e.key)
	if err != nil {
		return "", err
	}
	if err := e.client.SendTransaction(ctx, signed); err != nil {
		return "", err
	}
	return signed.Hash().Hex(), nil
}

func (e *EthClient) VerifyProjectAddedTx(ctx context.Context, txHash string, expected ProjectAddedExpectation) (uint64, error) {
	if e.client == nil {
		return 0, fmt.Errorf("rpc not configured")
	}
	if e.cfg.ProjectRegistryAddress == "" {
		return 0, fmt.Errorf("project registry address not configured")
	}
	receipt, err := e.client.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return 0, err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return 0, fmt.Errorf("transaction failed")
	}

	registryAddress := common.HexToAddress(e.cfg.ProjectRegistryAddress)
	expectedOwner := common.HexToAddress(expected.Owner)
	expectedHash := common.HexToHash(expected.ContentHash)
	parsed, err := abi.JSON(strings.NewReader(projectRegistryABI))
	if err != nil {
		return 0, err
	}
	event := parsed.Events["ProjectAdded"]
	projectAddedSig := crypto.Keccak256Hash([]byte("ProjectAdded(uint256,address,string,bytes32)"))
	for _, log := range receipt.Logs {
		if log.Address != registryAddress || len(log.Topics) < 3 || log.Topics[0] != projectAddedSig {
			continue
		}
		owner := common.BytesToAddress(log.Topics[2].Bytes())
		if owner != expectedOwner {
			return 0, fmt.Errorf("project owner mismatch")
		}
		values := map[string]interface{}{}
		if err := event.Inputs.NonIndexed().UnpackIntoMap(values, log.Data); err != nil {
			return 0, err
		}
		ipfsCID, ok := values["ipfsCID"].(string)
		if !ok || ipfsCID != expected.IPFSCID {
			return 0, fmt.Errorf("project ipfs cid mismatch")
		}
		contentHash, ok := values["contentHash"].([32]byte)
		if !ok || common.BytesToHash(contentHash[:]) != expectedHash {
			return 0, fmt.Errorf("project content hash mismatch")
		}
		return log.Topics[1].Big().Uint64(), nil
	}
	return 0, fmt.Errorf("ProjectAdded event not found")
}

func (e *EthClient) VerifyDIDRegisteredTx(ctx context.Context, txHash string, expected DIDRegisteredExpectation) error {
	if e.client == nil {
		return fmt.Errorf("rpc not configured")
	}
	if e.cfg.DIDRegistryAddress == "" {
		return fmt.Errorf("DID registry address not configured")
	}
	receipt, err := e.client.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("transaction failed")
	}

	registryAddress := common.HexToAddress(e.cfg.DIDRegistryAddress)
	expectedOwner := common.HexToAddress(expected.Owner)
	parsed, err := abi.JSON(strings.NewReader(didRegistryABI))
	if err != nil {
		return err
	}
	event := parsed.Events["DIDRegistered"]
	didRegisteredSig := crypto.Keccak256Hash([]byte("DIDRegistered(address,string,string)"))
	for _, log := range receipt.Logs {
		if log.Address != registryAddress || len(log.Topics) < 2 || log.Topics[0] != didRegisteredSig {
			continue
		}
		owner := common.BytesToAddress(log.Topics[1].Bytes())
		if owner != expectedOwner {
			return fmt.Errorf("DID owner mismatch")
		}
		values := map[string]interface{}{}
		if err := event.Inputs.NonIndexed().UnpackIntoMap(values, log.Data); err != nil {
			return err
		}
		did, ok := values["did"].(string)
		if !ok || did != expected.DID {
			return fmt.Errorf("DID value mismatch")
		}
		return nil
	}
	return fmt.Errorf("DIDRegistered event not found")
}

func (e *EthClient) VerifySkillMintedTx(ctx context.Context, txHash string, expected SkillMintedExpectation) (uint64, error) {
	if e.client == nil {
		return 0, fmt.Errorf("rpc not configured")
	}
	if e.cfg.SkillSBTAddress == "" {
		return 0, fmt.Errorf("SkillSBT address not configured")
	}
	receipt, err := e.client.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return 0, err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return 0, fmt.Errorf("transaction failed")
	}

	sbtAddress := common.HexToAddress(e.cfg.SkillSBTAddress)
	expectedTo := common.HexToAddress(expected.ToAddress)
	parsed, err := abi.JSON(strings.NewReader(skillSBTABI))
	if err != nil {
		return 0, err
	}
	event := parsed.Events["SkillMinted"]
	skillMintedSig := crypto.Keccak256Hash([]byte("SkillMinted(address,uint256,uint256)"))
	for _, log := range receipt.Logs {
		if log.Address != sbtAddress || len(log.Topics) < 3 || log.Topics[0] != skillMintedSig {
			continue
		}
		to := common.BytesToAddress(log.Topics[1].Bytes())
		if to != expectedTo {
			return 0, fmt.Errorf("SBT recipient mismatch")
		}
		tokenID := log.Topics[2].Big().Uint64()
		if expected.TokenID != 0 && tokenID != expected.TokenID {
			return 0, fmt.Errorf("SBT token id mismatch")
		}
		values := map[string]interface{}{}
		if err := event.Inputs.NonIndexed().UnpackIntoMap(values, log.Data); err != nil {
			return 0, err
		}
		skillType, ok := values["skillType"].(*big.Int)
		if !ok || skillType.Uint64() != expected.SkillType {
			return 0, fmt.Errorf("SBT skill type mismatch")
		}
		if err := e.verifySBTMetadata(ctx, parsed, sbtAddress, receipt.BlockNumber, tokenID, expected); err != nil {
			return 0, err
		}
		return tokenID, nil
	}
	return 0, fmt.Errorf("SkillMinted event not found")
}

func (e *EthClient) verifySBTMetadata(ctx context.Context, parsed abi.ABI, sbtAddress common.Address, blockNumber *big.Int, tokenID uint64, expected SkillMintedExpectation) error {
	data, err := parsed.Pack("skillMetadata", new(big.Int).SetUint64(tokenID))
	if err != nil {
		return err
	}
	result, err := e.client.CallContract(ctx, ethereum.CallMsg{
		To:   &sbtAddress,
		Data: data,
	}, blockNumber)
	if err != nil {
		return err
	}
	values, err := parsed.Unpack("skillMetadata", result)
	if err != nil {
		return err
	}
	if len(values) < 2 {
		return fmt.Errorf("invalid SBT metadata response")
	}
	skillType, ok := values[0].(*big.Int)
	if !ok || skillType.Uint64() != expected.SkillType {
		return fmt.Errorf("SBT metadata skill type mismatch")
	}
	metadataCID, ok := values[1].(string)
	if !ok || metadataCID != expected.MetadataCID {
		return fmt.Errorf("SBT metadata cid mismatch")
	}
	return nil
}

func Keccak256Hex(data []byte) string {
	return crypto.Keccak256Hash(data).Hex()
}

func MustDecodeHex(s string) []byte {
	s = strings.TrimPrefix(s, "0x")
	b, _ := hex.DecodeString(s)
	return b
}

func uint64Field(value interface{}, name string) uint64 {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return 0
	}
	field := v.FieldByName(name)
	if !field.IsValid() {
		return 0
	}
	if field.Kind() == reflect.Pointer && !field.IsNil() {
		if n, ok := field.Interface().(*big.Int); ok {
			return n.Uint64()
		}
	}
	if field.CanUint() {
		return field.Uint()
	}
	return 0
}

func (e *EthClient) RecentContractEvents(ctx context.Context, lookback uint64) ([]ContractEventLog, error) {
	if e.client == nil {
		return nil, fmt.Errorf("rpc not configured")
	}
	latest, err := e.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	if lookback == 0 {
		lookback = 5000
	}
	latestNumber := latest.Number.Uint64()
	fromBlock := uint64(0)
	if latestNumber > lookback {
		fromBlock = latestNumber - lookback
	}

	contractNames := map[common.Address]string{}
	for name, address := range map[string]string{
		"DIDProfile":       firstNonEmpty(e.cfg.DIDProfileAddress, e.cfg.DIDRegistryAddress),
		"EvidenceRegistry": firstNonEmpty(e.cfg.EvidenceRegistryAddress, e.cfg.ProjectRegistryAddress),
		"CredentialSBT":    firstNonEmpty(e.cfg.CredentialSBTAddress, e.cfg.SkillSBTAddress),
		"Reputation":       e.cfg.ReputationAddress,
	} {
		if address != "" {
			contractNames[common.HexToAddress(address)] = name
		}
	}
	addresses := make([]common.Address, 0, len(contractNames))
	for address := range contractNames {
		addresses = append(addresses, address)
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("no contracts configured")
	}

	logs, err := e.client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		ToBlock:   new(big.Int).SetUint64(latestNumber),
		Addresses: addresses,
	})
	if err != nil {
		return nil, err
	}
	eventNames := map[common.Hash]string{
		crypto.Keccak256Hash([]byte("DIDRegistered(address,string,string)")):                    "DIDRegistered",
		crypto.Keccak256Hash([]byte("ProfileRegistered(address,string,string)")):                "ProfileRegistered",
		crypto.Keccak256Hash([]byte("ProfileUpdated(address,string)")):                          "ProfileUpdated",
		crypto.Keccak256Hash([]byte("ProjectAdded(uint256,address,string,bytes32)")):            "ProjectAdded",
		crypto.Keccak256Hash([]byte("EvidenceCreated(uint256,address,bytes32,bytes32,string)")): "EvidenceCreated",
		crypto.Keccak256Hash([]byte("SkillMinted(address,uint256,uint256)")):                    "SkillMinted",
		crypto.Keccak256Hash([]byte("CredentialMinted(address,uint256,uint256,string)")):        "CredentialMinted",
		crypto.Keccak256Hash([]byte("ScoreUpdated(address,uint256,uint256)")):                   "ScoreUpdated",
	}

	events := make([]ContractEventLog, 0, len(logs))
	for _, item := range logs {
		topics := make([]string, 0, len(item.Topics))
		for _, topic := range item.Topics {
			topics = append(topics, topic.Hex())
		}
		eventName := "Unknown"
		if len(item.Topics) > 0 {
			if name, ok := eventNames[item.Topics[0]]; ok {
				eventName = name
			}
		}
		events = append(events, ContractEventLog{
			ContractName: contractNames[item.Address],
			EventName:    eventName,
			Address:      item.Address.Hex(),
			TxHash:       item.TxHash.Hex(),
			BlockNumber:  item.BlockNumber,
			LogIndex:     item.Index,
			Topics:       topics,
			Data:         hex.EncodeToString(item.Data),
		})
	}
	return events, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
