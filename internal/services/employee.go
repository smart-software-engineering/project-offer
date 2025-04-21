package services

import (
	"database/sql"
	"project-offer/internal/models"
)

type EmployeeService struct {
	db *sql.DB
}

func NewEmployeeService(db *sql.DB) *EmployeeService {
	return &EmployeeService{db: db}
}

func (s *EmployeeService) CreateEmployee(emp *models.Employee) error {
	query := `
		INSERT INTO employees (name, email, role, salary)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	return s.db.QueryRow(query, emp.Name, emp.Email, emp.Role, emp.YearlySalary).Scan(&emp.ID)
}

func (s *EmployeeService) GetEmployees() ([]models.Employee, error) {
	rows, err := s.db.Query(`SELECT id, name, email, role, salary FROM employees`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var emp models.Employee
		if err := rows.Scan(&emp.ID, &emp.Name, &emp.Email, &emp.Role, &emp.YearlySalary); err != nil {
			return nil, err
		}
		employees = append(employees, emp)
	}
	return employees, nil
}

func (s *EmployeeService) UpdateEmployee(emp *models.Employee) error {
	_, err := s.db.Exec(`
		UPDATE employees 
		SET name = $1, email = $2, role = $3, salary = $4
		WHERE id = $5`,
		emp.Name, emp.Email, emp.Role, emp.YearlySalary, emp.ID)
	return err
}

func (s *EmployeeService) DeleteEmployee(id int64) error {
	_, err := s.db.Exec(`DELETE FROM employees WHERE id = $1`, id)
	return err
}
