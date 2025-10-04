package crypto

type CryptoProvider interface {
	Encrypt(data []byte, key []byte) ([]byte, error)
	Decrypt(data []byte, key []byte) ([]byte, error)
	GenerateKey() ([]byte, error)
}