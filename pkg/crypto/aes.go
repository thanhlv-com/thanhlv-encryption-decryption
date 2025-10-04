package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"thanhlv-encryption-decryption/pkg/utils"
)

type AESProvider struct{}

func (a *AESProvider) Encrypt(data []byte, key []byte) ([]byte, error) {
	utils.DebugLogf("AES Encrypt: Input data size: %d bytes, key size: %d bytes", len(data), len(key))
	// Ensure key is 32 bytes for AES-256
	keyHash := sha256.Sum256(key)
	finalKey := keyHash[:]
	utils.DebugLog("AES Encrypt: Generated 256-bit key hash")

	block, err := aes.NewCipher(finalKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Pad data to block size first
	paddedData := pkcs7Pad(data, aes.BlockSize)

	// Generate random IV and allocate buffer for IV + padded data
	ciphertext := make([]byte, aes.BlockSize+len(paddedData))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %w", err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], paddedData)

	// Apply byte transfer using the original key
	utils.DebugLogf("AES Encrypt: Applying byte transfer, ciphertext size: %d bytes", len(ciphertext))
	finalResult := ApplyByteTransfer(ciphertext, key)
	utils.DebugLogf("AES Encrypt: Encryption completed, final result size: %d bytes", len(finalResult))

	return finalResult, nil
}

func (a *AESProvider) Decrypt(data []byte, key []byte) ([]byte, error) {
	utils.DebugLogf("AES Decrypt: Input data size: %d bytes, key size: %d bytes", len(data), len(key))
	// First reverse the byte transfer using the original key
	utils.DebugLog("AES Decrypt: Reversing byte transfer")
	reversedData := ReverseByteTransfer(data, key)
	utils.DebugLogf("AES Decrypt: Byte transfer reversed, data size: %d bytes", len(reversedData))

	if len(reversedData) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Ensure key is 32 bytes for AES-256
	keyHash := sha256.Sum256(key)
	finalKey := keyHash[:]

	block, err := aes.NewCipher(finalKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	iv := reversedData[:aes.BlockSize]
	ciphertext := reversedData[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove padding
	unpaddedData, err := pkcs7Unpad(plaintext, aes.BlockSize)
	if err != nil {
		return nil, fmt.Errorf("failed to remove padding: %w", err)
	}

	return unpaddedData, nil
}

func (a *AESProvider) GenerateKey() ([]byte, error) {
	key := make([]byte, 32) // 256 bits
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}
	return key, nil
}

// PKCS7 padding
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("invalid padding")
	}

	unpadding := int(data[length-1])
	if unpadding > blockSize || unpadding == 0 {
		return nil, fmt.Errorf("invalid padding")
	}

	for i := length - unpadding; i < length; i++ {
		if data[i] != byte(unpadding) {
			return nil, fmt.Errorf("invalid padding")
		}
	}

	return data[:(length - unpadding)], nil
}