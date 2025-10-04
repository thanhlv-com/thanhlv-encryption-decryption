package cmd

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"thanhlv-encryption-decryption/pkg/crypto"
	"thanhlv-encryption-decryption/pkg/utils"
)

var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate encryption keys",
	Long:  `Generate encryption keys for various algorithms.`,
	Run:   runKeygen,
}

var (
	keygenAlgorithm   string
	keygenPrivateFile string
	keygenPublicFile  string
	keygenBase64      bool
)

func init() {
	keygenCmd.Flags().StringVarP(&keygenAlgorithm, "algorithm", "a", "aes-256-cbc", "Key generation algorithm (aes-256-cbc, rsa)")
	keygenCmd.Flags().StringVarP(&keygenPrivateFile, "private", "p", "", "Private key output file (RSA only)")
	keygenCmd.Flags().StringVarP(&keygenPublicFile, "public", "u", "", "Public key output file (RSA only)")
	keygenCmd.Flags().BoolVarP(&keygenBase64, "base64", "b", false, "Output key in base64 format")
}

func runKeygen(cmd *cobra.Command, args []string) {
	switch keygenAlgorithm {
	case "aes-256-cbc":
		provider := &crypto.AESProvider{}
		key, err := provider.GenerateKey()
		if err != nil {
			fmt.Printf("Error generating AES key: %v\n", err)
			os.Exit(1)
		}

		if keygenBase64 {
			fmt.Printf("Generated AES-256-CBC key (base64): %s\n", base64.StdEncoding.EncodeToString(key))
		} else {
			fmt.Printf("Generated AES-256-CBC key (hex): %x\n", key)
		}

	case "rsa":
		privateKey, publicKey, err := crypto.GenerateRSAKeyPair()
		if err != nil {
			fmt.Printf("Error generating RSA keys: %v\n", err)
			os.Exit(1)
		}

		privateFile := keygenPrivateFile
		if privateFile == "" {
			privateFile = "private_key_rsa.pem"
		}

		publicFile := keygenPublicFile
		if publicFile == "" {
			publicFile = "public_key_rsa.pem"
		}

		err = utils.WriteFile(privateFile, privateKey)
		if err != nil {
			fmt.Printf("Error writing private key: %v\n", err)
			os.Exit(1)
		}

		err = utils.WriteFile(publicFile, publicKey)
		if err != nil {
			fmt.Printf("Error writing public key: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("RSA key pair generated:\n")
		fmt.Printf("Private key: %s\n", privateFile)
		fmt.Printf("Public key: %s\n", publicFile)

		if keygenBase64 {
			fmt.Printf("\nPrivate key (base64): %s\n", base64.StdEncoding.EncodeToString(privateKey))
			fmt.Printf("Public key (base64): %s\n", base64.StdEncoding.EncodeToString(publicKey))
		}

	default:
		fmt.Printf("Unsupported algorithm: %s\n", keygenAlgorithm)
		os.Exit(1)
	}
}