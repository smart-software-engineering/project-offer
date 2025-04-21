package services

import (
	"database/sql"
	"time"

	"project-offer/internal/models"
)

// OfferService handles business logic for offers
type OfferService struct {
	db *sql.DB
}

// NewOfferService creates a new offer service
func NewOfferService(db *sql.DB) *OfferService {
	return &OfferService{db: db}
}

// CreateOffer creates a new offer and calculates pricing
func (s *OfferService) CreateOffer(offer *models.Offer) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Set default values
	offer.Status = "draft"
	offer.CreatedAt = time.Now()
	offer.UpdatedAt = time.Now()

	// Insert offer
	err = s.insertOffer(tx, offer)
	if err != nil {
		return err
	}

	// Insert employee assignments
	if len(offer.Employees) > 0 {
		err = s.assignEmployeesToOffer(tx, offer.ID, offer.Employees)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// CalculateOfferPrice calculates the total price for an offer
func (s *OfferService) CalculateOfferPrice(offer *models.Offer) (float64, error) {
	var totalCost float64

	// Get employee costs
	rows, err := s.db.Query(`
		SELECT e.salary, ec.cost_per_year 
		FROM employees e
		CROSS JOIN employee_costs ec
		WHERE e.id = ANY($1)
	`, offer.Employees)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var salary, costs float64
		if err := rows.Scan(&salary, &costs); err != nil {
			return 0, err
		}
		// Calculate per-project cost: ((yearly salary + yearly costs) / 24) * timeframe
		totalCost += ((salary + costs) / 24) * float64(offer.TimeFrame)
	}

	// Apply risk factor
	totalCost *= offer.RiskFactor

	// Apply discount if any
	if offer.Discount != nil {
		totalCost -= offer.Discount.Amount
	}

	return totalCost, nil
}

func (s *OfferService) insertOffer(tx *sql.Tx, offer *models.Offer) error {
	query := `
		INSERT INTO offers (
			client_id, timeframe, requirements, risk_multiplier, 
			discount_amount, discount_explanation, status, 
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	var discountAmount sql.NullFloat64
	var discountExplanation sql.NullString
	if offer.Discount != nil {
		discountAmount.Float64 = offer.Discount.Amount
		discountAmount.Valid = true
		discountExplanation.String = offer.Discount.Explanation
		discountExplanation.Valid = true
	}

	return tx.QueryRow(
		query,
		offer.ClientID,
		offer.TimeFrame,
		offer.RequirementsURL,
		offer.RiskFactor,
		discountAmount,
		discountExplanation,
		offer.Status,
		offer.CreatedAt,
		offer.UpdatedAt,
	).Scan(&offer.ID)
}

func (s *OfferService) assignEmployeesToOffer(tx *sql.Tx, offerID int64, employeeIDs []int64) error {
	for _, empID := range employeeIDs {
		_, err := tx.Exec(
			"INSERT INTO offer_employees (offer_id, employee_id) VALUES ($1, $2)",
			offerID, empID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
