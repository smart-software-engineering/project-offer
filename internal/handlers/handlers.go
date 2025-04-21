package handlers

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"net/http"
	"project-offer/internal/models"
	"project-offer/internal/services"
	"strconv"
)

// Handler wraps common dependencies for all handlers
type Handler struct {
	DB              *sql.DB
	Tmpl            *template.Template
	Config          *Config
	EmployeeService *services.EmployeeService
	ClientService   *services.ClientService
	OfferService    *services.OfferService
	DocumentService *services.DocumentService
}

// Config contains application configuration
type Config struct {
	MaxFileSize    int64  // Maximum file size for uploads in bytes
	TemplatesPath  string // Path to HTML templates
	SkribbleAPIKey string // API key for Skribble integration
	API2PDFKey     string // API key for API2PDF service
}

// NewHandler creates a new handler with dependencies
func NewHandler(db *sql.DB, tmpl *template.Template, cfg *Config) *Handler {
	return &Handler{
		DB:              db,
		Tmpl:            tmpl,
		Config:          cfg,
		EmployeeService: services.NewEmployeeService(db),
		ClientService:   services.NewClientService(db),
		OfferService:    services.NewOfferService(db),
		DocumentService: services.NewDocumentService(tmpl),
	}
}

// Employee handlers
func (h *Handler) HandleEmployeeList(w http.ResponseWriter, r *http.Request) {
	employees, err := h.EmployeeService.GetEmployees()
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch employees")
		return
	}
	h.respondWithHTML(w, "employee-list", employees)
}

func (h *Handler) HandleEmployeeCreate(w http.ResponseWriter, r *http.Request) {
	var emp models.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.EmployeeService.CreateEmployee(&emp); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create employee")
		return
	}

	h.respondWithHTMXTrigger(w, "employeeCreated", emp)
}

// Client handlers
func (h *Handler) HandleClientList(w http.ResponseWriter, r *http.Request) {
	clients, err := h.ClientService.GetClients()
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch clients")
		return
	}
	h.respondWithHTML(w, "client-list", clients)
}

func (h *Handler) HandleClientCreate(w http.ResponseWriter, r *http.Request) {
	var client models.Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.ClientService.CreateClient(&client); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create client")
		return
	}

	h.respondWithHTMXTrigger(w, "clientCreated", client)
}

// Offer handlers
func (h *Handler) HandleOfferList(w http.ResponseWriter, r *http.Request) {
	offers, err := h.OfferService.GetOffers()
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch offers")
		return
	}
	h.respondWithHTML(w, "offer-list", offers)
}

func (h *Handler) HandleOfferCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(h.Config.MaxFileSize); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "File too large")
		return
	}

	offer := &models.Offer{
		ClientID:   parseID(r.FormValue("client_id")),
		TimeFrame:  models.TimeFrame(parseID(r.FormValue("timeframe"))),
		RiskFactor: parseFloat(r.FormValue("risk_factor")),
	}

	// Handle discount if present
	if amount := parseFloat(r.FormValue("discount_amount")); amount > 0 {
		offer.Discount = &models.Discount{
			Amount:      amount,
			Explanation: r.FormValue("discount_explanation"),
		}
	}

	// Handle requirements file
	file, header, err := r.FormFile("requirements")
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Missing requirements file")
		return
	}
	defer file.Close()

	// TODO: Save file and set RequirementsURL

	if err := h.OfferService.CreateOffer(offer); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create offer")
		return
	}

	h.respondWithHTMXTrigger(w, "offerCreated", offer)
}

// Helper functions

func (h *Handler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (h *Handler) respondWithHTML(w http.ResponseWriter, template string, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	return h.Tmpl.ExecuteTemplate(w, template, data)
}

func (h *Handler) respondWithHTMXTrigger(w http.ResponseWriter, trigger string, data interface{}) {
	w.Header().Set("HX-Trigger", trigger)
	h.respondWithJSON(w, http.StatusOK, data)
}

func (h *Handler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

// Utility functions
func parseID(s string) int64 {
	id, _ := strconv.ParseInt(s, 10, 64)
	return id
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
