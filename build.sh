#!/bin/bash

# Build script for cross-platform compilation
APP_NAME="thanhlv-ed"
VERSION=${1:-"1.0.0"}

echo "Building $APP_NAME version $VERSION for multiple platforms..."

# Clean up previous builds
rm -rf build/
mkdir -p build

# Build for different platforms
platforms=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    output_name="build/${APP_NAME}-${GOOS}-${GOARCH}"
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    echo "Building for $GOOS/$GOARCH..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o $output_name main.go

    if [ $? -ne 0 ]; then
        echo "An error has occurred! Aborting the script execution..."
        exit 1
    fi
done

echo "Build completed successfully!"
echo "Binaries are available in the build/ directory:"
ls -la build/