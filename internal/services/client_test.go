package services

import (
	"project-offer/internal/models"
	"project-offer/internal/testutil"
	"testing"
	"time"
)

func TestClientService(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	service := NewClientService(db)

	t.Run("CreateClient", func(t *testing.T) {
		client := &models.Client{
			Name:    "Test Client",
			Email:   "client@test.com",
			Address: "123 Test St\nTest City\nRomania",
		}

		err := service.CreateClient(client)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		if client.ID == 0 {
			t.Error("Expected client ID to be set after creation")
		}
		if client.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}
	})

	t.Run("GetClients", func(t *testing.T) {
		clients, err := service.GetClients()
		if err != nil {
			t.Fatalf("Failed to get clients: %v", err)
		}

		if len(clients) == 0 {
			t.Error("Expected to get test clients from database")
		}

		foundTestClient := false
		for _, c := range clients {
			if c.Email == "client@test.com" {
				foundTestClient = true
				if c.Name != "Test Client" {
					t.Errorf("Expected name %s, got %s", "Test Client", c.Name)
				}
				if c.Address != "123 Test St\nTest City\nRomania" {
					t.Errorf("Expected address not matching")
				}
			}
		}

		if !foundTestClient {
			t.Error("Did not find created test client in results")
		}
	})

	t.Run("GetClient", func(t *testing.T) {
		clients, _ := service.GetClients()
		var testClientID int64
		for _, c := range clients {
			if c.Email == "client@test.com" {
				testClientID = c.ID
				break
			}
		}

		client, err := service.GetClient(testClientID)
		if err != nil {
			t.Fatalf("Failed to get client: %v", err)
		}

		if client.Email != "client@test.com" {
			t.Errorf("Expected email %s, got %s", "client@test.com", client.Email)
		}
	})

	t.Run("UpdateClient", func(t *testing.T) {
		clients, _ := service.GetClients()
		var testClient *models.Client
		for _, c := range clients {
			if c.Email == "client@test.com" {
				testClient = &c
				break
			}
		}

		testClient.Name = "Updated Client"
		testClient.Address = "456 Update St\nNew City\nRomania"
		testClient.UpdatedAt = time.Now()

		err := service.UpdateClient(testClient)
		if err != nil {
			t.Fatalf("Failed to update client: %v", err)
		}

		// Verify update
		updated, err := service.GetClient(testClient.ID)
		if err != nil {
			t.Fatalf("Failed to get updated client: %v", err)
		}

		if updated.Name != "Updated Client" {
			t.Errorf("Expected name %s, got %s", "Updated Client", updated.Name)
		}
		if updated.Address != "456 Update St\nNew City\nRomania" {
			t.Error("Address not updated correctly")
		}
	})

	t.Run("DeleteClient", func(t *testing.T) {
		clients, _ := service.GetClients()
		var testClientID int64
		for _, c := range clients {
			if c.Email == "client@test.com" {
				testClientID = c.ID
				break
			}
		}

		err := service.DeleteClient(testClientID)
		if err != nil {
			t.Fatalf("Failed to delete client: %v", err)
		}

		// Verify deletion
		_, err = service.GetClient(testClientID)
		if err == nil {
			t.Error("Expected error when getting deleted client")
		}
	})
}