package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"spark-heimdall/internal/logger"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Device represents any connectable device (PC, server, etc.)
type Device struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IPAddress   string `json:"ip_address"`
	Protocol    string `json:"protocol"` // "vnc", "rdp", etc.
	Port        int    `json:"port"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	FullScreen  bool   `json:"full_screen"`
	Description string `json:"description,omitempty"`
	Screen      string `json:"screen,omitempty"`
}

type Devices []Device

// Config holds the application configuration
type Config struct {
	// Server configuration
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`

	// Connection configuration
	Connection struct {
		AutoStart   bool   `mapstructure:"auto_start"`
		AutoStartID string `mapstructure:"auto_start_id"`
	} `mapstructure:"connection"`

	// Client applications
	Clients struct {
		VncViewer       string `mapstructure:"vnc_viewer"`
		VncPasswordFile string `mapstructure:"vnc_password_file"`
		RdpViewer       string `mapstructure:"rdp_viewer"`
	} `mapstructure:"clients"`

	// Devices list
	Devices []Device `mapstructure:"devices"`

	// Security configuration
	Security struct {
		KeyFile string `mapstructure:"key_file"`
	} `mapstructure:"security"`

	// Logger configuration
	Logging struct {
		Level  string `mapstructure:"level"`
		Format string `mapstructure:"format"`
	} `mapstructure:"logging"`

	// Internal fields
	configFile string
	v          *viper.Viper
}

// LoadConfig loads the configuration from file, environment variables, and flags
func LoadConfig(configFlag string) (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Load configuration from file
	configFile := configFlag
	if configFile == "" {
		configFile = os.Getenv("HEIMDALL_CONFIG")
	}
	if configFile == "" {
		configFile = "config.json"
	}

	// Check if file exists
	if _, err := os.Stat(configFile); err == nil {
		v.SetConfigFile(configFile)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Load configuration from environment variables
	v.SetEnvPrefix("HEIMDALL")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Parse the configuration
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Store viper and config file for later use
	config.v = v
	config.configFile = configFile

	// Validate and normalize configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Configure logging based on the config
	setupLogging(config.Logging.Level, config.Logging.Format)

	return &config, nil
}

// setDefaults sets the default values for configuration
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", 8080)

	// Connection defaults
	v.SetDefault("connection.auto_start", false)
	v.SetDefault("connection.auto_start_id", "")

	// Client defaults
	v.SetDefault("clients.vnc_viewer", "vncviewer")

	homeDir, err := os.UserHomeDir()
	if err == nil {
		v.SetDefault("clients.vnc_password_file", filepath.Join(homeDir, ".vnc", "passwd"))
	} else {
		v.SetDefault("clients.vnc_password_file", ".vnc/passwd")
	}

	v.SetDefault("clients.rdp_viewer", "")

	// Security defaults
	v.SetDefault("security.key_file", ".heimdall/encryption.key")

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "text")
}

// setupLogging configures the logging system
func setupLogging(level, format string) {
	// Set log level
	var logLevel logger.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = logger.DebugLevel
	case "info":
		logLevel = logger.InfoLevel
	case "warn", "warning":
		logLevel = logger.WarnLevel
	case "error":
		logLevel = logger.ErrorLevel
	default:
		logLevel = logger.InfoLevel
	}
	logger.SetLevel(logLevel)

	// Set log format
	if strings.ToLower(format) == "json" {
		logger.EnableJSON()
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate server port
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return errors.New("server port must be a valid port number")
	}

	// Validate auto-start ID if auto-start is enabled
	if c.Connection.AutoStart && c.Connection.AutoStartID != "" {
		found := false
		for _, dev := range c.Devices {
			if dev.ID == c.Connection.AutoStartID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("auto-start ID %s does not reference a valid device", c.Connection.AutoStartID)
		}
	}

	// Validate VNC viewer
	if c.Clients.VncViewer == "" {
		return errors.New("VNC viewer cannot be empty")
	}

	// Check for duplicate device IDs
	deviceIDMap := make(map[string]bool)
	for _, dev := range c.Devices {
		if deviceIDMap[dev.ID] {
			return fmt.Errorf("duplicate device ID: %s", dev.ID)
		}
		deviceIDMap[dev.ID] = true
	}

	return nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	// Update viper with current config values
	c.v.Set("server.port", c.Server.Port)
	c.v.Set("connection.auto_start", c.Connection.AutoStart)
	c.v.Set("connection.auto_start_id", c.Connection.AutoStartID)
	c.v.Set("clients.vnc_viewer", c.Clients.VncViewer)
	c.v.Set("clients.vnc_password_file", c.Clients.VncPasswordFile)
	c.v.Set("clients.rdp_viewer", c.Clients.RdpViewer)
	c.v.Set("devices", c.Devices)
	c.v.Set("security.key_file", c.Security.KeyFile)
	c.v.Set("logging.level", c.Logging.Level)
	c.v.Set("logging.format", c.Logging.Format)

	// Validate before saving
	if err := c.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(c.configFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Save to file
	if err := c.v.WriteConfigAs(c.configFile); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// AddDevice adds a new device to the configuration
func (c *Config) AddDevice(dev Device) error {
	// Validate device ID is not already in use
	dev.ID = strconv.Itoa(len(c.Devices))
	for _, existingDev := range c.Devices {
		if existingDev.ID == dev.ID {
			return fmt.Errorf("device with ID %s already exists", dev.ID)
		}
	}

	// Add device to list
	c.Devices = append(c.Devices, dev)

	// Save configuration
	return c.Save()
}

// UpdateDevice updates an existing device
func (c *Config) UpdateDevice(dev Device) error {
	for i, existingDev := range c.Devices {
		if existingDev.ID == dev.ID {
			c.Devices[i] = dev
			return c.Save()
		}
	}

	return fmt.Errorf("device with ID %s not found", dev.ID)
}

// DeleteDevice removes a device from the configuration
func (c *Config) DeleteDevice(id string) error {
	found := false
	newDevices := make([]Device, 0, len(c.Devices))

	for _, dev := range c.Devices {
		if dev.ID == id {
			found = true
		} else {
			newDevices = append(newDevices, dev)
		}
	}

	if !found {
		return fmt.Errorf("device with ID %s not found", id)
	}

	c.Devices = newDevices

	// If auto-start device was deleted, clear auto-start ID
	if c.Connection.AutoStartID == id {
		c.Connection.AutoStartID = ""
	}

	return c.Save()
}

// GetDevice retrieves a device by ID
func (c *Config) GetDevice(id string) (Device, bool) {
	for _, dev := range c.Devices {
		if dev.ID == id {
			return dev, true
		}
	}

	return Device{}, false
}

// GetAllDevices returns all devices
func (c *Config) GetAllDevices() []Device {
	return c.Devices
}

// UpdateConfig updates configuration settings
func (c *Config) UpdateConfig(port int, autoStart bool, autoStartID, vncViewer, vncPasswordFile, rdpViewer string) error {
	c.Server.Port = port
	c.Connection.AutoStart = autoStart
	c.Connection.AutoStartID = autoStartID
	c.Clients.VncViewer = vncViewer
	c.Clients.VncPasswordFile = vncPasswordFile
	c.Clients.RdpViewer = rdpViewer

	return c.Save()
}
