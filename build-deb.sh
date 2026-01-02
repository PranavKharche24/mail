#!/bin/bash

# Gomail .deb Package Builder
# Creates a portable Debian package with no external dependencies

set -e

VERSION="1.0.0"
PACKAGE_NAME="gomail"
MAINTAINER="Pranav Kharche <pranavkharche24@gmail.com>"
DESCRIPTION="Professional Email Utility - CLI and Web interface for sending emails via Gmail"
ARCH="amd64"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}Building Gomail .deb package v${VERSION}${NC}"
echo "============================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

# Create build directory
BUILD_DIR="build"
DEB_DIR="${BUILD_DIR}/${PACKAGE_NAME}_${VERSION}_${ARCH}"

echo -e "${YELLOW}[1/6] Cleaning previous builds...${NC}"
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

echo -e "${YELLOW}[2/6] Compiling binary (static, no CGO)...${NC}"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=${VERSION}" -o ${BUILD_DIR}/gomail .

echo -e "${YELLOW}[3/6] Creating package structure...${NC}"
mkdir -p ${DEB_DIR}/DEBIAN
mkdir -p ${DEB_DIR}/usr/bin
mkdir -p ${DEB_DIR}/usr/share/gomail/templates
mkdir -p ${DEB_DIR}/usr/share/doc/gomail
mkdir -p ${DEB_DIR}/etc/gomail

# Copy binary
cp ${BUILD_DIR}/gomail ${DEB_DIR}/usr/bin/
chmod 755 ${DEB_DIR}/usr/bin/gomail

# Copy templates
cp templates/*.html ${DEB_DIR}/usr/share/gomail/templates/

# Copy documentation
cp README.md ${DEB_DIR}/usr/share/doc/gomail/
cp .env.example ${DEB_DIR}/etc/gomail/gomail.env.example

echo -e "${YELLOW}[4/6] Creating control file...${NC}"
cat > ${DEB_DIR}/DEBIAN/control << EOF
Package: ${PACKAGE_NAME}
Version: ${VERSION}
Section: mail
Priority: optional
Architecture: ${ARCH}
Maintainer: ${MAINTAINER}
Description: ${DESCRIPTION}
 Gomail is a professional email utility for sending emails via Gmail SMTP.
 Features include CLI and Web interfaces, HTML email support, and attachments.
 .
 No external dependencies required.
Homepage: https://github.com/pranavKharche24/mail
EOF

echo -e "${YELLOW}[5/6] Creating post-install script...${NC}"
cat > ${DEB_DIR}/DEBIAN/postinst << 'EOF'
#!/bin/bash
set -e

# Create gomail config directory in user home if running as regular user
if [ -n "$SUDO_USER" ]; then
    USER_HOME=$(getent passwd "$SUDO_USER" | cut -d: -f6)
    mkdir -p "$USER_HOME/.config/gomail"
    mkdir -p "$USER_HOME/.config/gomail/templates"
    mkdir -p "$USER_HOME/.config/gomail/uploads"
    
    # Copy templates
    cp -n /usr/share/gomail/templates/*.html "$USER_HOME/.config/gomail/templates/" 2>/dev/null || true
    
    # Copy example env if not exists
    if [ ! -f "$USER_HOME/.config/gomail/.env" ]; then
        cp /etc/gomail/gomail.env.example "$USER_HOME/.config/gomail/.env"
    fi
    
    # Set ownership
    chown -R "$SUDO_USER:$SUDO_USER" "$USER_HOME/.config/gomail"
fi

echo ""
echo "Gomail installed successfully!"
echo ""
echo "Quick start:"
echo "  1. Configure credentials: nano ~/.config/gomail/.env"
echo "  2. Run: gomail"
echo ""
echo "Usage:"
echo "  gomail          - Start both CLI and Web interface"
echo "  gomail cli      - CLI mode only"
echo "  gomail web      - Web mode only"
echo "  gomail help     - Show help"
echo ""

exit 0
EOF
chmod 755 ${DEB_DIR}/DEBIAN/postinst

# Create prerm script
cat > ${DEB_DIR}/DEBIAN/prerm << 'EOF'
#!/bin/bash
set -e
echo "Removing Gomail..."
exit 0
EOF
chmod 755 ${DEB_DIR}/DEBIAN/prerm

echo -e "${YELLOW}[6/6] Building .deb package...${NC}"
dpkg-deb --build ${DEB_DIR}

# Move to dist folder
mkdir -p dist
mv ${DEB_DIR}.deb dist/

# Calculate checksums
cd dist
sha256sum ${PACKAGE_NAME}_${VERSION}_${ARCH}.deb > ${PACKAGE_NAME}_${VERSION}_${ARCH}.deb.sha256
cd ..

# Get file size
SIZE=$(du -h dist/${PACKAGE_NAME}_${VERSION}_${ARCH}.deb | cut -f1)

echo ""
echo -e "${GREEN}Build complete!${NC}"
echo "============================================="
echo "Package: dist/${PACKAGE_NAME}_${VERSION}_${ARCH}.deb"
echo "Size: ${SIZE}"
echo ""
echo "To install:"
echo "  sudo dpkg -i dist/${PACKAGE_NAME}_${VERSION}_${ARCH}.deb"
echo ""
echo "To uninstall:"
echo "  sudo dpkg -r gomail"
echo ""
