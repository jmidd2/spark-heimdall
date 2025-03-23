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

## Installation

### From Source

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/spark-heimdall.git
   cd spark-heimdall
   ```

2. Build the application:
   ```
   go build -o heimdall
   ```

3. Run the application:
   ```
   ./heimdall
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
   ./heimdall
   ```

2. Open your web browser and navigate to `http://localhost:8080` (or whatever port you configured)

3. Add your remote computers through the web interface

4. Click on a computer to connect to it

## Managing Devices

### Device Properties

Each device (remote computer) has the following properties:

- **Name**: A friendly name for the device
- **IP Address**: The IP address or hostname of the remote computer
- **Port**: The port number for the remote connection (default: 5900 for VNC)
- **Protocol**: Either "vnc" or "rdp"
- **Username**: Username for RDP connections (optional)
- **Password**: Password for RDP connections (optional)
- **Full Screen**: Whether to start the connection in full-screen mode

## Security Considerations

- Heimdall stores connection details including passwords in the configuration file. Ensure this file has appropriate permissions.
- For VNC connections, Heimdall uses a VNC password file.
- The web interface does not include authentication, so it should only be run on trusted networks.

## Development

### Project Structure

- `main.go` - Entry point of the application
- `internal/heimdall/server.go` - HTTP server implementation
- `internal/config/config.go` - Configuration management
- `internal/device/device.go` - Device management (not shown in provided files)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.