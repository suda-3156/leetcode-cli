#!/bin/bash

set -euo pipefail

DEFAULT_BINARY_NAME="lcli"
BINARY_NAME="${1:-$DEFAULT_BINARY_NAME}"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Building ${BINARY_NAME}...${NC}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_ROOT"

if go build -o "$BINARY_NAME" .; then
    echo -e "${GREEN}✓ Build successful${NC}"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi

GOPATH="${GOPATH:-$HOME/go}"
INSTALL_DIR="$GOPATH/bin"

if [ ! -d "$INSTALL_DIR" ]; then
    echo -e "${YELLOW}Creating directory: $INSTALL_DIR${NC}"
    mkdir -p "$INSTALL_DIR"
fi

echo -e "${YELLOW}Installing to $INSTALL_DIR/$BINARY_NAME...${NC}"
if cp "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"; then
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    echo -e "${GREEN}✓ Successfully installed to $INSTALL_DIR/$BINARY_NAME${NC}"
else
    echo -e "${RED}✗ Installation failed${NC}"
    exit 1
fi

rm -f "$BINARY_NAME"

echo -e "${GREEN}Done! You can now run: ${BINARY_NAME}${NC}"
echo -e "${YELLOW}Make sure $INSTALL_DIR is in your PATH${NC}"

