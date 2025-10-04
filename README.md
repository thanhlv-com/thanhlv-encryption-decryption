# Cross-Platform Encryption/Decryption Tool

A Go application that supports encryption and decryption using various algorithms like AES-256-CBC and RSA. Supports both text and file encryption/decryption across macOS, Windows, and Linux.

## Features

- **Multiple Algorithms**: Support for AES-256-CBC and RSA encryption
- **Cross-Platform**: Runs on macOS, Windows, and Linux (x64 and ARM64)
- **Text & File Support**: Encrypt/decrypt both text strings and files
- **Base64 Key Support**: Input keys in base64 format (automatically converted)
- **Extensible Design**: Easy to add new encryption algorithms

## Installation

### Option 1: Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd thanhlv-encryption-decryption

# Build for current platform
make build

# Cross-compile for all platforms
make cross-compile
```

### Option 2: Download Pre-built Binaries

Download the appropriate binary for your platform from the `build/` directory after running `make cross-compile`.

## Usage

### Key Generation

Generate encryption keys for different algorithms:

```bash
# Generate AES-256-CBC key (outputs base64)
./thanhlv-ed keygen -a aes-256-cbc -b

# Generate RSA key pair
./thanhlv-ed keygen -a rsa -b
```

### Text Encryption/Decryption

#### AES-256-CBC

```bash
# Encrypt text
./thanhlv-ed encrypt -a aes-256-cbc -k "MTIzZGY=" -t "Hello World!"

# Decrypt text
./thanhlv-ed decrypt -a aes-256-cbc -k "MTIzZGY=" -t "<base64-encrypted-text>"
```

#### RSA

```bash
# Encrypt with public key
./thanhlv-ed encrypt -a rsa -k "<base64-public-key>" -t "Hello RSA World!"

# Decrypt with private key
./thanhlv-ed decrypt -a rsa -k "<base64-private-key>" -t "<base64-encrypted-text>"
```

### File Encryption/Decryption

#### AES-256-CBC

```bash
# Encrypt file
./thanhlv-ed encrypt -a aes-256-cbc -k "MTIzZGY=" -f input.txt

# Decrypt file
./thanhlv-ed decrypt -a aes-256-cbc -k "MTIzZGY=" -f input.txt.encrypted
```

#### RSA

```bash
# Encrypt file with public key
./thanhlv-ed encrypt -a rsa -k "<base64-public-key>" -f document.pdf

# Decrypt file with private key
./thanhlv-ed decrypt -a rsa -k "<base64-private-key>" -f document.pdf.encrypted
```

### Command Options

#### Common Flags

- `-a, --algorithm`: Encryption algorithm (`aes-256-cbc`, `rsa`)
- `-k, --key`: Encryption/decryption key (base64 encoded)
- `-t, --text`: Text to encrypt/decrypt
- `-f, --file`: File to encrypt/decrypt
- `-o, --output`: Output file (optional)

#### Key Generation Flags

- `-b, --base64`: Output key in base64 format
- `-p, --private`: Private key output file (RSA only)
- `-u, --public`: Public key output file (RSA only)

## Key Format Examples

### AES-256-CBC Key

The application accepts base64-encoded keys. For example, if you have the string `123df`, encode it to base64:

```bash
echo -n "123df" | base64
# Output: MTIzZGY=
```

Use `MTIzZGY=` as your key parameter.

### RSA Keys

RSA keys are generated in PEM format and can be used directly in base64 encoded form:

```bash
# Generate RSA keys
./thanhlv-ed keygen -a rsa -b

# This outputs both private and public keys in base64 format
```

## Algorithm Details

### AES-256-CBC

- **Key Size**: 256-bit (automatically derived from input using SHA-256)
- **Block Size**: 128-bit
- **Padding**: PKCS#7
- **IV**: Randomly generated for each encryption

### RSA

- **Key Size**: 2048-bit
- **Padding**: OAEP with SHA-256
- **Format**: PEM (PKCS#1 for private keys, PKIX for public keys)
- **Chunking**: Automatically handles large data by splitting into chunks

## Examples

### Complete AES Workflow

```bash
# 1. Generate a key
./thanhlv-ed keygen -a aes-256-cbc -b
# Output: Generated AES-256-CBC key (base64): <your-key>

# 2. Encrypt text
./thanhlv-ed encrypt -a aes-256-cbc -k "MTExZHZzZHZzeGN2eGN2" -t "Secret message"
# Output: Encrypted text (base64): <encrypted-data>

# 3. Decrypt text
./thanhlv-ed decrypt -a aes-256-cbc -k "MTExZHZzZHZzeGN2eGN2" -t "J+FWPLdwO+N6BSaRo2o8vCUImk50kHYi4SkHrzeLP9Q="
# Output: Decrypted text: Secret message
```

### Complete RSA Workflow

```bash
# 1. Generate RSA key pair
./thanhlv-ed keygen -a rsa -b
# Outputs private and public keys in base64

# 2. Encrypt with public key
./thanhlv-ed encrypt -a rsa -k "<public-key-base64>" -t "Confidential data"

# 3. Decrypt with private key
./thanhlv-ed decrypt -a rsa -k "<private-key-base64>" -t "<encrypted-data>"
```

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional, for convenience)

### Building

```bash
# Install dependencies
make deps

# Build for development
make dev

# Run tests
make test

# Clean build artifacts
make clean
```

### Adding New Algorithms

1. Implement the `CryptoProvider` interface in `pkg/crypto/`
2. Add the algorithm to the factory function in `pkg/crypto/init.go`
3. Update the CLI help text and documentation

## Platform Support

- **macOS**: Intel (x64) and Apple Silicon (ARM64)
- **Linux**: x64 and ARM64
- **Windows**: x64

## Security Notes

- AES keys are derived using SHA-256 for consistent 256-bit length
- RSA uses OAEP padding with SHA-256 for security
- Random IVs are generated for each AES encryption
- Private keys should be kept secure and never shared
- This tool is for educational and development purposes

## License

[Add your license information here]