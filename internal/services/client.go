package services

import (
	"database/sql"
	"project-offer/internal/models"
	"time"
)

type ClientService struct {
	db *sql.DB
}

func NewClientService(db *sql.DB) *ClientService {
	return &ClientService{db: db}
}

func (s *ClientService) CreateClient(client *models.Client) error {
	query := `
		INSERT INTO clients (name, email, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id`

	now := time.Now()
	return s.db.QueryRow(query, client.Name, client.Email, client.Address, now).Scan(&client.ID)
}

func (s *ClientService) GetClients() ([]models.Client, error) {
	rows, err := s.db.Query(`
		SELECT id, name, email, address, created_at 
		FROM clients 
		ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var client models.Client
		if err := rows.Scan(&client.ID, &client.Name, &client.Email, &client.Address, &client.CreatedAt); err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func (s *ClientService) GetClient(id int64) (*models.Client, error) {
	client := &models.Client{}
	err := s.db.QueryRow(`
		SELECT id, name, email, address, created_at 
		FROM clients 
		WHERE id = $1`, id).Scan(
		&client.ID, &client.Name, &client.Email, &client.Address, &client.CreatedAt)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *ClientService) UpdateClient(client *models.Client) error {
	_, err := s.db.Exec(`
		UPDATE clients 
		SET name = $1, email = $2, address = $3, updated_at = $4
		WHERE id = $5`,
		client.Name, client.Email, client.Address, time.Now(), client.ID)
	return err
}

func (s *ClientService) DeleteClient(id int64) error {
	_, err := s.db.Exec(`DELETE FROM clients WHERE id = $1`, id)
	return err
}
