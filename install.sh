#!/usr/bin/env bash

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}  Claudes - Claude CLI Profile Manager${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed.${NC}"
    echo ""
    echo "Please install Go first:"
    echo "  macOS:   brew install go"
    echo "  Linux:   sudo apt install golang-go"
    echo "  Windows: Download from https://go.dev/dl/"
    echo ""
    exit 1
fi

# Check Go version (at least 1.16)
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
GO_MAJOR=$(echo "$GO_VERSION" | cut -d. -f1)
GO_MINOR=$(echo "$GO_VERSION" | cut -d. -f2)

if [ "$GO_MAJOR" -lt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 16 ]); then
    echo -e "${RED}Error: Go 1.16+ is required. You have Go $GO_VERSION.${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Go $GO_VERSION found${NC}"

# Create temporary directory
TMP_DIR=$(mktemp -d)
echo -e "${GREEN}✓ Created temporary directory: $TMP_DIR${NC}"

# Cleanup on exit
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

# Clone the repository
echo ""
echo -e "${YELLOW}Cloning repository...${NC}"
cd "$TMP_DIR"
git clone --branch develop https://github.com/yasaricli/claudes.git
cd claudes

# Build the binary
echo ""
echo -e "${YELLOW}Building claudes...${NC}"
go build -o claudes .

# Determine installation directory
if [ "$EUID" -eq 0 ]; then
    # Running as root, install to /usr/local/bin
    INSTALL_DIR="/usr/local/bin"
else
    # Check if ~/.local/bin exists and is in PATH
    if [ -d "$HOME/.local/bin" ] && [[ ":$PATH:" == *":$HOME/.local/bin:"* ]]; then
        INSTALL_DIR="$HOME/.local/bin"
    else
        # Default to /usr/local/bin (may need sudo)
        INSTALL_DIR="/usr/local/bin"
    fi
fi

echo ""
echo -e "${YELLOW}Installing to $INSTALL_DIR...${NC}"

# Install the binary
if [ "$INSTALL_DIR" = "/usr/local/bin" ] && [ "$EUID" -ne 0 ]; then
    # Need sudo for /usr/local/bin
    if command -v sudo &> /dev/null; then
        echo -e "${YELLOW}Using sudo for installation to /usr/local/bin${NC}"
        sudo mv claudes "$INSTALL_DIR/"
    else
        echo -e "${RED}Error: Cannot install to /usr/local/bin without root privileges.${NC}"
        echo ""
        echo "Please run this script with sudo, or install Go and build manually:"
        echo "  git clone https://github.com/yasaricli/claudes.git"
        echo "  cd claudes"
        echo "  go build -o claudes ."
        echo "  mv claudes ~/.local/bin/  # or somewhere in your PATH"
        echo ""
        exit 1
    fi
else
    mv claudes "$INSTALL_DIR/"
fi

# Verify installation
echo ""
echo -e "${YELLOW}Verifying installation...${NC}"
if command -v claudes &> /dev/null; then
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  Installation successful!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo "Installed at: $(which claudes)"
    echo ""
    echo "Quick start:"
    echo "  claudes help          # Show available commands"
    echo "  claudes add           # Create a new profile"
    echo "  claudes list          # List all profiles"
    echo "  claudes <profile>     # Run Claude with a profile"
    echo ""
    exit 0
else
    echo -e "${RED}Error: Installation failed. claudes not found in PATH.${NC}"
    exit 1
fi
