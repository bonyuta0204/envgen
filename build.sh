#!/bin/bash
set -e

# Build for multiple platforms
build() {
    local output="$1"
    local os="$2"
    local arch="$3"
    
    echo "Building for $os/$arch..."
    GOOS=$os GOARCH=$arch go build -o "$output" .
}

# Create build directory
mkdir -p build

# Build for common platforms
build "build/envgen-darwin-amd64" "darwin" "amd64"
build "build/envgen-darwin-arm64" "darwin" "arm64"
build "build/envgen-linux-amd64" "linux" "amd64"
build "build/envgen-linux-arm64" "linux" "arm64"
build "build/envgen-windows-amd64.exe" "windows" "amd64"

echo "Build complete. Binaries are in the build directory."
