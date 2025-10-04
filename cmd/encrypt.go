package cmd

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"thanhlv-encryption-decryption/pkg/crypto"
	"thanhlv-encryption-decryption/pkg/utils"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt text or files",
	Long:  `Encrypt text or files using various algorithms like AES-256-CBC or RSA.`,
	Run:   runEncrypt,
}

var (
	encryptAlgorithm string
	encryptInput     string
	encryptOutput    string
	encryptKey       string
	encryptKeyEnv    string
	encryptText      string
	encryptFile      string
)

func init() {
	encryptCmd.Flags().StringVarP(&encryptAlgorithm, "algorithm", "a", "aes-256-cbc", "Encryption algorithm (aes-256-cbc, rsa)")
	encryptCmd.Flags().StringVarP(&encryptKey, "key", "k", "", "Encryption key (base64 encoded)")
	encryptCmd.Flags().StringVarP(&encryptKeyEnv, "key-env", "e", "", "Environment variable name containing the encryption key (base64 encoded)")
	encryptCmd.Flags().StringVarP(&encryptText, "text", "t", "", "Text to encrypt")
	encryptCmd.Flags().StringVarP(&encryptFile, "file", "f", "", "File to encrypt")
	encryptCmd.Flags().StringVarP(&encryptOutput, "output", "o", "", "Output file (optional)")
}

func runEncrypt(cmd *cobra.Command, args []string) {
	utils.DebugLogf("Starting encryption with algorithm: %s", encryptAlgorithm)
	utils.DebugLogf("Text input: %t, File input: %t", encryptText != "", encryptFile != "")

	if encryptText == "" && encryptFile == "" {
		fmt.Println("Error: Either --text or --file must be specified")
		os.Exit(1)
	}

	if encryptText != "" && encryptFile != "" {
		fmt.Println("Error: Cannot specify both --text and --file")
		os.Exit(1)
	}

	// Validate key input
	if encryptKey == "" && encryptKeyEnv == "" {
		fmt.Println("Error: Either --key or --key-env must be specified")
		os.Exit(1)
	}

	if encryptKey != "" && encryptKeyEnv != "" {
		fmt.Println("Error: Cannot specify both --key and --key-env")
		os.Exit(1)
	}

	// Get the key value
	var keyValue string
	if encryptKeyEnv != "" {
		keyValue = os.Getenv(encryptKeyEnv)
		if keyValue == "" {
			fmt.Printf("Error: Environment variable '%s' is not set or empty\n", encryptKeyEnv)
			os.Exit(1)
		}
		utils.DebugLogf("Using key from environment variable: %s", encryptKeyEnv)
	} else {
		keyValue = encryptKey
		utils.DebugLogf("Using key from command line flag")
	}

	// Decode base64 key
	utils.DebugLogf("Decoding base64 key of length: %d", len(keyValue))
	keyBytes, err := base64.StdEncoding.DecodeString(keyValue)
	if err != nil {
		fmt.Printf("Error decoding base64 key: %v\n", err)
		os.Exit(1)
	}
	utils.DebugLogf("Successfully decoded key, byte length: %d", len(keyBytes))

	// Initialize crypto provider
	utils.DebugLogf("Initializing crypto provider for algorithm: %s", encryptAlgorithm)
	provider, err := crypto.NewCryptoProvider(encryptAlgorithm)
	if err != nil {
		fmt.Printf("Error initializing crypto provider: %v\n", err)
		os.Exit(1)
	}
	utils.DebugLog("Crypto provider initialized successfully")

	var result []byte

	if encryptText != "" {
		// Encrypt text
		utils.DebugLogf("Encrypting text of length: %d", len(encryptText))
		result, err = provider.Encrypt([]byte(encryptText), keyBytes)
		if err != nil {
			fmt.Printf("Error encrypting text: %v\n", err)
			os.Exit(1)
		}

		if encryptOutput != "" {
			err = utils.WriteFile(encryptOutput, result)
			if err != nil {
				fmt.Printf("Error writing to output file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Encrypted text written to: %s\n", encryptOutput)
		} else {
			fmt.Printf("Encrypted text (base64): %s\n", base64.StdEncoding.EncodeToString(result))
		}
	} else {
		// Encrypt file
		utils.DebugLogf("Reading file for encryption: %s", encryptFile)
		data, err := utils.ReadFile(encryptFile)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		utils.DebugLogf("File read successfully, size: %d bytes", len(data))
		utils.DebugLog("Starting file encryption")
		result, err = provider.Encrypt(data, keyBytes)
		if err != nil {
			fmt.Printf("Error encrypting file: %v\n", err)
			os.Exit(1)
		}
		utils.DebugLogf("File encrypted successfully, result size: %d bytes", len(result))

		outputFile := encryptOutput
		if outputFile == "" {
			outputFile = encryptFile + ".encrypted"
		}

		err = utils.WriteFile(outputFile, result)
		if err != nil {
			fmt.Printf("Error writing encrypted file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("File encrypted and saved to: %s\n", outputFile)
	}
}