#!/bin/bash

set -e

echo "ServerPanel Remote Installer"
echo "============================="

# Check if running as root
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root or with sudo"
  exit 1
fi

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
  x86_64)
    PACKAGE="serverpanel-linux-amd64.tar.gz"
    ;;
  aarch64|arm64)
    PACKAGE="serverpanel-linux-arm64.tar.gz"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

echo "Downloading ServerPanel..."
# Download latest release
if ! wget -q https://github.com/binbankm/My/releases/latest/download/$PACKAGE; then
  echo "Failed to download ServerPanel package"
  echo "Please check your internet connection and try again"
  rm -rf "$TMP_DIR"
  exit 1
fi

echo "Extracting package..."
tar -zxf $PACKAGE

echo "Installing ServerPanel..."
./install.sh

# Cleanup
cd /
rm -rf "$TMP_DIR"

echo ""
echo "Remote installation complete!"
