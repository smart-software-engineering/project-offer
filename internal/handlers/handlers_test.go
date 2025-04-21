package handlers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"project-offer/internal/models"
	"project-offer/internal/testutil"
	"strings"
	"testing"
)

func setupTestHandler(t *testing.T) *Handler {
	db := testutil.SetupTestDB(t)
	t.Cleanup(func() {
		testutil.CleanupTestDB(t, db)
	})

	tmpl, err := template.ParseGlob("../../templates/*.html")
	if err != nil {
		t.Fatalf("Failed to parse templates: %v", err)
	}

	config := &Config{
		MaxFileSize:    1 << 20, // 1MB
		TemplatesPath:  "../../templates",
		SkribbleAPIKey: "test-key",
		API2PDFKey:     "test-key",
	}

	return NewHandler(db, tmpl, config)
}

func createMultipartFormRequest(t *testing.T, url string, fields map[string]string, files map[string]string) (*http.Request, string) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Add regular fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatalf("Failed to write form field: %v", err)
		}
	}

	// Add files
	for fieldName, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			t.Fatalf("Failed to open test file: %v", err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}
		if _, err := io.Copy(part, file); err != nil {
			t.Fatalf("Failed to copy file content: %v", err)
		}
	}

	writer.Close()

	req := httptest.NewRequest("POST", url, &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, writer.FormDataContentType()
}

func TestHandleEmployeeList(t *testing.T) {
	h := setupTestHandler(t)

	t.Run("List Employees", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/employees", nil)
		w := httptest.NewRecorder()

		h.HandleEmployeeList(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		if !strings.Contains(w.Body.String(), "table") {
			t.Error("Expected HTML table in response")
		}
	})
}

func TestHandleEmployeeCreate(t *testing.T) {
	h := setupTestHandler(t)

	t.Run("Create Valid Employee", func(t *testing.T) {
		emp := models.Employee{
			Name:         "Test Handler Employee",
			Email:        "handler.test@example.com",
			Role:         models.RoleSenior,
			YearlySalary: 65000,
		}

		body, err := json.Marshal(emp)
		if err != nil {
			t.Fatalf("Failed to marshal employee: %v", err)
		}

		req := httptest.NewRequest("POST", "/api/employees/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.HandleEmployeeCreate(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		if !strings.Contains(w.Header().Get("HX-Trigger"), "employeeCreated") {
			t.Error("Expected HX-Trigger header with employeeCreated event")
		}
	})

	t.Run("Create Invalid Employee", func(t *testing.T) {
		emp := models.Employee{
			Name:  "Invalid",
			Email: "not-an-email",
			Role:  "InvalidRole",
		}

		body, _ := json.Marshal(emp)
		req := httptest.NewRequest("POST", "/api/employees/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.HandleEmployeeCreate(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestHandleOfferCreate(t *testing.T) {
	h := setupTestHandler(t)

	t.Run("Create Valid Offer", func(t *testing.T) {
		// Create a temporary markdown file
		tmpFile, err := os.CreateTemp("", "requirements-*.md")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.WriteString("# Test Requirements\n\n- Requirement 1\n- Requirement 2"); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		tmpFile.Close()

		fields := map[string]string{
			"client_id":            "1",
			"timeframe":            "2",
			"risk_factor":          "1.5",
			"discount_amount":      "1000",
			"discount_explanation": "Test discount",
		}

		files := map[string]string{
			"requirements": tmpFile.Name(),
		}

		req, contentType := createMultipartFormRequest(t, "/api/offers/create", fields, files)
		w := httptest.NewRecorder()

		h.HandleOfferCreate(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		if !strings.Contains(w.Header().Get("HX-Trigger"), "offerCreated") {
			t.Error("Expected HX-Trigger header with offerCreated event")
		}
	})

	t.Run("Create Offer With Invalid File", func(t *testing.T) {
		// Create a file that's too large
		tmpFile, err := os.CreateTemp("", "large-*.md")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		// Write more than MaxFileSize bytes
		if err := tmpFile.Truncate(2 << 20); err != nil {
			t.Fatalf("Failed to create large file: %v", err)
		}
		tmpFile.Close()

		fields := map[string]string{
			"client_id":   "1",
			"timeframe":   "2",
			"risk_factor": "1.5",
		}

		files := map[string]string{
			"requirements": tmpFile.Name(),
		}

		req, contentType := createMultipartFormRequest(t, "/api/offers/create", fields, files)
		w := httptest.NewRecorder()

		h.HandleOfferCreate(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestHandleOfferList(t *testing.T) {
	h := setupTestHandler(t)

	t.Run("List Offers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/offers", nil)
		w := httptest.NewRecorder()

		h.HandleOfferList(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		if !strings.Contains(w.Body.String(), "grid") {
			t.Error("Expected offer grid in response")
		}
	})
}

func TestHandleClientList(t *testing.T) {
	h := setupTestHandler(t)

	t.Run("List Clients", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/clients", nil)
		w := httptest.NewRecorder()

		h.HandleClientList(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestHandleClientCreate(t *testing.T) {
	h := setupTestHandler(t)

	t.Run("Create Valid Client", func(t *testing.T) {
		client := models.Client{
			Name:    "Test Handler Client",
			Email:   "handler.client@example.com",
			Address: "Test Address",
		}

		body, _ := json.Marshal(client)
		req := httptest.NewRequest("POST", "/api/clients/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.HandleClientCreate(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		if !strings.Contains(w.Header().Get("HX-Trigger"), "clientCreated") {
			t.Error("Expected HX-Trigger header with clientCreated event")
		}
	})

	t.Run("Create Invalid Client", func(t *testing.T) {
		client := models.Client{
			Name:  "Invalid",
			Email: "not-an-email",
		}

		body, _ := json.Marshal(client)
		req := httptest.NewRequest("POST", "/api/clients/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.HandleClientCreate(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}