package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

type DocumentService struct {
	templates   *template.Template
	api2pdfKey  string
	skribbleKey string
}

func NewDocumentService(templates *template.Template) *DocumentService {
	return &DocumentService{
		templates:   templates,
		api2pdfKey:  os.Getenv("API2PDF_KEY"),
		skribbleKey: os.Getenv("SKRIBBLE_API_KEY"),
	}
}

// GeneratePDF generates a PDF from an offer using API2PDF
func (s *DocumentService) GeneratePDF(offerData interface{}, markdownContent string) (string, error) {
	// Convert markdown to HTML
	html := s.convertMarkdownToHTML(markdownContent)

	// Render the offer template with the offer data and requirements HTML
	var buf bytes.Buffer
	data := struct {
		Offer            interface{}
		RequirementsHTML template.HTML
	}{
		Offer:            offerData,
		RequirementsHTML: template.HTML(html),
	}

	if err := s.templates.ExecuteTemplate(&buf, "offer.html", data); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	// Call API2PDF to generate PDF
	pdfURL, err := s.callAPI2PDF(buf.String())
	if err != nil {
		return "", fmt.Errorf("PDF generation failed: %w", err)
	}

	return pdfURL, nil
}

// UploadToSkribble uploads a PDF to Skribble for e-signing
func (s *DocumentService) UploadToSkribble(pdfURL string, signerEmail string) (string, error) {
	payload := map[string]interface{}{
		"document_url": pdfURL,
		"signer": map[string]string{
			"email": signerEmail,
		},
		"message": "Please review and sign the project offer",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.skribble.com/v1/documents", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+s.skribbleKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		SigningURL string `json:"signing_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.SigningURL, nil
}

func (s *DocumentService) convertMarkdownToHTML(md string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	html := markdown.ToHTML([]byte(md), p, nil)
	return string(html)
}

func (s *DocumentService) callAPI2PDF(html string) (string, error) {
	payload := map[string]interface{}{
		"html": html,
		"options": map[string]interface{}{
			"margin": "20mm",
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.api2pdf.com/v2/html", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", s.api2pdfKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.URL, nil
}
