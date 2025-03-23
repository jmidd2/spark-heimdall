package heimdall

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// ANSI color codes
const (
	Reset      = "\033[0m"
	Red        = "\033[31m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Purple     = "\033[35m"
	Cyan       = "\033[36m"
	White      = "\033[37m"
	BoldRed    = "\033[1;31m"
	BoldGreen  = "\033[1;32m"
	BoldYellow = "\033[1;33m"
	BoldBlue   = "\033[1;34m"
)

// Custom response writer to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

// Enhanced logging middleware with emojis and colors
func loggingMiddleware(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Get a response writer that captures the status code
		rw := newResponseWriter(w)

		// Call the next handler
		next(rw, r)

		// Calculate request duration
		duration := time.Since(startTime)

		// Choose emoji based on status code
		var emoji string
		var color string
		switch {
		case rw.statusCode >= 500:
			emoji = "❌ "
			color = BoldRed
		case rw.statusCode >= 400:
			emoji = "⚠️ "
			color = Yellow
		case rw.statusCode >= 300:
			emoji = "🔄 "
			color = Cyan
		case rw.statusCode >= 200:
			emoji = "✅ "
			color = Green
		default:
			emoji = "❓ "
			color = Blue
		}

		// Choose HTTP method color
		var methodColor string
		switch r.Method {
		case "GET":
			methodColor = Green
		case "POST":
			methodColor = Blue
		case "PUT":
			methodColor = Yellow
		case "DELETE":
			methodColor = Red
		case "PATCH":
			methodColor = Cyan
		default:
			methodColor = White
		}

		// Choose path color based on path components
		pathColor := Purple
		if strings.HasPrefix(r.URL.Path, "/api") {
			pathColor = BoldBlue
		} else if strings.HasPrefix(r.URL.Path, "/admin") {
			pathColor = BoldRed
		} else if strings.HasPrefix(r.URL.Path, "/static") {
			pathColor = BoldYellow
		}

		// Format the colored output first
		coloredOutput := fmt.Sprintf("%s[%s%s%s] %s%s%s %s%d%s %s%.2fms%s",
			emoji,
			methodColor, r.Method, Reset,
			pathColor, r.URL.Path, Reset,
			color, rw.statusCode, Reset,
			Blue, duration.Seconds()*1000, Reset,
		)

		// Use log package to include timestamp
		log.Println(coloredOutput)
	})
}
