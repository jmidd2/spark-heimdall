package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"spark-heimdall/internal/device"
	"strconv"
)

// Manager defines the interface for configuration operations
type Manager interface {
	load() error
	save() error
	Validate() error
	AddDevice(d device.Device) error
	UpdateDevice(d device.Device) error
	DeleteDevice(id string) error
	GetDevice(id string) (device.Device, bool)
	Update(config UpdateConfig) error
}

type UpdateConfig struct {
	ListenPort      int    `json:"listen_port"`
	AutoStart       bool   `json:"auto_start"`
	AutoStartID     string `json:"auto_start_id"`
	VncViewer       string `json:"vnc_viewer"`
	VncPasswordFile string `json:"vnc_password_file"`
	RdpViewer       string `json:"rdp_viewer"`
}

// Ensure Config implements Manager
var _ Manager = (*Config)(nil)

func (c *Config) Update(config UpdateConfig) error {
	c.ListenPort = config.ListenPort
	c.AutoStart = config.AutoStart
	c.AutoStartID = config.AutoStartID

	return c.save()
}

// Config holds the application configuration
type Config struct {
	// FilePath is the path to the current configuration file
	FilePath string `json:"-"`
	// ListenPort determines the port of the HTTP server
	ListenPort int `json:"listen_port"`

	AutoStart   bool   `json:"auto_start"`
	AutoStartID string `json:"auto_start_id"`

	VncViewer       string `json:"vnc_viewer"`
	VncPasswordFile string `json:"vnc_password_file"`
	RdpViewer       string `json:"rdp_viewer"`

	HighestDeviceId string `json:"-"`
	device.Store
}

func NewConfig(path string, vncPasswdFile string) *Config {
	return &Config{
		FilePath:        path,
		ListenPort:      8080,
		Store:           device.Store{Devices: []device.Device{}},
		VncPasswordFile: vncPasswdFile,
	}
}

// LoadConfigFromFlags loads configuration from flags or environment variables
func LoadConfigFromFlags() (*Config, error) {
	configFilePtr := flag.String("config", getEnvString("HEIMDALL_CONFIG", "config.json"), "Path to configuration file")
	portPtr := flag.Int("port", getEnvInt("HEIMDALL_PORT", 8080), "Port to listen on")
	vncViewerPtr := flag.String("vnc", getEnvString("HEIMDALL_VNC_VIEWER", "vncviewer"), "VNC viewer executable")
	vncPasswordFilePtr := flag.String("vnc-password-file", getEnvString("HEIMDALL_VNC_PASSWORD_FILE", fmt.Sprintf("%s/.vnc/passwd", getUserHomeDir())), "VNC password file")
	rdpViewerPtr := flag.String("rdp", getEnvString("HEIMDALL_RDP_VIEWER", ""), "RDP viewer executable")
	flag.Parse()

	config := NewConfig(*configFilePtr, *vncPasswordFilePtr)

	// load from file
	err := config.load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Override with command line flags if provided
	if *portPtr < 1 {
		config.ListenPort = *portPtr
	}

	if *vncViewerPtr != "" {
		config.VncViewer = *vncViewerPtr
	}

	if *rdpViewerPtr != "" {
		config.RdpViewer = *rdpViewerPtr
	}

	if *vncPasswordFilePtr != "" {
		config.VncPasswordFile = *vncPasswordFilePtr
	}

	return config, nil
}

func getUserHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// Helper functions for environment variables
func getEnvString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func (c *Config) load() error {
	// Check if file exists
	_, err := os.Stat(c.FilePath)
	if os.IsNotExist(err) {
		// Create directory if it doesn't exist
		dir := filepath.Dir(c.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		// save default config
		return c.save()
	} else if err != nil {
		return fmt.Errorf("failed to check config file: %w", err)
	}

	// Read file
	data, err := os.ReadFile(c.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := c.Validate(); err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	return nil
}

func (c *Config) save() error {
	if err := c.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	dir := filepath.Dir(c.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	if err := os.WriteFile(c.FilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (c *Config) Validate() error {
	if c.ListenPort <= 0 || c.ListenPort > 65535 {
		return errors.New("listen address must be a valid port number")
	}

	deviceIdMap := make(map[string]bool)
	for _, pc := range c.Store.Devices {
		if deviceIdMap[pc.ID] {
			return fmt.Errorf("duplicate PC ID: %s", pc.ID)
		}
		deviceIdMap[pc.ID] = true
	}

	// Verify AutoStartID references a valid PC
	if c.AutoStart && c.AutoStartID != "" {
		valid := false
		for _, pc := range c.Store.Devices {
			if pc.ID == c.AutoStartID {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("auto start ID %s does not reference a valid PC", c.AutoStartID)
		}
	}

	if c.VncViewer == "" {
		c.VncViewer = "vncviewer"
	}

	if c.VncPasswordFile == "" {
		c.VncPasswordFile = fmt.Sprintf("%s/.vnc/passwd", getUserHomeDir())
	}

	if c.RdpViewer == "" {
		c.RdpViewer = "" // TODO: Figure out a good default rdp viewer
	}

	return nil
}

func (c *Config) AddDevice(device device.Device) error {
	err := c.Store.Add(device)
	if err != nil {
		return err
	}
	log.Printf("Added new device: (%s) %s", device.ID, device.Name)
	return c.save()
}

func (c *Config) UpdateDevice(d device.Device) error {
	err := c.Store.Update(d)
	if err != nil {
		return err
	}

	return c.save()
}

func (c *Config) DeleteDevice(id string) error {
	err := c.Store.Delete(id)
	if err != nil {
		return err
	}

	// Update references
	if c.AutoStartID == id {
		c.AutoStartID = ""
	}

	return c.save()
}

func (c *Config) GetDevice(id string) (d device.Device, found bool) {
	if d, found = c.Store.Get(id); found {
		return d, true
	}

	return d, false
}
