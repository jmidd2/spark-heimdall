package heimdall

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"spark-heimdall/internal/config"
	"spark-heimdall/internal/logger"
	"strings"
	"sync"
	"time"
)

// Server handles HTTP requests and manages connections
type Server struct {
	configFile      *config.Config
	staticFS        fs.FS
	templatesFS     fs.FS
	httpServer      *http.Server
	cmdLock         sync.Mutex
	currentCmd      *exec.Cmd
	currentDeviceId string
}

// APIResponse provides a standardized API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, staticFS fs.FS, templatesFS fs.FS) (*Server, error) {
	return &Server{
		configFile:  cfg,
		staticFS:    staticFS,
		templatesFS: templatesFS,
	}, nil
}

// Start initializes and starts the HTTP server
func (s *Server) Start() error {
	// Set up router with all routes
	router := s.setupRouter()

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.configFile.Server.Port),
		Handler: router,
	}

	// Start the server
	logger.Info("Starting HTTP server").WithField("port", s.configFile.Server.Port).Send()
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown() error {
	// Disconnect any active connections
	s.DisconnectCurrentDevice()

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the HTTP server
	if s.httpServer != nil {
		logger.Info("Shutting down HTTP server").Send()
		return s.httpServer.Shutdown(ctx)
	}

	return nil
}

// setupRouter initializes all HTTP routes
func (s *Server) setupRouter() http.Handler {
	// Create router
	mux := http.NewServeMux()

	// API Routes
	mux.HandleFunc("/api/devices", s.handleDevices)
	mux.HandleFunc("/api/devices/", s.handleDeviceByID)
	mux.HandleFunc("/api/config", s.handleConfig)

	// Connection routes
	mux.HandleFunc("/connect/", s.handleConnect)
	mux.HandleFunc("/disconnect", s.handleDisconnect)

	// Serve static files
	mux.Handle("/", s.spaHandler())

	// Apply middleware
	return logMiddleware(corsMiddleware(mux))
}

// spaHandler handles serving the SPA and static files
func (s *Server) spaHandler() http.Handler {
	// File server for static files
	fileServer := http.FileServer(http.FS(s.staticFS))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the requested file
		path := r.URL.Path

		// Check if file exists in static FS
		_, err := fs.Stat(s.staticFS, path[1:]) // Remove leading slash
		if os.IsNotExist(err) && path != "/" && !contains(path, []string{".js", ".css", ".ico", ".png", ".svg", ".jpg", ".jpeg", ".gif"}) {
			// For SPA routing, serve index.html for routes that don't exist as files
			// but don't match specific file extensions
			r.URL.Path = "/"
		}

		// Serve the file
		fileServer.ServeHTTP(w, r)
	})
}

