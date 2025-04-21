package services

import (
	"project-offer/internal/models"
	"project-offer/internal/testutil"
	"testing"
	"time"
)

func TestOfferService(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	service := NewOfferService(db)

	// Helper function to create a test client
	createTestClient := func(t *testing.T) int64 {
		clientService := NewClientService(db)
		client := &models.Client{
			Name:    "Test Client",
			Email:   "test@client.com",
			Address: "Test Address",
		}
		if err := clientService.CreateClient(client); err != nil {
			t.Fatalf("Failed to create test client: %v", err)
		}
		return client.ID
	}

	// Helper function to create test employees
	createTestEmployees := func(t *testing.T) []int64 {
		employeeService := NewEmployeeService(db)
		employees := []models.Employee{
			{
				Name:         "Senior Dev",
				Email:        "senior@test.com",
				Role:         models.RoleSenior,
				YearlySalary: 60000,
			},
			{
				Name:         "Junior Dev",
				Email:        "junior@test.com",
				Role:         models.RoleJunior,
				YearlySalary: 30000,
			},
		}

		var ids []int64
		for _, emp := range employees {
			emp := emp // Create new variable for each iteration
			if err := employeeService.CreateEmployee(&emp); err != nil {
				t.Fatalf("Failed to create test employee: %v", err)
			}
			ids = append(ids, emp.ID)
		}
		return ids
	}

	t.Run("CreateOffer", func(t *testing.T) {
		clientID := createTestClient(t)
		employeeIDs := createTestEmployees(t)

		offer := &models.Offer{
			ClientID:        clientID,
			TimeFrame:       models.TwoWeeks,
			RequirementsURL: "test-requirements.md",
			RiskFactor:     1.5,
			Employees:      employeeIDs,
			Discount: &models.Discount{
				Amount:      1000,
				Explanation: "Test discount",
			},
		}

		err := service.CreateOffer(offer)
		if err != nil {
			t.Fatalf("Failed to create offer: %v", err)
		}

		if offer.ID == 0 {
			t.Error("Expected offer ID to be set after creation")
		}
		if offer.Status != "draft" {
			t.Errorf("Expected status 'draft', got '%s'", offer.Status)
		}
	})

	t.Run("CalculateOfferPrice", func(t *testing.T) {
		clientID := createTestClient(t)
		employeeIDs := createTestEmployees(t)

		offer := &models.Offer{
			ClientID:        clientID,
			TimeFrame:       models.TwoWeeks,
			RequirementsURL: "test-requirements.md",
			RiskFactor:     1.5,
			Employees:      employeeIDs,
		}

		price, err := service.CalculateOfferPrice(offer)
		if err != nil {
			t.Fatalf("Failed to calculate offer price: %v", err)
		}

		if price <= 0 {
			t.Error("Expected price to be greater than 0")
		}

		// Test with discount
		offer.Discount = &models.Discount{
			Amount:      1000,
			Explanation: "Test discount",
		}

		discountedPrice, err := service.CalculateOfferPrice(offer)
		if err != nil {
			t.Fatalf("Failed to calculate discounted offer price: %v", err)
		}

		if discountedPrice >= price {
			t.Error("Expected discounted price to be less than original price")
		}
		if price-discountedPrice != 1000 {
			t.Errorf("Expected discount of 1000, got %f", price-discountedPrice)
		}
	})

	t.Run("Invalid TimeFrame", func(t *testing.T) {
		offer := &models.Offer{
			TimeFrame: 3, // Invalid timeframe
		}

		err := service.CreateOffer(offer)
		if err == nil {
			t.Error("Expected error for invalid timeframe")
		}
	})

	t.Run("Invalid RiskFactor", func(t *testing.T) {
		offer := &models.Offer{
			TimeFrame:   models.TwoWeeks,
			RiskFactor: 0.5, // Too low
		}

		err := service.CreateOffer(offer)
		if err == nil {
			t.Error("Expected error for risk factor < 1.0")
		}

		offer.RiskFactor = 2.5 // Too high
		err = service.CreateOffer(offer)
		if err == nil {
			t.Error("Expected error for risk factor > 2.0")
		}
	})

	t.Run("Discount Validation", func(t *testing.T) {
		clientID := createTestClient(t)
		offer := &models.Offer{
			ClientID:   clientID,
			TimeFrame:  models.TwoWeeks,
			Discount: &models.Discount{
				Amount: 1000,
				// Missing explanation
			},
		}

		err := service.CreateOffer(offer)
		if err == nil {
			t.Error("Expected error for discount without explanation")
		}
	})
}

func TestOfferIntegration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	employeeService := NewEmployeeService(db)
	clientService := NewClientService(db)
	offerService := NewOfferService(db)

	// Create test data
	client := &models.Client{
		Name:    "Integration Test Client",
		Email:   "integration@test.com",
		Address: "Integration Test Address",
	}
	if err := clientService.CreateClient(client); err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}

	emp1 := &models.Employee{
		Name:         "Senior Integration",
		Email:        "senior.integration@test.com",
		Role:         models.RoleSenior,
		YearlySalary: 65000,
	}
	if err := employeeService.CreateEmployee(emp1); err != nil {
		t.Fatalf("Failed to create first test employee: %v", err)
	}

	emp2 := &models.Employee{
		Name:         "Junior Integration",
		Email:        "junior.integration@test.com",
		Role:         models.RoleJunior,
		YearlySalary: 35000,
	}
	if err := employeeService.CreateEmployee(emp2); err != nil {
		t.Fatalf("Failed to create second test employee: %v", err)
	}

	// Create and test offer
	offer := &models.Offer{
		ClientID:        client.ID,
		TimeFrame:       models.SixWeeks,
		RequirementsURL: "integration-test.md",
		RiskFactor:     1.8,
		Employees:      []int64{emp1.ID, emp2.ID},
		Discount: &models.Discount{
			Amount:      2000,
			Explanation: "Integration test discount",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := offerService.CreateOffer(offer); err != nil {
		t.Fatalf("Failed to create integration test offer: %v", err)
	}

	price, err := offerService.CalculateOfferPrice(offer)
	if err != nil {
		t.Fatalf("Failed to calculate price for integration test offer: %v", err)
	}

	// Verify the calculations make sense
	expectedMinimum := ((65000+35000)/24)*6*1.8 - 2000 // Rough minimum without additional costs
	if price < expectedMinimum {
		t.Errorf("Price %f is lower than expected minimum %f", price, expectedMinimum)
	}
}