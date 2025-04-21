package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"project-offer/internal/handlers"

	_ "github.com/lib/pq"
)

const (
	defaultPort = "8080"
	maxFileSize = 1 << 20 // 1MB
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize database connection
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize templates
	tmpl, err := template.ParseGlob(filepath.Join("templates", "*.html"))
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	// Create handler config
	config := &handlers.Config{
		MaxFileSize:    maxFileSize,
		TemplatesPath:  "templates",
		SkribbleAPIKey: os.Getenv("SKRIBBLE_API_KEY"),
		API2PDFKey:     os.Getenv("API2PDF_KEY"),
	}

	// Initialize handler with dependencies
	h := handlers.NewHandler(db, tmpl, config)

	// Initialize router with handlers
	router := setupRoutes(h)

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func initDB() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/project_offer?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func setupRoutes(h *handlers.Handler) http.Handler {
	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// API routes
	mux.HandleFunc("/api/employees", h.HandleEmployeeList)
	mux.HandleFunc("/api/employees/create", h.HandleEmployeeCreate)
	mux.HandleFunc("/api/clients", h.HandleClientList)
	mux.HandleFunc("/api/clients/create", h.HandleClientCreate)
	mux.HandleFunc("/api/offers", h.HandleOfferList)
	mux.HandleFunc("/api/offers/create", h.HandleOfferCreate)

	// Page routes
	mux.HandleFunc("/", handleHome(h))
	mux.HandleFunc("/offers/new", handleNewOffer(h))
	mux.HandleFunc("/offers/", handleViewOffer(h))
	mux.HandleFunc("/employees", handleEmployees(h))

	return mux
}

// Page handler factories
func handleHome(h *handlers.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		h.Tmpl.ExecuteTemplate(w, "layout", nil)
	}
}

func handleNewOffer(h *handlers.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Tmpl.ExecuteTemplate(w, "new-offer", nil)
	}
}

func handleViewOffer(h *handlers.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement offer viewing
		http.NotFound(w, r)
	}
}

func handleEmployees(h *handlers.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Tmpl.ExecuteTemplate(w, "employees", nil)
	}
}
