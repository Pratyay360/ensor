#!/bin/sh
set -eu

REPO="Pratyay360/ensor"
VERSION="${ENSOR_VERSION:-latest}"

OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
  Linux)  OS="Linux" ;;
  Darwin) OS="Darwin" ;;
  *) echo "Unsupported OS: $OS" >&2; exit 1 ;;
esac

case "$ARCH" in
  x86_64|amd64) ARCH="x86_64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  i386|i686)     ARCH="i386" ;;
  *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

EXT="tar.gz"
case "$OS" in
  Linux) ;;
  Darwin) ;;
esac
# Windows is not supported by this installer (use the zip from releases)

ASSET="ensor_${OS}_${ARCH}.${EXT}"

if [ "$VERSION" = "latest" ]; then
  URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"
else
  # ensure leading v
  case "$VERSION" in
    v*) ;;
    *) VERSION="v${VERSION}" ;;
  esac
  URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET}"
fi

INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
mkdir -p "$INSTALL_DIR"

TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT INT TERM

echo "Downloading ensor (${OS}/${ARCH}) ${VERSION}..."
echo "  -> ${URL}"

if ! curl -fL "$URL" -o "$TMPDIR/${ASSET}"; then
  echo "Failed to download ${URL}" >&2
  echo "Check https://github.com/${REPO}/releases for available assets." >&2
  exit 1
fi

tar -xzf "$TMPDIR/${ASSET}" -C "$TMPDIR"
# archive contains binary named 'ensor'
if [ ! -f "$TMPDIR/ensor" ]; then
  echo "Binary not found inside archive" >&2; exit 1
else
  cp "$TMPDIR/ensor" "$INSTALL_DIR/ensor"
fi

chmod +x "$INSTALL_DIR/ensor"

echo "Installed ensor to $INSTALL_DIR/ensor"
echo "Run: $INSTALL_DIR/ensor --help"
if ! echo "$PATH" | tr ':' '\n' | grep -qx "$INSTALL_DIR"; then
  echo "Note: $INSTALL_DIR is not in your PATH. Add it with:"
  echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
fi
