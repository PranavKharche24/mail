.PHONY: build clean deb release install uninstall

VERSION := 1.0.0
BINARY := gomail
PACKAGE := gomail

# Build flags
LDFLAGS := -ldflags="-s -w -X main.version=$(VERSION)"

# Default target
all: build

# Build binary
build:
	@echo "Building $(BINARY)..."
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BINARY) .

# Build for all platforms
build-all: build-linux build-darwin build-windows

build-linux:
	@echo "Building for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-amd64 .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-arm64 .

build-darwin:
	@echo "Building for macOS..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-arm64 .

build-windows:
	@echo "Building for Windows..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-windows-amd64.exe .

# Build .deb package
deb: build-linux
	@echo "Building .deb package..."
	@chmod +x build-deb.sh
	@./build-deb.sh

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BINARY) build dist

# Install locally
install: build
	@echo "Installing to /usr/local/bin..."
	sudo cp $(BINARY) /usr/local/bin/
	sudo chmod 755 /usr/local/bin/$(BINARY)
	@echo "Installed successfully"

# Uninstall
uninstall:
	@echo "Uninstalling..."
	sudo rm -f /usr/local/bin/$(BINARY)
	@echo "Uninstalled successfully"

# Run
run: build
	./$(BINARY)

# Test
test:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Create release tag
release:
	@if [ -z "$(v)" ]; then echo "Usage: make release v=1.0.0"; exit 1; fi
	@echo "Creating release v$(v)..."
	git tag -a v$(v) -m "Release v$(v)"
	git push origin v$(v)
	@echo "Release v$(v) created. GitHub Actions will build and publish."

# Help
help:
	@echo "Gomail Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build      - Build binary"
	@echo "  make build-all  - Build for all platforms"
	@echo "  make deb        - Build .deb package"
	@echo "  make install    - Install to /usr/local/bin"
	@echo "  make uninstall  - Remove from /usr/local/bin"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make run        - Build and run"
	@echo "  make release v=X.X.X - Create release tag"
	@echo ""
