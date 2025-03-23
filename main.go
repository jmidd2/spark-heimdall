package main

import (
	"embed"
	"html/template"
	"log"
	configuration "spark-heimdall/internal/config"
	"spark-heimdall/internal/heimdall"
)

// Embed the templates directory
//
//go:embed templates
var templateFS embed.FS

func main() {
	configFile, err := configuration.LoadConfigFromFlags()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Parse HTML templates
	templates, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse HTML templates: %v", err)
	}

	server := heimdall.NewServer(configFile, templates)
	server.SetupRoutes()

	err = server.Start()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
