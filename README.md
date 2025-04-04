# Heimdall - Remote Desktop Connection Manager

Heimdall is a lightweight web-based application that helps you manage and connect to remote computers using VNC and RDP protocols. It provides an easy-to-use interface for organizing your remote connections and quickly accessing them from a central dashboard.

## Features

- Web-based management interface accessible from any browser
- Support for both VNC and RDP protocols
- Save connection details for quick access
- Auto-start option for frequently used connections
- Configurable through CLI flags, environment variables, or configuration file

## Requirements

- Go 1.18 or later
- A VNC viewer (default: `vncviewer`)
- An RDP client (configured through settings)
- Just command runner (for build automation)

## Installation

### From Source

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/spark-heimdall.git
   cd spark-heimdall
   ```

2. Build the application using Just:
   ```
   just build
   ```

3. Run the application:
   ```
   just run
   ```

### Using Pre-built Binaries

1. Download the appropriate binary for your platform from the [Releases](https://github.com/yourusername/spark-heimdall/releases) page.

2. Make the file executable (Linux/macOS):
   ```
   chmod +x heimdall-*
   ```

3. Run the application:
   ```
   ./heimdall-*
   ```

## Configuration

Heimdall can be configured through command-line flags, environment variables, or a configuration file. The priority order is:

1. Command-line flags
2. Environment variables
3. Configuration file
4. Default values

### Configuration File

By default, Heimdall looks for a configuration file at `./config.json`. The configuration is stored in JSON format and includes the following options:

```json
{
  "listen_port": 8080,
  "auto_start": false,
  "auto_start_id": "",
  "vnc_viewer": "vncviewer",
  "vnc_password_file": "/home/user/.vnc/passwd",
  "rdp_viewer": "xfreerdp",
  "devices": [
    {
      "id": "unique-id",
      "name": "My PC",
      "ip_address": "192.168.1.100",
      "port": 5900,
      "protocol": "vnc",
      "username": "",
      "password": "",
      "full_screen": false
    }
  ]
}
```

### Command-Line Options

| Flag | Description | Default | Environment Variable |
|------|-------------|---------|---------------------|
| `-config` | Path to configuration file | `config.json` | `HEIMDALL_CONFIG` |
| `-port` | HTTP server port | `8080` | `HEIMDALL_PORT` |
| `-vnc` | Path to VNC viewer executable | `vncviewer` | `HEIMDALL_VNC_VIEWER` |
| `-vnc-password-file` | Path to VNC password file | `$HOME/.vnc/passwd` | `HEIMDALL_VNC_PASSWORD_FILE` |
| `-rdp` | Path to RDP client executable | `""` | `HEIMDALL_RDP_VIEWER` |

### Environment Variables

You can use environment variables instead of command-line flags:

```
HEIMDALL_CONFIG=/path/to/config.json
HEIMDALL_PORT=8080
HEIMDALL_VNC_VIEWER=/usr/bin/vncviewer
HEIMDALL_VNC_PASSWORD_FILE=/home/user/.vnc/passwd
HEIMDALL_RDP_VIEWER=/usr/bin/xfreerdp
```

## Usage

1. Start Heimdall:
   ```
   just run
   ```

2. Open your web browser and navigate to `http://localhost:8080` (or whatever port you configured)

3. Add your remote computers through the web interface

4. Click on a computer to connect to it

## Development

### Build Tools

Heimdall uses [Just](https://github.com/casey/just) as a command runner for build automation. Install Just on your system before developing.

Common Just commands:

```
just                    # List all available commands
just build              # Build for current platform
just build-all          # Build for all platforms
just run                # Build and run
just release-notes 1.0.0 # Generate release notes for version 1.0.0
just full-release 1.0.0  # Tag a new release and generate notes
```

### Commit Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) standard for commit messages. This enables automatic generation of release notes and changelogs.

Examples:
- `feat: add support for RDP connections`
- `fix: resolve connection timeout issue`
- `docs: update installation instructions`
- `chore: update dependencies`

### Pull Request Process

1. Ensure your code follows the project's coding standards
2. Update the README.md with details of changes where appropriate
3. Follow the Conventional Commits standard in your commit messages
4. The PR title should summarize the changes and follow the Conventional Commits format
5. Link any related issues in the PR description

### Project Structure

- `main.go` - Entry point of the application
- `internal/heimdall/server.go` - HTTP server implementation
- `internal/config/config.go` - Configuration management
- `internal/device/device.go` - Device management

## Security Considerations

- Heimdall stores connection details including passwords in the configuration file. Ensure this file has appropriate permissions.
- For VNC connections, Heimdall uses a VNC password file.
- The web interface does not include authentication, so it should only be run on trusted networks.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request following our contribution guidelines above.