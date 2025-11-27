#!/usr/bin/env bash
set -e

# Detect if emojigate is already installed
if ! command -v emojigate &> /dev/null; then
    # Install emojigate if not found
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

    # Try to get version from git tag in the pre-commit repo
    VERSION="latest"
    if [ -d "$SCRIPT_DIR/../.git" ]; then
        GIT_TAG=$(cd "$SCRIPT_DIR/.." && git describe --tags --exact-match 2>/dev/null || echo "")
        if [ -n "$GIT_TAG" ]; then
            VERSION="${GIT_TAG#v}"
        fi
    fi

    "$SCRIPT_DIR/install.sh" "$VERSION" "$HOME/.local/bin"
    export PATH="$PATH:$HOME/.local/bin"
fi

# Run emojigate with all arguments
emojigate lint "$@"
