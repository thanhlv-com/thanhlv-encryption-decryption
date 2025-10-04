package crypto

import (
	"fmt"
	"strings"
	"thanhlv-encryption-decryption/pkg/utils"
)

func NewCryptoProvider(algorithm string) (CryptoProvider, error) {
	utils.DebugLogf("NewCryptoProvider: initializing provider for algorithm: %s", algorithm)
	switch strings.ToLower(algorithm) {
	case "aes-256-cbc":
		return &AESProvider{}, nil
	case "rsa":
		return &RSAProvider{}, nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}