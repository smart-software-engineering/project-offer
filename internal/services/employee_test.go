package services

import (
	"project-offer/internal/models"
	"project-offer/internal/testutil"
	"testing"
)

func TestEmployeeService(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	service := NewEmployeeService(db)

	t.Run("CreateEmployee", func(t *testing.T) {
		emp := &models.Employee{
			Name:         "Test Employee",
			Email:        "test@example.com",
			Role:         models.RoleSenior,
			YearlySalary: 60000,
		}

		err := service.CreateEmployee(emp)
		if err != nil {
			t.Fatalf("Failed to create employee: %v", err)
		}

		if emp.ID == 0 {
			t.Error("Expected employee ID to be set after creation")
		}
	})

	t.Run("GetEmployees", func(t *testing.T) {
		employees, err := service.GetEmployees()
		if err != nil {
			t.Fatalf("Failed to get employees: %v", err)
		}

		// Check test data from 02_test_data.sql
		if len(employees) == 0 {
			t.Error("Expected to get test employees from database")
		}

		foundTestEmployee := false
		for _, emp := range employees {
			if emp.Email == "test@example.com" {
				foundTestEmployee = true
				if emp.Role != models.RoleSenior {
					t.Errorf("Expected role %s, got %s", models.RoleSenior, emp.Role)
				}
				if emp.YearlySalary != 60000 {
					t.Errorf("Expected salary %f, got %f", 60000.0, emp.YearlySalary)
				}
			}
		}

		if !foundTestEmployee {
			t.Error("Did not find created test employee in results")
		}
	})

	t.Run("UpdateEmployee", func(t *testing.T) {
		emp := &models.Employee{
			Name:         "Updated Employee",
			Email:        "test@example.com",
			Role:         models.RolePrincipal,
			YearlySalary: 85000,
		}

		employees, _ := service.GetEmployees()
		for _, e := range employees {
			if e.Email == emp.Email {
				emp.ID = e.ID
				break
			}
		}

		err := service.UpdateEmployee(emp)
		if err != nil {
			t.Fatalf("Failed to update employee: %v", err)
		}

		// Verify update
		employees, _ = service.GetEmployees()
		found := false
		for _, e := range employees {
			if e.ID == emp.ID {
				found = true
				if e.Role != models.RolePrincipal {
					t.Errorf("Expected role %s, got %s", models.RolePrincipal, e.Role)
				}
				if e.YearlySalary != 85000 {
					t.Errorf("Expected salary %f, got %f", 85000.0, e.YearlySalary)
				}
			}
		}

		if !found {
			t.Error("Updated employee not found")
		}
	})

	t.Run("DeleteEmployee", func(t *testing.T) {
		employees, _ := service.GetEmployees()
		var testEmpID int64
		for _, e := range employees {
			if e.Email == "test@example.com" {
				testEmpID = e.ID
				break
			}
		}

		err := service.DeleteEmployee(testEmpID)
		if err != nil {
			t.Fatalf("Failed to delete employee: %v", err)
		}

		// Verify deletion
		employees, _ = service.GetEmployees()
		for _, e := range employees {
			if e.ID == testEmpID {
				t.Error("Employee still exists after deletion")
			}
		}
	})
}