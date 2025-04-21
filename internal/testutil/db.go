package testutil

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

const (
	testDBHost = "localhost"
	testDBPort = 5432
	testDBUser = "postgres"
	testDBPass = "postgres"
	testDBName = "project_offer_test"
)

// SetupTestDB creates a test database and returns a connection
func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// Connect to default postgres database to create test database
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
		testDBHost, testDBPort, testDBUser, testDBPass,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to postgres: %v", err)
	}

	// Drop test database if it exists and create it fresh
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
	if err != nil {
		t.Fatalf("Failed to drop test database: %v", err)
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName))
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Close connection to postgres database
	db.Close()

	// Connect to newly created test database
	testConnStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		testDBHost, testDBPort, testDBUser, testDBPass, testDBName,
	)

	testDB, err := sql.Open("postgres", testConnStr)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Initialize schema
	if err := initializeSchema(testDB); err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	return testDB
}

// initializeSchema executes all SQL files in the sql directory
func initializeSchema(db *sql.DB) error {
	sqlFiles, err := filepath.Glob("../../sql/*.sql")
	if err != nil {
		return fmt.Errorf("failed to find SQL files: %v", err)
	}

	for _, file := range sqlFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read SQL file %s: %v", file, err)
		}

		// Split the file into individual statements
		statements := strings.Split(string(content), ";")

		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			if _, err := db.Exec(stmt); err != nil {
				return fmt.Errorf("failed to execute SQL statement from %s: %v", file, err)
			}
		}
	}

	return nil
}

// CleanupTestDB drops the test database
func CleanupTestDB(t *testing.T, db *sql.DB) {
	t.Helper()

	// Close connection to test database
	if err := db.Close(); err != nil {
		log.Printf("Warning: Failed to close test database connection: %v", err)
	}

	// Connect to default postgres database to drop test database
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
		testDBHost, testDBPort, testDBUser, testDBPass,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Warning: Failed to connect to postgres: %v", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
	if err != nil {
		log.Printf("Warning: Failed to drop test database: %v", err)
	}
}
