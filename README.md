# Cross-Platform Encryption/Decryption Tool

A Go application that supports encryption and decryption using various algorithms like AES-256-CBC and RSA. Supports both text and file encryption/decryption across macOS, Windows, and Linux.

## Features

- **Multiple Algorithms**: Support for AES-256-CBC and RSA encryption
- **Cross-Platform**: Runs on macOS, Windows, and Linux (x64 and ARM64)
- **Text & File Support**: Encrypt/decrypt both text strings and files
- **Base64 Key Support**: Input keys in base64 format (automatically converted)
- **Environment Variable Keys**: Support for reading keys from environment variables
- **Extensible Design**: Easy to add new encryption algorithms

## Installation

### Option 1: Automated Installation (Recommended)

#### Linux/macOS
```bash
# Install latest version
curl -sSL https://raw.githubusercontent.com/thanhlv-com/thanhlv-encryption-decryption/main/install.sh | bash

# Install specific version
curl -sSL https://raw.githubusercontent.com/thanhlv-com/thanhlv-encryption-decryption/main/install.sh | bash -s -- v1.1.0
```

#### Windows (PowerShell as Administrator)
```powershell
# Install latest version
PowerShell -ExecutionPolicy Bypass -Command "iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/thanhlv-com/thanhlv-encryption-decryption/main/install.ps1'))"

# Install specific version
$env:VERSION='v1.2.0'; PowerShell -ExecutionPolicy Bypass -Command "iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/thanhlv-com/thanhlv-encryption-decryption/main/install.ps1'))"
```

The automated installer will:
- Detect your platform (OS and architecture)
- Download the appropriate binary from GitHub releases
- Install to `/usr/local/bin` (Linux/macOS) or `Program Files` (Windows)
- Add to system PATH
- Verify the installation

### Option 2: Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd thanhlv-encryption-decryption

# Build for current platform
make build

# build-all for all platforms
make build-all
```

### Option 3: Download Pre-built Binaries

Download the appropriate binary for your platform from the [GitHub releases](https://github.com/thanhlv-com/thanhlv-encryption-decryption/releases) page.

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
# Encrypt text with key flag
./thanhlv-ed encrypt -a aes-256-cbc -k "<base64-text-key>" -t "Hello World!"

# Encrypt text with environment variable
export ENCRYPTION_KEY="<base64-text-key>"
./thanhlv-ed encrypt -a aes-256-cbc -e ENCRYPTION_KEY -t "Hello World!"

# Decrypt text with key flag
./thanhlv-ed decrypt -a aes-256-cbc -k "<base64-text-key>" -t "<base64-encrypted-text>"

# Decrypt text with environment variable
export ENCRYPTION_KEY="<base64-text-key>"
./thanhlv-ed decrypt -a aes-256-cbc -e ENCRYPTION_KEY -t "<base64-encrypted-text>"
```

##### AES-256-CBC demo

```bash
# Encrypt text
./thanhlv-ed encrypt -a aes-256-cbc -k "MTIzZGY=" -t "Hello World!"

# Decrypt text
./thanhlv-ed decrypt -a aes-256-cbc -k "MTIzZGY=" -t "<base64-encrypted-text>"
```

#### RSA

```bash
# Encrypt with public key (key flag)
./thanhlv-ed encrypt -a rsa -k "<base64-public-key>" -t "Hello RSA World!"

# Encrypt with public key (environment variable)
export RSA_PUBLIC_KEY="<base64-public-key>"
./thanhlv-ed encrypt -a rsa -e RSA_PUBLIC_KEY -t "Hello RSA World!"

# Decrypt with private key (key flag)
./thanhlv-ed decrypt -a rsa -k "<base64-private-key>" -t "<base64-encrypted-text>"

# Decrypt with private key (environment variable)
export RSA_PRIVATE_KEY="<base64-private-key>"
./thanhlv-ed decrypt -a rsa -e RSA_PRIVATE_KEY -t "<base64-encrypted-text>"
```

### File Encryption/Decryption

#### AES-256-CBC

```bash
# Encrypt file with key flag
./thanhlv-ed encrypt -a aes-256-cbc -k "MTIzZGY=" -f input.txt

# Encrypt file with environment variable
export ENCRYPTION_KEY="MTIzZGY="
./thanhlv-ed encrypt -a aes-256-cbc -e ENCRYPTION_KEY -f input.txt

# Decrypt file with key flag
./thanhlv-ed decrypt -a aes-256-cbc -k "MTIzZGY=" -f input.txt.encrypted

# Decrypt file with environment variable
export ENCRYPTION_KEY="MTIzZGY="
./thanhlv-ed decrypt -a aes-256-cbc -e ENCRYPTION_KEY -f input.txt.encrypted
```

#### RSA

```bash
# Encrypt file with public key (key flag)
./thanhlv-ed encrypt -a rsa -k "<base64-public-key>" -f document.pdf

# Encrypt file with public key (environment variable)
export RSA_PUBLIC_KEY="<base64-public-key>"
./thanhlv-ed encrypt -a rsa -e RSA_PUBLIC_KEY -f document.pdf

# Decrypt file with private key (key flag)
./thanhlv-ed decrypt -a rsa -k "<base64-private-key>" -f document.pdf.encrypted

# Decrypt file with private key (environment variable)
export RSA_PRIVATE_KEY="<base64-private-key>"
./thanhlv-ed decrypt -a rsa -e RSA_PRIVATE_KEY -f document.pdf.encrypted
```

### Command Options

#### Common Flags

- `-a, --algorithm`: Encryption algorithm (`aes-256-cbc`, `rsa`)
- `-k, --key`: Encryption/decryption key (base64 encoded)
- `-e, --key-env`: Environment variable name containing the key (base64 encoded)
- `-t, --text`: Text to encrypt/decrypt
- `-f, --file`: File to encrypt/decrypt
- `-o, --output`: Output file (optional)

**Note**: Either `--key` or `--key-env` must be specified (but not both).

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

# 2. Encrypt text using key flag
./thanhlv-ed encrypt -a aes-256-cbc -k "MTExZHZzZHZzeGN2eGN2" -t "Secret message"
# Output: Encrypted text (base64): <encrypted-data>

# 2. Alternative: Encrypt text using environment variable
export MY_AES_KEY="MTExZHZzZHZzeGN2eGN2"
./thanhlv-ed encrypt -a aes-256-cbc -e MY_AES_KEY -t "Secret message"
# Output: Encrypted text (base64): <encrypted-data>

# 3. Decrypt text using key flag
./thanhlv-ed decrypt -a aes-256-cbc -k "MTExZHZzZHZzeGN2eGN2" -t "J+FWPLdwO+N6BSaRo2o8vCUImk50kHYi4SkHrzeLP9Q="
# Output: Decrypted text: Secret message

# 3. Alternative: Decrypt text using environment variable
export MY_AES_KEY="MTExZHZzZHZzeGN2eGN2"
./thanhlv-ed decrypt -a aes-256-cbc -e MY_AES_KEY -t "J+FWPLdwO+N6BSaRo2o8vCUImk50kHYi4SkHrzeLP9Q="
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
