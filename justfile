# Justfile for Heimdall

# Default recipe to run when just is called without arguments
default:
    @just --list

# Build the application for the current platform
build:
    @echo "Building Heimdall..."
    go build -o bin/heimdall main.go

# Run the application
run: build
    @echo "Running Heimdall..."
    ./bin/heimdall

# Clean build artifacts
clean:
    @echo "Cleaning build artifacts..."
    rm -rf bin dist

# Install dependencies
deps:
    @echo "Installing dependencies..."
    go mod download
    go mod tidy

# Run tests
test:
    @echo "Running tests..."
    go test -v ./...

# Build binary for Linux or macOS
build-unix os arch:
    @echo "Building for {{os}}/{{arch}}..."
    mkdir -p dist
    GOOS={{os}} GOARCH={{arch}} go build -o "dist/spark-heimdall-{{os}}-{{arch}}" -ldflags "-s -w" ./main.go

# Build binary for Windows
build-windows arch:
    @echo "Building for windows/{{arch}}..."
    mkdir -p dist
    GOOS=windows GOARCH={{arch}} go build -o "dist/spark-heimdall-windows-{{arch}}.exe" -ldflags "-s -w" ./main.go

# Build binaries for all platforms (called by build.sh)
build-all: clean
    @echo "Building for all platforms..."
    mkdir -p dist
    just build-unix linux amd64
    just build-unix linux arm64
    just build-unix darwin amd64
    just build-unix darwin arm64
    just build-windows amd64

# Create release archives
package: build-all
    @echo "Creating release packages..."
    cd dist && tar -czf spark-heimdall-linux-amd64.tar.gz spark-heimdall-linux-amd64
    cd dist && tar -czf spark-heimdall-linux-arm64.tar.gz spark-heimdall-linux-arm64
    cd dist && tar -czf spark-heimdall-darwin-amd64.tar.gz spark-heimdall-darwin-amd64
    cd dist && tar -czf spark-heimdall-darwin-arm64.tar.gz spark-heimdall-darwin-arm64
    cd dist && zip spark-heimdall-windows-amd64.zip spark-heimdall-windows-amd64.exe

# Tag and release a new version
release VERSION:
    @echo "Tagging release v{{VERSION}}..."
    git tag v{{VERSION}}
    git push origin v{{VERSION}}
    @echo "Release v{{VERSION}} tagged. GitHub Actions will build and publish the release after a PR."

# Run with specific config file
run-with-config CONFIG:
    @echo "Running with config {{CONFIG}}..."
    ./bin/heimdall -config {{CONFIG}}

# Dev mode - build and run with hot reloading (requires air: https://github.com/cosmtrek/air)
#dev:
#    @echo "Running in development mode with hot reload..."
#    air -c .air.toml

# Check for lint issues (requires golangci-lint)
#lint:
#    @echo "Running linter..."
#    golangci-lint run

# Format code
fmt:
    @echo "Formatting code..."
    go fmt ./...

# Build for current platform with version info
build-version VERSION:
    @echo "Building version v{{VERSION}}..."
    go build -o bin/heimdall -ldflags "-s -w -X main.version=v{{VERSION}}" main.go

# Generate release notes from git history
release-notes VERSION:
    @echo "Generating release notes for v{{VERSION}}..."
    @echo "# Heimdall v{{VERSION}} Release" > release-notes.md
    @echo "" >> release-notes.md
    @echo "## Overview" >> release-notes.md
    @echo "Heimdall is a web-based remote desktop connection manager that supports VNC and RDP protocols." >> release-notes.md
    @echo "" >> release-notes.md
    @echo "## What's New" >> release-notes.md

    @# Get features from commits
    @echo "" >> release-notes.md
    @git log --pretty=format:"- %s" $(git describe --tags --abbrev=0 2>/dev/null || echo HEAD^)..HEAD | grep -E "^- feat(\([^)]+\))?:" >> release-notes.md || true

    @echo "" >> release-notes.md
    @echo "## Bug Fixes" >> release-notes.md
    @git log --pretty=format:"- %s" $(git describe --tags --abbrev=0 2>/dev/null || echo HEAD^)..HEAD | grep -E "^- fix(\([^)]+\))?:" >> release-notes.md || true

    @echo "" >> release-notes.md
    @echo "## Installation" >> release-notes.md
    @echo "" >> release-notes.md
    @echo "### Binary Downloads" >> release-notes.md
    @echo "Pre-built binaries are available for:" >> release-notes.md
    @echo "- Windows (64-bit)" >> release-notes.md
    @echo "- Linux (64-bit, ARM64)" >> release-notes.md
    @echo "- macOS (Intel, Apple Silicon)" >> release-notes.md
    @echo "" >> release-notes.md
    @echo "### Quick Start" >> release-notes.md
    @echo "1. Download the appropriate binary for your platform" >> release-notes.md
    @echo "2. Make the file executable (Linux/macOS): \`chmod +x heimdall-*\`" >> release-notes.md
    @echo "3. Run the application: \`./heimdall-*\`" >> release-notes.md
    @echo "4. Access the web interface at \`http://localhost:8080\`" >> release-notes.md
    @echo "" >> release-notes.md
    @echo "## Configuration" >> release-notes.md
    @echo "See our [README](https://github.com/yourusername/spark-heimdall#configuration) for configuration options and details." >> release-notes.md
    @echo "" >> release-notes.md
    @cat release-notes.md

# Create a release with generated notes
full-release VERSION: (release-notes VERSION)
    @just release {{VERSION}}
    @echo "Release v{{VERSION}} created with generated notes."