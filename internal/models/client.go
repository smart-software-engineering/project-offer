package models

import "time"

// Client represents a client in the system
type Client struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Validate checks if the client data is valid
func (c *Client) Validate() error {
	// TODO: Add validation rules
	return nil
}
