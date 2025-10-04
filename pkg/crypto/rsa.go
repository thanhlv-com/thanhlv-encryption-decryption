package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"thanhlv-encryption-decryption/pkg/utils"
)

type RSAProvider struct{}

func (r *RSAProvider) Encrypt(data []byte, key []byte) ([]byte, error) {
	utils.DebugLogf("RSAProvider.Encrypt: encrypting %d bytes of data", len(data))
	// Apply byte transfer to the original data before RSA encryption
	transferredData := ApplyByteTransfer(data, key)

	block, _ := pem.Decode(key)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not an RSA public key")
	}

	// For large data, we need to chunk it since RSA has size limitations
	maxChunkSize := publicKey.Size() - 2*sha256.Size - 2
	var encryptedData []byte

	for i := 0; i < len(transferredData); i += maxChunkSize {
		end := i + maxChunkSize
		if end > len(transferredData) {
			end = len(transferredData)
		}

		chunk := transferredData[i:end]
		encryptedChunk, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, chunk, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt chunk: %w", err)
		}

		encryptedData = append(encryptedData, encryptedChunk...)
	}

	return encryptedData, nil
}

func (r *RSAProvider) Decrypt(data []byte, key []byte) ([]byte, error) {
	utils.DebugLogf("RSAProvider.Decrypt: decrypting %d bytes of data", len(data))
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	var privateKey *rsa.PrivateKey
	var err error

	// Try PKCS1 format first
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8 format
		parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		var ok bool
		privateKey, ok = parsedKey.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("key is not an RSA private key")
		}
	}

	// Decrypt in chunks
	chunkSize := privateKey.Size()
	var decryptedData []byte

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}

		chunk := data[i:end]
		decryptedChunk, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, chunk, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt chunk: %w", err)
		}

		decryptedData = append(decryptedData, decryptedChunk...)
	}

	// Reverse the byte transfer applied during encryption
	finalResult := ReverseByteTransfer(decryptedData, key)

	return finalResult, nil
}

func (r *RSAProvider) GenerateKey() ([]byte, error) {
	utils.DebugLog("RSAProvider.GenerateKey: generating new RSA key pair")
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	return privateKeyPEM, nil
}

// Helper function to generate RSA key pair and return both public and private keys
func GenerateRSAKeyPair() ([]byte, []byte, error) {
	utils.DebugLog("GenerateRSAKeyPair: generating RSA key pair")
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// Private key
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Public key
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal public key: %w", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return privateKeyPEM, publicKeyPEM, nil
}