#!/bin/bash

REPO_OWNER="ketches"
REPO_NAME="kube-recycle-bin"
BINARY_NAME="krb-cli"
INSTALL_DIR="/usr/local/bin"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64) ARCH="arm64" ;;
    i386) ARCH="386" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

LATEST_VERSION=$(basename $(curl -s -w %{redirect_url} https://github.com/$REPO_OWNER/$REPO_NAME/releases/latest) | sed 's/^v//')

if [ -z "$LATEST_VERSION" ]; then
    echo "Failed to get latest version"
    exit 1
fi

DOWNLOAD_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/v$LATEST_VERSION/${BINARY_NAME}_${LATEST_VERSION}_${OS}_${ARCH}"

TMP_DIR=$(mktemp -d)

echo "Installing $BINARY_NAME v$LATEST_VERSION ($OS/$ARCH)"
echo "Downloading from: $DOWNLOAD_URL"

curl -L $DOWNLOAD_URL -o "$TMP_DIR/$BINARY_NAME" || {
    echo "Download failed"
    rm -rf "$TMP_DIR"
    exit 1
}

sudo install -d "$INSTALL_DIR"
sudo install -m 755 "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME" || {
    echo "Installation failed"
    rm -rf "$TMP_DIR"
    exit 1
}

rm -rf "$TMP_DIR"

if command -v $BINARY_NAME >/dev/null; then
    echo "Successfully installed $BINARY_NAME to $INSTALL_DIR"
    echo "$($BINARY_NAME version 2>/dev/null || echo 'version check not supported')"
else
    echo "Installation completed but binary not found in PATH"
    exit 1
fi