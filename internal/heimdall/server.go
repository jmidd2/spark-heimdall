package heimdall

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	configuration "spark-heimdall/internal/config"
	"spark-heimdall/internal/device"
	"strconv"
	"sync"
)

type Server struct {
	configFile      *configuration.Config
	templates       *template.Template
	cmdLock         sync.Mutex
	currentCmd      *exec.Cmd
	currentDeviceId string
	Store           *device.Store
}

func NewServer(configFile *configuration.Config, templates *template.Template) *Server {
	return &Server{
		configFile: configFile,
		templates:  templates,
		Store:      &device.Store{Devices: configFile.Devices},
	}
}

func (s *Server) SetupRoutes() {
	http.HandleFunc("/", loggingMiddleware(s.HandleIndex))
	http.HandleFunc("/connect/", loggingMiddleware(s.HandleConnect))
	http.HandleFunc("/disconnect", loggingMiddleware(s.HandleDisconnect))
	http.HandleFunc("/api/pcs", loggingMiddleware(s.HandleGetPCs))
	http.HandleFunc("/api/pcs/add", loggingMiddleware(s.HandleAddPC))
	http.HandleFunc("/api/pcs/edit", loggingMiddleware(s.HandleEditPC))
	http.HandleFunc("/api/pcs/delete", loggingMiddleware(s.HandleDeletePC))
	http.HandleFunc("/api/config", loggingMiddleware(s.HandleGetConfig))
	http.HandleFunc("/api/config/update", loggingMiddleware(s.HandleUpdateConfig))

	// Serve static files (CSS, JS) if they exist
	if _, err := os.Stat("static"); !os.IsNotExist(err) {
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	}
}

func (s *Server) Start() error {
	// Auto-start if configured
	log.Println("Starting server...")
	if s.configFile.AutoStart && s.configFile.AutoStartID != "" {
		log.Printf("Auto-starting ID %s", s.configFile.AutoStartID)
		for _, pc := range s.Store.Devices {
			if pc.ID == s.configFile.AutoStartID {
				go s.connectToPC(pc)
				break
			}
		}
	}

	log.Printf("Starting server on port %d", s.configFile.ListenPort)
	log.Printf("Open http://localhost:%d in your browser", s.configFile.ListenPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.configFile.ListenPort), nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	return nil
}

func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := struct {
		PCs              device.Devices
		CurrentlyPlaying string
	}{
		PCs:              s.configFile.Store.Devices,
		CurrentlyPlaying: s.currentDeviceId,
	}

	w.Header().Set("Content-Type", "text/html")
	err := s.templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (s *Server) HandleConnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/connect/"):]

	for _, d := range s.Store.Devices {
		if d.ID == id {
			go s.connectToPC(d)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	http.Error(w, "PC not found", http.StatusNotFound)
}

func (s *Server) HandleDisconnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	s.disconnectCurrentPC()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) HandleGetPCs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.Store.Devices)
}

