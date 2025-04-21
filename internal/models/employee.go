package models

// Role represents the employee role type
type Role string

const (
	RolePrincipal    Role = "Principal"
	RoleSenior       Role = "Senior"
	RoleProfessional Role = "Professional"
	RoleJunior       Role = "Junior"
)

// Employee represents an employee in the system
type Employee struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Role         Role    `json:"role"`
	YearlySalary float64 `json:"yearly_salary"`
}

// ValidateRole checks if the role is valid
func (r Role) ValidateRole() bool {
	switch r {
	case RolePrincipal, RoleSenior, RoleProfessional, RoleJunior:
		return true
	default:
		return false
	}
}
