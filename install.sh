#!/bin/bash
set -e

REPO="ferivoq/animedl"
VERSION="latest"
BINARY_NAME="animedrive-dl"

if [ "$(uname)" == "Darwin" ]; then
    OS="darwin"
elif [ "$(uname)" == "Linux" ]; then
    OS="linux"
elif [[ "$(uname -s)" == *"_NT"* ]]; then
    OS="windows"
else
    echo "Unsupported OS"
    exit 1
fi

ARCH=$(uname -m)
if [ "$ARCH" == "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" == "aarch64" ]; then
    ARCH="arm64"
else
    echo "Unsupported architecture"
    exit 1
fi

URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY_NAME-$OS-$ARCH"

echo "Downloading $BINARY_NAME from $URL"
curl -L "$URL" -o /usr/local/bin/$BINARY_NAME
chmod +x /usr/local/bin/$BINARY_NAME

echo "$BINARY_NAME installed successfully"
