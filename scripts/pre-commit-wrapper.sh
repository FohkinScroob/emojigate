#!/bin/bash
set -e

# Detect if emojigate is already installed
if ! command -v emojigate &> /dev/null; then
    # Install emojigate if not found
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    "$SCRIPT_DIR/install.sh" latest "$HOME/.local/bin"
    export PATH="$PATH:$HOME/.local/bin"
fi

# Run emojigate with all arguments
emojigate lint "$@"
