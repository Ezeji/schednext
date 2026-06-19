#!/usr/bin/bash

set -e

###############################################################################
# CONFIGURATION
###############################################################################

REPO_OWNER="Ezeji"
REPO_NAME="schednext"

INSTALL_DIR="/opt/schednext-runtime"
BINARY_DIR="/opt/schednext-runtime/bin"
STATELENS_DIR="/statelens"

###############################################################################
# ROOT CHECK
###############################################################################

if [ "$(id -u)" -ne 0 ]; then
    echo ""
    echo "ERROR: SchedNext installation must be run as root."
    echo ""
    echo "Example:"
    echo "  sudo bash install.sh"
    echo ""
    exit 1
fi

###############################################################################
# VERSION RESOLUTION
###############################################################################

REQUESTED_VERSION="${1:-latest}"

if [ "$REQUESTED_VERSION" = "latest" ]; then

    VERSION=$(curl -fsSL \
        "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" \
        | grep '"tag_name"' \
        | sed -E 's/.*"([^"]+)".*/\1/')

else

    VERSION="$REQUESTED_VERSION"

fi

echo ""
echo "Installing SchedNext Runtime"
echo "Version: ${VERSION}"
echo ""

###############################################################################
# ARCH DETECTION
###############################################################################

ARCH=$(uname -m)

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "Architecture: ${ARCH}"

###############################################################################
# CREATE DIRECTORIES
###############################################################################

mkdir -p "$INSTALL_DIR"
mkdir -p "$BINARY_DIR"
mkdir -p "$STATELENS_DIR"

###############################################################################
# DOWNLOAD BINARIES
###############################################################################

AGENT_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${VERSION}/schednext-agent-linux-${ARCH}"

CLI_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${VERSION}/schednext-linux-${ARCH}"

echo ""
echo "Downloading agent..."
curl -fL "$AGENT_URL" -o "${BINARY_DIR}/schednext-agent"

echo "Downloading cli..."
curl -fL "$CLI_URL" -o "${BINARY_DIR}/schednext"

chmod +x "${BINARY_DIR}/schednext-agent"
chmod +x "${BINARY_DIR}/schednext"

###############################################################################
# SAMPLE CONFIG
###############################################################################

CONFIG_FILE="${INSTALL_DIR}/schednext.config"

if [ ! -f "$CONFIG_FILE" ]; then

cat > "$CONFIG_FILE" <<EOF
{
  "version": 1,
  "jobs": [
    {
      "id": "heartbeat",
      "binary": "heartbeat.sh",
      "cron": "* * * * *",
      "enabled": true,
      "maxRuntimeSeconds": 5
    }
  ]
}
EOF

fi

###############################################################################
# SAMPLE JOB
###############################################################################

HEARTBEAT_JOB="${INSTALL_DIR}/heartbeat.sh"

if [ ! -f "$HEARTBEAT_JOB" ]; then

cat > "$HEARTBEAT_JOB" <<EOF
#!/usr/bin/bash

echo "[heartbeat] alive at \$(date)"
sleep 10
exit 0
EOF

chmod +x "$HEARTBEAT_JOB"

fi

###############################################################################
# SUCCESS MESSAGE
###############################################################################

echo ""
echo "=============================================="
echo "SchedNext Runtime installed successfully"
echo "=============================================="
echo ""
echo "Install directory:"
echo "  ${INSTALL_DIR}"
echo ""
echo "Binary directory:"
echo "  ${BINARY_DIR}"
echo ""
echo "StateLens mount directory:"
echo "  ${STATELENS_DIR}"
echo ""
echo "Agent:"
echo "  ${BINARY_DIR}/schednext-agent"
echo ""
echo "CLI:"
echo "  ${BINARY_DIR}/schednext-cli"
echo ""
echo "Config:"
echo "  ${INSTALL_DIR}/schednext.config"
echo ""
echo "Example:"
echo ""
echo "  cd ${BINARY_DIR}"
echo "  ./schednext-agent"
echo ""