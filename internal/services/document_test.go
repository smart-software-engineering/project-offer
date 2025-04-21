package services

import (
	"html/template"
	"os"
	"strings"
	"testing"
)

func TestDocumentService(t *testing.T) {
	// Create test templates
	tmpl := template.Must(template.New("test").Parse(`
		{{define "offer.html"}}
		<h1>Test Offer</h1>
		<div>{{.Offer.Title}}</div>
		<div>{{.RequirementsHTML}}</div>
		{{end}}
	`))

	service := NewDocumentService(tmpl)

	t.Run("ConvertMarkdownToHTML", func(t *testing.T) {
		markdown := "# Test Heading\n\nThis is a test paragraph."
		html := service.convertMarkdownToHTML(markdown)

		expectedHTML := "<h1 id=\"test-heading\">Test Heading</h1>\n\n<p>This is a test paragraph.</p>\n"
		if html != expectedHTML {
			t.Errorf("Expected HTML:\n%s\nGot:\n%s", expectedHTML, html)
		}
	})

	t.Run("GeneratePDF", func(t *testing.T) {
		if os.Getenv("API2PDF_KEY") == "" {
			t.Skip("Skipping PDF generation test - API2PDF_KEY not set")
		}

		offerData := struct {
			Title string
		}{
			Title: "Test Offer",
		}

		markdownContent := "# Requirements\n\n- Requirement 1\n- Requirement 2"

		pdfURL, err := service.GeneratePDF(offerData, markdownContent)
		if err != nil {
			t.Fatalf("Failed to generate PDF: %v", err)
		}

		if !strings.HasPrefix(pdfURL, "https://") || !strings.HasSuffix(pdfURL, ".pdf") {
			t.Errorf("Invalid PDF URL format: %s", pdfURL)
		}
	})

	t.Run("UploadToSkribble", func(t *testing.T) {
		if os.Getenv("SKRIBBLE_API_KEY") == "" {
			t.Skip("Skipping Skribble upload test - SKRIBBLE_API_KEY not set")
		}

		pdfURL := "https://example.com/test.pdf"
		signerEmail := "test@example.com"

		signingURL, err := service.UploadToSkribble(pdfURL, signerEmail)
		if err != nil {
			t.Fatalf("Failed to upload to Skribble: %v", err)
		}

		if !strings.HasPrefix(signingURL, "https://") {
			t.Errorf("Invalid signing URL format: %s", signingURL)
		}
	})
}