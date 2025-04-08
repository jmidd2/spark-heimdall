package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"spark-heimdall/internal/config"
	"spark-heimdall/internal/heimdall"
	"spark-heimdall/internal/logger"
	"syscall"
)

// Version information - populated at build time
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

//go:embed static
var staticFS embed.FS

//go:embed templates
var templateFS embed.FS

func main() {
	// Parse command line flags
	versionFlag := flag.Bool("version", false, "Print version information")
	configFlag := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	// Show version information and exit if requested
	if *versionFlag {
		fmt.Printf("Heimdall %s (commit: %s, built: %s)\n", version, commit, date)
		return
	}

	// Set configuration file from flag if provided
	if *configFlag != "" {
		os.Setenv("HEIMDALL_CONFIG", *configFlag)
	}

	// Initialize logger
	logger.Info("Starting Heimdall...").WithField("version", version).Send()

	// Load configuration
	cfg, err := config.LoadConfig(*configFlag)
	if err != nil {
		logger.Fatal("Failed to load configuration").WithError(err).Send()
	}

	// Set up static file system
	staticSubFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		logger.Fatal("Failed to create sub file system for static files").WithError(err).Send()
	}

	// Create server
	server, err := heimdall.NewServer(cfg, staticSubFS, templateFS)
	if err != nil {
		logger.Fatal("Failed to create server").WithError(err).Send()
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	errChan := make(chan error)
	go func() {
		logger.Info("Starting server").
			WithField("port", cfg.Server.Port).
			Send()

		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Auto-start connection if configured
	if cfg.Connection.AutoStart && cfg.Connection.AutoStartID != "" {
		device, found := cfg.GetDevice(cfg.Connection.AutoStartID)
		if found {
			logger.Info("Auto-starting connection").
				WithField("device", device.Name).
				WithField("id", device.ID).
				Send()

			go server.ConnectToDevice(device)
		} else {
			logger.Warn("Auto-start device not found").
				WithField("id", cfg.Connection.AutoStartID).
				Send()
		}
	}

	// Wait for server error or shutdown signal
	select {
	case err := <-errChan:
		logger.Fatal("Server error").WithError(err).Send()
	case signal := <-sigChan:
		logger.Info("Received shutdown signal").
			WithField("signal", signal.String()).
			Send()
	}

	// Perform cleanup
	logger.Info("Shutting down server...").Send()
	server.Shutdown()
	logger.Info("Server stopped").Send()
}
