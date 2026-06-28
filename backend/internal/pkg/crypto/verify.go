package crypto

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func VerifyPersonalSign(message, signature, expectedAddress string) error {
	sig, err := hexutil.Decode(signature)
	if err != nil {
		return fmt.Errorf("invalid signature encoding")
	}
	if len(sig) != 65 {
		return fmt.Errorf("invalid signature length")
	}
	if sig[64] >= 27 {
		sig[64] -= 27
	}
	hash := accounts.TextHash([]byte(message))
	pubKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return fmt.Errorf("signature recovery failed")
	}
	recovered := strings.ToLower(crypto.PubkeyToAddress(*pubKey).Hex())
	expected := strings.ToLower(common.HexToAddress(expectedAddress).Hex())
	if recovered != expected {
		return fmt.Errorf("signature address mismatch")
	}
	return nil
}
