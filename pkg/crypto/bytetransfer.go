package crypto

import (
	"thanhlv-encryption-decryption/pkg/utils"
)

// ApplyByteTransfer applies byte-level transformation to data using key bytes in rotation
// For encryption: adds key bytes to data bytes
func ApplyByteTransfer(data []byte, key []byte) []byte {
	utils.DebugLogf("ApplyByteTransfer: processing %d bytes of data with key length %d", len(data), len(key))
	if len(key) == 0 {
		return data
	}

	result := make([]byte, len(data))
	keyLen := len(key)

	for i, dataByte := range data {
		keyByte := key[i%keyLen] // Rotate through key bytes
		result[i] = dataByte + keyByte
	}

	return result
}

// ReverseByteTransfer reverses the byte-level transformation
// For decryption: subtracts key bytes from data bytes
func ReverseByteTransfer(data []byte, key []byte) []byte {
	utils.DebugLogf("ReverseByteTransfer: processing %d bytes of data with key length %d", len(data), len(key))
	if len(key) == 0 {
		return data
	}

	result := make([]byte, len(data))
	keyLen := len(key)

	for i, dataByte := range data {
		keyByte := key[i%keyLen] // Rotate through key bytes
		result[i] = dataByte - keyByte
	}

	return result
}