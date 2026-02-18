#!/bin/bash

set -e

echo "Uninstalling ServerPanel..."

# Check if running as root
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit 1
fi

# Stop and disable service
systemctl stop serverpanel || true
systemctl disable serverpanel || true

# Remove service file
rm -f /etc/systemd/system/serverpanel.service

# Reload systemd
systemctl daemon-reload

# Remove installation directory
rm -rf /opt/serverpanel

# Remove symlink
rm -f /usr/local/bin/serverpanel

echo "ServerPanel has been uninstalled"
