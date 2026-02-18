#!/bin/bash

set -e

echo "ServerPanel Installer"
echo "===================="

# Check if running as root
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit 1
fi

# Installation directory
INSTALL_DIR="/opt/serverpanel"
BIN_DIR="/usr/local/bin"
SERVICE_FILE="/etc/systemd/system/serverpanel.service"

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
  x86_64)
    BINARY="serverpanel-linux-amd64"
    ;;
  aarch64|arm64)
    BINARY="serverpanel-linux-arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Create installation directory
mkdir -p $INSTALL_DIR

# Stop service if it's running (to avoid "Text file busy" error during re-installation)
if systemctl is-active --quiet serverpanel 2>/dev/null; then
  echo "Stopping existing ServerPanel service..."
  systemctl stop serverpanel
fi

# Copy binary
if [ -f "$BINARY" ]; then
  cp $BINARY $INSTALL_DIR/serverpanel
  chmod +x $INSTALL_DIR/serverpanel
else
  echo "Binary $BINARY not found"
  exit 1
fi

# Create symlink
ln -sf $INSTALL_DIR/serverpanel $BIN_DIR/serverpanel

# Create systemd service
cat > $SERVICE_FILE << 'EOF'
[Unit]
Description=ServerPanel - Linux Server Management
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/serverpanel
ExecStart=/opt/serverpanel/serverpanel
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd
systemctl daemon-reload

# Enable and start service
systemctl enable serverpanel
systemctl start serverpanel

echo ""
echo "Installation complete!"
echo "ServerPanel is running on port 8888"
echo ""
echo "Access it at: http://your-server-ip:8888"
echo "Default credentials: admin / admin123"
echo ""
echo "Commands:"
echo "  Start:   systemctl start serverpanel"
echo "  Stop:    systemctl stop serverpanel"
echo "  Status:  systemctl status serverpanel"
echo "  Logs:    journalctl -u serverpanel -f"