// Helper function to check if a string contains any of the given substrings
func contains(s string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// handleDevices handles GET and POST requests for the /api/devices endpoint
func (s *Server) handleDevices(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")

	case http.MethodGet:
		// Return all devices
		devices := s.configFile.GetAllDevices()
		s.respondWithJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Data:    devices,
		})

	case http.MethodPost:
		// Add a new device
		var dev config.Device
		if err := json.NewDecoder(r.Body).Decode(&dev); err != nil {
			logger.Error("Error decoding device").WithError(err).Send()
			s.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Handle password encryption
		//if dev.Password != "" {
		//	encryptedPwd, err := s.passwordManager.Encrypt(dev.Password)
		//	if err != nil {
		//		logger.Error("Error encrypting password").WithError(err).Send()
		//		s.respondWithError(w, http.StatusInternalServerError, "Failed to encrypt password")
		//		return
		//	}
		//	dev.Password = encryptedPwd
		//}
		newDev, err := s.configFile.AddDevice(dev)
		if err != nil {
			logger.Error("Error adding device").WithError(err).Send()
			s.respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		s.respondWithJSON(w, http.StatusCreated, APIResponse{
			Success: true,
			Data:    newDev,
		})

	default:
		s.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleDeviceByID handles requests for /api/devices/:id
func (s *Server) handleDeviceByID(w http.ResponseWriter, r *http.Request) {
	// Extract device ID from URL
	id := strings.TrimPrefix(r.URL.Path, "/api/devices/")
	if id == "" {
		s.respondWithError(w, http.StatusBadRequest, "Device ID is required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Get device by ID
		device, found := s.configFile.GetDevice(id)
		if !found {
			s.respondWithError(w, http.StatusNotFound, "Device not found")
			return
		}

		// Don't send encrypted password
		//device.Password = ""

		s.respondWithJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Data:    device,
		})

	case http.MethodPut:
		// Update existing device
		var dev config.Device
		if err := json.NewDecoder(r.Body).Decode(&dev); err != nil {
			s.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Ensure ID in URL matches ID in body
		dev.ID = id

		// Handle password encryption
		//if dev.Password != "" {
		//	encryptedPwd, err := s.passwordManager.Encrypt(dev.Password)
		//	if err != nil {
		//		logger.Error("Error encrypting password").WithError(err).Send()
		//		s.respondWithError(w, http.StatusInternalServerError, "Failed to encrypt password")
		//		return
		//	}
		//	dev.Password = encryptedPwd
		//}

		if err := s.configFile.UpdateDevice(dev); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		s.respondWithJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Data:    dev,
		})

	case http.MethodDelete:
		// Delete device
		if err := s.configFile.DeleteDevice(id); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		s.respondWithJSON(w, http.StatusOK, APIResponse{
			Success: true,
		})

	default:
		s.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleConfig handles requests for /api/config
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Return current configuration (excluding sensitive info)
		safeConfig := struct {
			Server struct {
				Port int `json:"port"`
			} `json:"server"`
			Connection struct {
				AutoStart   bool   `json:"auto_start"`
				AutoStartID string `json:"auto_start_id"`
			} `json:"connection"`
			Clients struct {
				VncViewer       string `json:"vnc_viewer"`
				VncPasswordFile string `json:"vnc_password_file"`
				RdpViewer       string `json:"rdp_viewer"`
			} `json:"clients"`
		}{
			Server: struct {
				Port int `json:"port"`
			}{
				Port: s.configFile.Server.Port,
			},
			Connection: struct {
				AutoStart   bool   `json:"auto_start"`
				AutoStartID string `json:"auto_start_id"`
			}{
				AutoStart:   s.configFile.Connection.AutoStart,
				AutoStartID: s.configFile.Connection.AutoStartID,
			},
			Clients: struct {
				VncViewer       string `json:"vnc_viewer"`
				VncPasswordFile string `json:"vnc_password_file"`
				RdpViewer       string `json:"rdp_viewer"`
			}{
				VncViewer:       s.configFile.Clients.VncViewer,
				VncPasswordFile: s.configFile.Clients.VncPasswordFile,
				RdpViewer:       s.configFile.Clients.RdpViewer,
			},
		}

		s.respondWithJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Data:    safeConfig,
		})

	case http.MethodPut:
		// Update configuration
		var configUpdate struct {
			Server struct {
				Port int `json:"port"`
			} `json:"server"`
			Connection struct {
				AutoStart   bool   `json:"auto_start"`
				AutoStartID string `json:"auto_start_id"`
			} `json:"connection"`
			Clients struct {
				VncViewer       string `json:"vnc_viewer"`
				VncPasswordFile string `json:"vnc_password_file"`
				RdpViewer       string `json:"rdp_viewer"`
			} `json:"clients"`
		}

		if err := json.NewDecoder(r.Body).Decode(&configUpdate); err != nil {
			s.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Update configuration
		if err := s.configFile.UpdateConfig(
			configUpdate.Server.Port,
			configUpdate.Connection.AutoStart,
			configUpdate.Connection.AutoStartID,
			configUpdate.Clients.VncViewer,
			configUpdate.Clients.VncPasswordFile,
			configUpdate.Clients.RdpViewer,
		); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		s.respondWithJSON(w, http.StatusOK, APIResponse{
			Success: true,
		})

	default:
		s.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleConnect handles connection requests
func (s *Server) handleConnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract device ID from URL
	id := strings.TrimPrefix(r.URL.Path, "/connect/")
	if id == "" {
		s.respondWithError(w, http.StatusBadRequest, "Device ID is required")
		return
	}

	// Get device
	device, found := s.configFile.GetDevice(id)
	if !found {
		s.respondWithError(w, http.StatusNotFound, "Device not found")
		return
	}

	// Start connection in background
	go s.ConnectToDevice(device)

	s.respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
	})
}

