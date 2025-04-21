# GitHub Copilot Instructions for Project Offering Tool

## Technology Stack
- **Backend**: Go with standard library and key packages:
  - `database/sql` with `lib/pq` for PostgreSQL connectivity
  - Built-in `net/http` for routing
  - `html/template` for templating
- **Frontend**: HTMX for dynamic interactions
  - No JavaScript framework required
  - Server-side rendered HTML
  - HTMX attributes for dynamic updates
- **Database**: PostgreSQL
  - Version 14+ recommended
  - JSONB support for flexible data storage
  - Proper indexing strategies

## Project Structure Requirements
When scaffolding this project, use the following Go project layout:

```
project-offer/
├── cmd/
│   └── server/             # Main application entry point
├── internal/
│   ├── models/            # Database models
│   ├── handlers/          # HTTP handlers
│   ├── services/          # Business logic
│   └── templates/         # HTML templates
├── sql/                   # SQL scripts for table creation and initial data
├── static/               # Static assets
│   ├── css/
│   └── js/              # HTMX only
└── templates/           # Go HTML templates
```

1. **Database Models** (in `internal/models`):
   - Client model (name, contact info)
   - Employee model (name, role [Principal/Senior/Professional/Junior], salary)
   - Offer model (client, employees, timeframe, requirements file, multiplier/discount)

2. **Core Features**:
   - Simple HTTP endpoints with HTMX integration
   - PDF generation using Go templates
   - Skribble API integration
   - Markdown to HTML conversion

## HTMX Integration Guidelines

1. **HTMX Patterns to Use**:
   ```html
   <!-- Example for dynamic employee loading -->
   <div hx-get="/api/employees" 
        hx-trigger="load"
        hx-swap="innerHTML">
   </div>
   ```

2. **Server Response Format**:
   - Return HTML fragments for HTMX requests
   - Use proper HTTP status codes
   - Include HX-Trigger headers for client-side events

### Setting Up New Components

1. **Creating Models**: 
   ```go
   // Example model
   type Client struct {
       ID        int64     
       Name      string    
       Email     string    
       CreatedAt time.Time
   }
   ```

2. **HTML Templates with HTMX**:
   ```html
   <!-- Example template -->
   {{define "employee-form"}}
   <form hx-post="/api/employees" 
         hx-swap="afterend">
       <!-- form fields -->
   </form>
   {{end}}
   ```

3. **Go Handlers**:
   ```go
   // Example handler pattern
   func handleEmployeeCreate(w http.ResponseWriter, r *http.Request) {
       // Process form
       // Return HTML fragment
   }
   ```

### Best Practices for This Project

1. **Database Operations**:
   - Use `database/sql` with prepared statements
   - Simple SQL migrations
   - Use transactions where needed
   - Connection pooling via `sql.DB`

2. **API Design**:
   - HTML-first approach with HTMX
   - Progressive enhancement
   - Server-side validation
   - Simple error responses with HTML fragments

3. **File Operations**:
   - Handle markdown file uploads securely
   - Implement proper file size validation
   - Include error handling for PDF generation

4. **Security Considerations**:
   - CSRF protection for forms
   - Input sanitization
   - Proper PostgreSQL connection pooling
   - Secure session handling

## Testing Guidelines

When requesting test code from Copilot:
- Include both unit and integration tests
- Test HTTP handlers
- Use table-driven tests
- Test HTML responses
- Use `httptest` package for handler testing

## Database Schema Guidelines

1. **Database Setup**:
   - SQL scripts in `sql/` directory for table creation
   - Simple initialization script to check and create tables if needed
   - Direct SQL execution for schema changes during development
   - Keep SQL scripts versioned in repository

2. **PostgreSQL Features to Use**:
   - JSONB for flexible data
   - Proper indexing strategies
   - Foreign key constraints

## Validation Requirements

When implementing validation, ensure:
- Employee roles are restricted to Principal, Senior, Professional, Junior
- Timeframe is either 2 or 6 weeks
- Requirement files are markdown only, max 1MB
- Discount amounts require explanation text

## Development Setup

Required tools:
- Go 1.24 or later
- PostgreSQL 14+
- HTMX (included via CDN)

By following these instructions, you can effectively use GitHub Copilot to develop the Project Offering Tool while maintaining consistent code quality and project requirements.