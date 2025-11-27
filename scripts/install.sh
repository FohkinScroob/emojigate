#!/bin/bash
set -e

VERSION="${1:-latest}"
INSTALL_DIR="${2:-$HOME/.local/bin}"

REPO="FohkinScroob/emojigate"
BINARY_NAME="emojigate"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

case "$OS" in
    darwin)
        OS="darwin"
        ;;
    linux)
        OS="linux"
        ;;
    mingw*|msys*|cygwin*)
        OS="windows"
        BINARY_NAME="emojigate.exe"
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

# Get latest version if not specified
if [ "$VERSION" = "latest" ]; then
    VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        echo "Failed to get latest version"
        exit 1
    fi
fi

# Remove 'v' prefix if present
VERSION=${VERSION#v}

# Build download URL
DOWNLOAD_URL="https://github.com/$REPO/releases/download/v${VERSION}/${BINARY_NAME}-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    DOWNLOAD_URL="${DOWNLOAD_URL}.exe"
fi

echo "Downloading emojigate v${VERSION} for ${OS}/${ARCH}..."
echo "URL: $DOWNLOAD_URL"

# Create install directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Download binary
TEMP_FILE=$(mktemp)
if command -v curl &> /dev/null; then
    curl -sL "$DOWNLOAD_URL" -o "$TEMP_FILE"
elif command -v wget &> /dev/null; then
    wget -q "$DOWNLOAD_URL" -O "$TEMP_FILE"
else
    echo "Error: curl or wget is required"
    exit 1
fi

# Move to install directory
mv "$TEMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo "✅ Successfully installed emojigate to $INSTALL_DIR/$BINARY_NAME"
echo ""
echo "Make sure $INSTALL_DIR is in your PATH:"
echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
