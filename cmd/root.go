package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"thanhlv-encryption-decryption/pkg/utils"
)

var (
	debugFlag bool
	rootCmd = &cobra.Command{
		Use:   "thanhlv-ed",
		Short: "A cross-platform encryption/decryption tool",
		Long: `A Go application that supports encryption and decryption using various algorithms
like AES-256-CBC and RSA. Supports both text and file encryption/decryption.`,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "enable debug logging")
	utils.SetDebugEnabledFunc(IsDebugEnabled)
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(keygenCmd)
}

func IsDebugEnabled() bool {
	return debugFlag
}