package cmd

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"thanhlv-encryption-decryption/pkg/crypto"
	"thanhlv-encryption-decryption/pkg/utils"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt text or files",
	Long:  `Decrypt text or files using various algorithms like AES-256-CBC or RSA.`,
	Run:   runDecrypt,
}

var (
	decryptAlgorithm string
	decryptInput     string
	decryptOutput    string
	decryptKey       string
	decryptKeyEnv    string
	decryptText      string
	decryptFile      string
)

func init() {
	decryptCmd.Flags().StringVarP(&decryptAlgorithm, "algorithm", "a", "aes-256-cbc", "Decryption algorithm (aes-256-cbc, rsa)")
	decryptCmd.Flags().StringVarP(&decryptKey, "key", "k", "", "Decryption key (base64 encoded)")
	decryptCmd.Flags().StringVarP(&decryptKeyEnv, "key-env", "e", "", "Environment variable name containing the decryption key (base64 encoded)")
	decryptCmd.Flags().StringVarP(&decryptText, "text", "t", "", "Base64 encoded encrypted text to decrypt")
	decryptCmd.Flags().StringVarP(&decryptFile, "file", "f", "", "Encrypted file to decrypt")
	decryptCmd.Flags().StringVarP(&decryptOutput, "output", "o", "", "Output file (optional)")
}

func runDecrypt(cmd *cobra.Command, args []string) {
	if decryptText == "" && decryptFile == "" {
		fmt.Println("Error: Either --text or --file must be specified")
		os.Exit(1)
	}

	if decryptText != "" && decryptFile != "" {
		fmt.Println("Error: Cannot specify both --text and --file")
		os.Exit(1)
	}

	// Validate key input
	if decryptKey == "" && decryptKeyEnv == "" {
		fmt.Println("Error: Either --key or --key-env must be specified")
		os.Exit(1)
	}

	if decryptKey != "" && decryptKeyEnv != "" {
		fmt.Println("Error: Cannot specify both --key and --key-env")
		os.Exit(1)
	}

	// Get the key value
	var keyValue string
	if decryptKeyEnv != "" {
		keyValue = os.Getenv(decryptKeyEnv)
		if keyValue == "" {
			fmt.Printf("Error: Environment variable '%s' is not set or empty\n", decryptKeyEnv)
			os.Exit(1)
		}
	} else {
		keyValue = decryptKey
	}

	// Decode base64 key
	keyBytes, err := base64.StdEncoding.DecodeString(keyValue)
	if err != nil {
		fmt.Printf("Error decoding base64 key: %v\n", err)
		os.Exit(1)
	}

	// Initialize crypto provider
	provider, err := crypto.NewCryptoProvider(decryptAlgorithm)
	if err != nil {
		fmt.Printf("Error initializing crypto provider: %v\n", err)
		os.Exit(1)
	}

	var result []byte

	if decryptText != "" {
		// Decode base64 encrypted text
		encryptedData, err := base64.StdEncoding.DecodeString(decryptText)
		if err != nil {
			fmt.Printf("Error decoding base64 encrypted text: %v\n", err)
			os.Exit(1)
		}

		// Decrypt text
		result, err = provider.Decrypt(encryptedData, keyBytes)
		if err != nil {
			fmt.Printf("Error decrypting text: %v\n", err)
			os.Exit(1)
		}

		if decryptOutput != "" {
			err = utils.WriteFile(decryptOutput, result)
			if err != nil {
				fmt.Printf("Error writing to output file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Decrypted text written to: %s\n", decryptOutput)
		} else {
			fmt.Printf("Decrypted text: %s\n", string(result))
		}
	} else {
		// Decrypt file
		data, err := utils.ReadFile(decryptFile)
		if err != nil {
			fmt.Printf("Error reading encrypted file: %v\n", err)
			os.Exit(1)
		}

		result, err = provider.Decrypt(data, keyBytes)
		if err != nil {
			fmt.Printf("Error decrypting file: %v\n", err)
			os.Exit(1)
		}

		outputFile := decryptOutput
		if outputFile == "" {
			if len(decryptFile) > 10 && decryptFile[len(decryptFile)-10:] == ".encrypted" {
				outputFile = decryptFile[:len(decryptFile)-10]
			} else {
				outputFile = decryptFile + ".decrypted"
			}
		}

		err = utils.WriteFile(outputFile, result)
		if err != nil {
			fmt.Printf("Error writing decrypted file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("File decrypted and saved to: %s\n", outputFile)
	}
}