func (s *Server) HandleAddPC(w http.ResponseWriter, r *http.Request) {
	log.Printf("Route: %s, Method: %s\n", r.URL.Path, r.Method)
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var d device.Device
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		log.Printf("Error decoding PC: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Adding new device...")
	err = s.configFile.AddDevice(d)
	if err != nil {
		log.Printf("Error adding device: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.Store.Add(d)
	if err != nil {
		log.Printf("Error adding device: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(d)
	if err != nil {
		log.Printf("Error encoding device: %v", err)
		return
	}
}
func (s *Server) HandleEditPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var d device.Device
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.configFile.UpdateDevice(d)
	if err != nil {
		log.Printf("Error updating device: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}
func (s *Server) HandleDeletePC(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		ID string `json:"id"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.configFile.DeleteDevice(data.ID)
	if err != nil {
		log.Printf("Error deleting device: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
}
func (s *Server) HandleGetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Create a copy of the config without sensitive data
	safeConfig := SafeEncodeConfig{
		ListenPort:    s.configFile.ListenPort,
		AutoStart:     s.configFile.AutoStart,
		AutoStartID:   s.configFile.AutoStartID,
		VncViewer:     s.configFile.VncViewer,
		VncPasswdFile: s.configFile.VncPasswordFile,
		RdpViewer:     s.configFile.RdpViewer,
	}

	json.NewEncoder(w).Encode(safeConfig)
}

type SafeEncodeConfig struct {
	ListenPort    int    `json:"listen_port"`
	AutoStart     bool   `json:"auto_start"`
	AutoStartID   string `json:"auto_start_id"`
	VncViewer     string `json:"vnc_viewer"`
	VncPasswdFile string `json:"vnc_passwd_file"`
	RdpViewer     string `json:"rdp_viewer"`
}

type SafeDecodeConfig struct {
	ListenPort    string `json:"listen_port"`
	AutoStart     bool   `json:"auto_start"`
	AutoStartID   string `json:"auto_start_id"`
	VncViewer     string `json:"vnc_viewer"`
	VncPasswdFile string `json:"vnc_passwd_file"`
	RdpViewer     string `json:"rdp_viewer"`
}

func (s *Server) HandleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var decodedConfig SafeDecodeConfig
	err := json.NewDecoder(r.Body).Decode(&decodedConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var newConfig configuration.UpdateConfig

	// Update only allowed fields
	if decodedConfig.ListenPort != "" {
		addr, err := strconv.ParseInt(decodedConfig.ListenPort, 10, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newConfig.ListenPort = int(addr)
	}

	newConfig.AutoStart = decodedConfig.AutoStart
	newConfig.AutoStartID = decodedConfig.AutoStartID
	newConfig.VncViewer = decodedConfig.VncViewer
	newConfig.RdpViewer = decodedConfig.RdpViewer
	newConfig.VncPasswordFile = decodedConfig.VncPasswdFile

	err = s.configFile.Update(newConfig)
	if err != nil {
		log.Printf("Error updating config: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
}

func (s *Server) connectToPC(pc device.Device) {
	s.cmdLock.Lock()
	defer s.cmdLock.Unlock()

	// First disconnect any current connection
	if s.currentCmd != nil && s.currentCmd.Process != nil {
		log.Printf("Killing current process")
		err := s.currentCmd.Process.Kill()
		if err != nil {
			log.Printf("Failed to kill process: %v", err)
		}
		s.currentCmd.Wait()
		s.currentCmd = nil
	}

	var cmd *exec.Cmd
	switch pc.Protocol {
	case "vnc":
		args := []string{pc.IPAddress}
		if pc.Port != 0 {
			args[0] = fmt.Sprintf("%s:%d", pc.IPAddress, pc.Port)
		} else {
			args[0] = fmt.Sprintf("%s:5900", pc.IPAddress)
		}

		if pc.FullScreen {
			args = append(args, "-FullScreen")
		}

		args = append(args, "-PasswordFile", s.configFile.VncPasswordFile)

		log.Println(args)

		cmd = exec.Command(s.configFile.VncViewer, args...)
	case "rdp":
		args := []string{"-u", pc.Username}

		if pc.Password != "" {
			args = append(args, "-p", pc.Password)
		}

		if pc.FullScreen {
			args = append(args, "-f")
		}

		args = append(args, pc.IPAddress)
		if pc.Port != 0 {
			args[len(args)-1] = fmt.Sprintf("%s:%d", pc.IPAddress, pc.Port)
		}

		cmd = exec.Command(s.configFile.RdpViewer, args...)
	default:
		log.Printf("Unknown protocol: %s", pc.Protocol)
		return
	}

	log.Printf("Running command: %v %v", cmd.Path, cmd.Args)
	log.Printf("Connecting to %s (%s)", pc.Name, pc.IPAddress)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Printf("Failed to start command: %v", err)
		return
	}

	s.currentCmd = cmd
	s.currentDeviceId = pc.ID

	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Printf("Command exited with error: %v", err)
			s.currentCmd = nil
			s.currentDeviceId = ""
		}

		s.cmdLock.Lock()
		defer s.cmdLock.Unlock()

		if s.currentCmd == cmd {
			s.currentCmd = nil
			s.currentDeviceId = ""
		}
	}()
}

func (s *Server) disconnectCurrentPC() {
	s.cmdLock.Lock()
	defer s.cmdLock.Unlock()

	if s.currentCmd != nil && s.currentCmd.Process != nil {
		log.Printf("Disconnecting current PC")
		err := s.currentCmd.Process.Kill()
		if err != nil {
			log.Printf("Failed to kill process: %v", err)
		}
		s.currentCmd.Wait()
		s.currentCmd = nil
		s.currentDeviceId = ""
	}
}
