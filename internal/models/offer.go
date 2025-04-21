package models

import "time"

// TimeFrame represents the project duration type
type TimeFrame int

const (
	TwoWeeks TimeFrame = 2
	SixWeeks TimeFrame = 6
)

// Offer represents a project offer in the system
type Offer struct {
	ID              int64     `json:"id"`
	ClientID        int64     `json:"client_id"`
	TimeFrame       TimeFrame `json:"time_frame"`
	RequirementsURL string    `json:"requirements_url"`
	RiskFactor      float64   `json:"risk_factor"`
	Discount        *Discount `json:"discount,omitempty"`
	Employees       []int64   `json:"employees"` // List of employee IDs
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	SkribbleLink    string    `json:"skribble_link,omitempty"`
}

// Discount represents a discount applied to an offer
type Discount struct {
	Amount      float64 `json:"amount"`
	Explanation string  `json:"explanation"`
}

// ValidateTimeFrame checks if the timeframe is valid
func (t TimeFrame) ValidateTimeFrame() bool {
	return t == TwoWeeks || t == SixWeeks
}
