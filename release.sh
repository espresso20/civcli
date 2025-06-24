#!/bin/bash

# Create a release package for CivIdleCli

VERSION="1.0.1"
PACKAGE_NAME="cividlecli-${VERSION}"

echo "Creating release package for CivIdleCli v${VERSION}..."

# Create release directory
mkdir -p "releases/${PACKAGE_NAME}"

# Build for different platforms
echo "Building for multiple platforms..."
GOOS=darwin GOARCH=amd64 go build -o "releases/${PACKAGE_NAME}/cividlecli-darwin-amd64" main.go
GOOS=darwin GOARCH=arm64 go build -o "releases/${PACKAGE_NAME}/cividlecli-darwin-arm64" main.go
GOOS=linux GOARCH=amd64 go build -o "releases/${PACKAGE_NAME}/cividlecli-linux-amd64" main.go
GOOS=windows GOARCH=amd64 go build -o "releases/${PACKAGE_NAME}/cividlecli-windows-amd64.exe" main.go

# Copy README and other files
cp README.md "releases/${PACKAGE_NAME}/"
mkdir -p "releases/${PACKAGE_NAME}/data/saves"

# Create zip archives
echo "Creating zip archives..."
cd releases
zip -r "${PACKAGE_NAME}-macos.zip" "${PACKAGE_NAME}/cividlecli-darwin-amd64" "${PACKAGE_NAME}/cividlecli-darwin-arm64" "${PACKAGE_NAME}/README.md" "${PACKAGE_NAME}/data"
zip -r "${PACKAGE_NAME}-linux.zip" "${PACKAGE_NAME}/cividlecli-linux-amd64" "${PACKAGE_NAME}/README.md" "${PACKAGE_NAME}/data"
zip -r "${PACKAGE_NAME}-windows.zip" "${PACKAGE_NAME}/cividlecli-windows-amd64.exe" "${PACKAGE_NAME}/README.md" "${PACKAGE_NAME}/data"

echo "Release packages created in releases/ directory!"
