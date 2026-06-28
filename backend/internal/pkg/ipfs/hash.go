package ipfs

import (
	"github.com/ethereum/go-ethereum/crypto"
)

func ContentHashHex(body []byte) string {
	hash := crypto.Keccak256Hash(body)
	return hash.Hex()
}
