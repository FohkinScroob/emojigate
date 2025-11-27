# 🎯 emojigate

**Enforce emoji usage in GitHub Actions workflows** - because workflows deserve to be expressive!

Emojigate lints your GitHub Actions YAML files to ensure every workflow, job, and step name starts with an emoji. Perfect for teams who want consistent, visually appealing workflows.

## ✨ Features

- 🔍 Lints workflow names, job names, and step names
- 📁 Supports multiple files and directory scanning
- 🪝 Pre-commit hook integration
- 🚀 Zero dependencies for CLI usage
- ⚡ Fast and lightweight

## 🚀 Installation

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/FohkinScroob/emojigate/releases):

1. Download the binary for your platform (e.g., `emojigate-linux-amd64`)
2. Rename it to `emojigate` (or `emojigate.exe` on Windows)
3. Make it executable: `chmod +x emojigate` (Linux/macOS)
4. Move it to a directory in your PATH (e.g., `/usr/local/bin` or `~/.local/bin`)

**Quick install script:**
```bash
curl -sSL https://raw.githubusercontent.com/FohkinScroob/emojigate/main/scripts/install.sh | bash
```

### Build from Source

```bash
go install github.com/FohkinScroob/emojigate/cmd/emojigate@latest
```

Or clone and build:

```bash
git clone https://github.com/FohkinScroob/emojigate.git
cd emojigate
make build
```

## 📖 Usage

### Lint all workflows

```bash
emojigate workflows
```

Automatically lints all `.yml` and `.yaml` files in `.github/workflows/`.

### Lint specific files

```bash
emojigate lint .github/workflows/ci.yml
emojigate lint .github/workflows/ci.yml .github/workflows/release.yml
```

### Get help

```bash
emojigate help
```

## 🪝 Pre-commit Integration

### Option 1: Auto-download binary (recommended, no Go required)

Add to your `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: https://github.com/FohkinScroob/emojigate
    rev: v1.0.0
    hooks:
      - id: emojigate-binary
```

Pre-commit will automatically download the appropriate binary for your platform from GitHub releases.

### Option 2: Auto-build with Go (requires Go installed)

```yaml
repos:
  - repo: https://github.com/FohkinScroob/emojigate
    rev: v1.0.0
    hooks:
      - id: emojigate
```

Pre-commit will compile the binary from source using your local Go installation.

## ✅ Example

**❌ Invalid workflow:**

```yaml
name: CI Pipeline
jobs:
  build:
    name: Build Application
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
```

**✅ Valid workflow:**

```yaml
name: 🚀 CI Pipeline
jobs:
  build:
    name: 🔨 Build Application
    steps:
      - name: 📥 Checkout code
        uses: actions/checkout@v3
```

## 🛠️ Development

### Prerequisites

- Go 1.25 or higher

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Test Coverage

```bash
make test-coverage
```

### Project Structure

```
emojigate/
├── cmd/emojigate/     # CLI entry point
├── internal/          # Core linting logic
│   ├── linter.go      # Workflow linter
│   ├── parser.go      # YAML parser
│   └── testdata/      # Test fixtures
├── Makefile           # Build tasks
└── README.md
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m '✨ Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 💡 Why Emojis?

Emojis make workflows easier to scan and understand at a glance. They provide visual anchors that help you quickly identify different parts of your CI/CD pipeline.

- 🚀 Deployments
- 🧪 Tests
- 🔨 Builds
- 📦 Releases
- 🔒 Security scans
- And many more!