// handleDisconnect handles disconnect requests
func (s *Server) handleDisconnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Disconnect current device
	s.DisconnectCurrentDevice()

	s.respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
	})
}

// ConnectToDevice connects to a device
func (s *Server) ConnectToDevice(device config.Device) {
	s.cmdLock.Lock()
	defer s.cmdLock.Unlock()

	// First disconnect any current connection
	if s.currentCmd != nil && s.currentCmd.Process != nil {
		logger.Info("Killing current connection process").Send()
		err := s.currentCmd.Process.Kill()
		if err != nil {
			logger.Error("Failed to kill process").WithError(err).Send()
		}
		s.currentCmd.Wait()
		s.currentCmd = nil
		s.currentDeviceId = ""
	}

	var cmd *exec.Cmd
	switch device.Protocol {
	case "vnc":
		args := []string{device.IPAddress}
		if device.Port != 0 {
			args[0] = fmt.Sprintf("%s:%d", device.IPAddress, device.Port)
		} else {
			args[0] = fmt.Sprintf("%s:5900", device.IPAddress)
		}

		if device.FullScreen {
			args = append(args, "-FullScreen")
		}

		args = append(args, "-PasswordFile", s.configFile.Clients.VncPasswordFile)

		logger.Debug("VNC command arguments").WithField("args", args).Send()

		cmd = exec.Command(s.configFile.Clients.VncViewer, args...)

	case "rdp":
		args := []string{"-u", device.Username}

		//if device.Password != "" {
		//	// Decrypt password
		//	password, err := s.passwordManager.Decrypt(device.Password)
		//	if err != nil {
		//		logger.Error("Failed to decrypt password").WithError(err).Send()
		//		return
		//	}
		//
		//	args = append(args, "-p", password)
		//}

		if device.FullScreen {
			args = append(args, "-f")
		}

		args = append(args, device.IPAddress)
		if device.Port != 0 {
			args[len(args)-1] = fmt.Sprintf("%s:%d", device.IPAddress, device.Port)
		}

		logger.Debug("RDP command arguments").WithField("args", args).Send()

		cmd = exec.Command(s.configFile.Clients.RdpViewer, args...)

	default:
		logger.Error("Unknown protocol").WithField("protocol", device.Protocol).Send()
		return
	}

	logger.Info("Connecting to device").
		WithField("name", device.Name).
		WithField("id", device.ID).
		WithField("ip", device.IPAddress).
		WithField("protocol", device.Protocol).
		Send()

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		logger.Error("Failed to start command").WithError(err).Send()
		return
	}

	s.currentCmd = cmd
	s.currentDeviceId = device.ID

	go func() {
		err := cmd.Wait()
		if err != nil {
			logger.Error("Command exited with error").WithError(err).Send()
		} else {
			logger.Info("Connection process exited").Send()
		}

		s.cmdLock.Lock()
		defer s.cmdLock.Unlock()

		if s.currentCmd == cmd {
			s.currentCmd = nil
			s.currentDeviceId = ""
		}
	}()
}

// DisconnectCurrentDevice disconnects the current device if any
func (s *Server) DisconnectCurrentDevice() {
	s.cmdLock.Lock()
	defer s.cmdLock.Unlock()

	if s.currentCmd != nil && s.currentCmd.Process != nil {
		logger.Info("Disconnecting current device").Send()
		err := s.currentCmd.Process.Kill()
		if err != nil {
			logger.Error("Failed to kill process").WithError(err).Send()
		}
		s.currentCmd.Wait()
		s.currentCmd = nil
		s.currentDeviceId = ""
	}
}

// Helper functions for JSON response handling
func (s *Server) respondWithError(w http.ResponseWriter, code int, message string) {
	s.respondWithJSON(w, code, APIResponse{
		Success: false,
		Error:   message,
	})
}

func (s *Server) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Error marshaling JSON").WithError(err).Send()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
}

// logMiddleware logs incoming HTTP requests
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Log the request
		duration := time.Since(start)
		logger.Info("HTTP Request").
			WithField("method", r.Method).
			WithField("path", r.URL.Path).
			WithField("status", rw.statusCode).
			WithField("duration", duration.Milliseconds()).
			WithField("ip", r.RemoteAddr).
			Send()
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
