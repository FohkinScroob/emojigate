#!/usr/bin/env bash
set -e

# Detect if emojigate is already installed
if ! command -v emojigate &> /dev/null; then
    # Install emojigate if not found
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

    # Try to get version from .version file
    VERSION="latest"
    if [ -f "$SCRIPT_DIR/../.version" ]; then
        VERSION=$(cat "$SCRIPT_DIR/../.version")
    fi

    "$SCRIPT_DIR/install.sh" "$VERSION" "$HOME/.local/bin"
    export PATH="$PATH:$HOME/.local/bin"
fi

# Run emojigate with all arguments
emojigate lint "$@"
