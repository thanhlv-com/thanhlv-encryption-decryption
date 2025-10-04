# thanhlv-encryption-decryption Auto Install Script for Windows
# Usage:
# PowerShell -ExecutionPolicy Bypass -Command "iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/thanhlv-com/thanhlv-encryption-decryption/main/install.ps1'))"
# Or with version:
# PowerShell -ExecutionPolicy Bypass -Command "$env:VERSION='v1.2.0'; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/thanhlv-com/thanhlv-encryption-decryption/main/install.ps1'))"

param(
    [string]$Version = $env:VERSION
)

# Configuration
$GitHubRepo = "thanhlv-com/thanhlv-encryption-decryption"
$AppName = "thanhlv-ed"
$DefaultVersion = "v1"
$InstallDir = "${env:ProgramFiles}\${AppName}"

# Colors for output (Windows PowerShell compatible)
function Write-ColorOutput {
    param(
        [string]$Message,
        [string]$Level = "INFO"
    )

    switch ($Level) {
        "INFO" { Write-Host "[INFO] $Message" -ForegroundColor Blue }
        "SUCCESS" { Write-Host "[SUCCESS] $Message" -ForegroundColor Green }
        "WARNING" { Write-Host "[WARNING] $Message" -ForegroundColor Yellow }
        "ERROR" {
            Write-Host "[ERROR] $Message" -ForegroundColor Red
            exit 1
        }
    }
}

# Check if running as administrator
function Test-Administrator {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

# Detect architecture
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    if ($arch -eq "AMD64" -or $arch -eq "x64") {
        return "amd64"
    }
    elseif ($arch -eq "ARM64") {
        return "arm64"
    }
    else {
        Write-ColorOutput "Unsupported architecture: $arch. Supported architectures: AMD64, ARM64" "ERROR"
    }
}

# Check prerequisites
function Test-Prerequisites {
    Write-ColorOutput "Checking prerequisites..."

    # Check PowerShell version
    if ($PSVersionTable.PSVersion.Major -lt 3) {
        Write-ColorOutput "PowerShell 3.0 or higher is required" "ERROR"
    }

    # Check .NET Framework for web requests
    try {
        [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType]::Tls12
    }
    catch {
        Write-ColorOutput "Failed to set TLS 1.2. Please ensure .NET Framework 4.5+ is installed" "ERROR"
    }

    Write-ColorOutput "Prerequisites check passed" "SUCCESS"
}

# Get latest release version from GitHub API
function Get-LatestVersion {
    param([string]$SpecifiedVersion)

    if ($SpecifiedVersion) {
        $script:Version = $SpecifiedVersion
        Write-ColorOutput "Using specified version: $script:Version"
        return
    }

    Write-ColorOutput "Fetching latest release version..."

    try {
        $apiUrl = "https://api.github.com/repos/$GitHubRepo/releases/latest"
        $response = Invoke-RestMethod -Uri $apiUrl -Method Get
        $script:Version = $response.tag_name
        Write-ColorOutput "Latest version: $script:Version"
    }
    catch {
        Write-ColorOutput "Could not fetch latest version, using default: $DefaultVersion" "WARNING"
        $script:Version = $DefaultVersion
    }
}

# Download and install binary
function Install-Binary {
    $arch = Get-Architecture
    $binaryName = "${AppName}-windows-${arch}.exe"
    $downloadUrl = "https://github.com/$GitHubRepo/releases/download/$script:Version/$binaryName"
    $tempPath = Join-Path $env:TEMP $binaryName
    $finalPath = Join-Path $InstallDir "${AppName}.exe"

    Write-ColorOutput "Downloading $AppName $script:Version for Windows $arch..."
    Write-ColorOutput "Download URL: $downloadUrl"

    try {
        # Download binary
        Write-ColorOutput "Downloading binary..."
        [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType]::Tls12
        $webClient = New-Object System.Net.WebClient
        $webClient.DownloadFile($downloadUrl, $tempPath)

        # Create install directory
        if (!(Test-Path $InstallDir)) {
            Write-ColorOutput "Creating installation directory: $InstallDir"
            New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        }

        # Move binary to install location
        Write-ColorOutput "Installing to $InstallDir..."
        Move-Item -Path $tempPath -Destination $finalPath -Force

        Write-ColorOutput "$AppName installed successfully to $finalPath" "SUCCESS"
    }
    catch {
        Write-ColorOutput "Failed to download or install binary: $($_.Exception.Message)" "ERROR"
    }
}

# Add to PATH
function Add-ToPath {
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "Machine")

    if ($currentPath -notlike "*$InstallDir*") {
        Write-ColorOutput "Adding $InstallDir to system PATH..."

        try {
            $newPath = "$currentPath;$InstallDir"
            [Environment]::SetEnvironmentVariable("PATH", $newPath, "Machine")
            Write-ColorOutput "Added to system PATH successfully" "SUCCESS"
            Write-ColorOutput "Please restart your command prompt or PowerShell to use the updated PATH" "WARNING"
        }
        catch {
            Write-ColorOutput "Failed to add to system PATH. You may need to add $InstallDir manually" "WARNING"
        }
    }
    else {
        Write-ColorOutput "Installation directory already in PATH"
    }
}

# Verify installation
function Test-Installation {
    Write-ColorOutput "Verifying installation..."

    $binaryPath = Join-Path $InstallDir "${AppName}.exe"

    if (Test-Path $binaryPath) {
        try {
            $version = & $binaryPath --version 2>$null
            if ($LASTEXITCODE -eq 0) {
                Write-ColorOutput "Installation verified! Run '$AppName --help' to get started." "SUCCESS"
                Write-ColorOutput "Installed version: $version"
            }
            else {
                Write-ColorOutput "Binary installed but version check failed" "WARNING"
            }
        }
        catch {
            Write-ColorOutput "Binary installed at $binaryPath" "SUCCESS"
            Write-ColorOutput "You can run: `$binaryPath --help"
        }
    }
    else {
        Write-ColorOutput "Installation verification failed - binary not found at $binaryPath" "ERROR"
    }
}

# Show usage information
function Show-Usage {
    Write-ColorOutput "Usage examples:"
    Write-Host "  # Generate AES key:" -ForegroundColor Cyan
    Write-Host "  $AppName keygen -a aes-256-cbc -b"
    Write-Host ""
    Write-Host "  # Encrypt text:" -ForegroundColor Cyan
    Write-Host "  $AppName encrypt -a aes-256-cbc -k `"MTIzZGY=`" -t `"Hello World!`""
    Write-Host ""
    Write-Host "  # Decrypt text:" -ForegroundColor Cyan
    Write-Host "  $AppName decrypt -a aes-256-cbc -k `"MTIzZGY=`" -t `"<encrypted-text>`""
    Write-Host ""
    Write-Host "  # Get help:" -ForegroundColor Cyan
    Write-Host "  $AppName --help"
}

# Main installation flow
function Main {
    Write-Host ""
    Write-ColorOutput "Starting $AppName installation..."
    Write-Host ""

    if (!(Test-Administrator)) {
        Write-ColorOutput "This script requires administrator privileges to install to Program Files and modify PATH" "ERROR"
    }

    Test-Prerequisites
    Get-LatestVersion $Version
    Install-Binary
    Add-ToPath
    Test-Installation

    Write-Host ""
    Write-ColorOutput "Installation completed!" "SUCCESS"
    Write-Host ""
    Show-Usage
    Write-Host ""
}

# Handle script interruption
trap {
    Write-ColorOutput "Installation interrupted" "ERROR"
}

# Run main function
Main