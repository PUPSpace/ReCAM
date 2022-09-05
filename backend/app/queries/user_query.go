package queries

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kaleemubarok/recam/backend/app/models"
)

// UserQueries struct for queries from User model.
type UserQueries struct {
	*sqlx.DB
}

// GetUserByID query for getting one User by given ID.
func (q *UserQueries) GetUserByID(id uuid.UUID) (models.User, error) {
	// Define User variable.
	user := models.User{}

	// Define query string.
	query := `SELECT * FROM t_user WHERE id = $1`

	// Send query to database.
	err := q.Get(&user, query, id)
	if err != nil {
		// Return empty object and error.
		return user, err
	}

	// Return query result.
	return user, nil
}

// GetUserByName query for getting one User by given Name.
func (q *UserQueries) GetUserByName(name string) (models.User, error) {
	// Define User variable.
	user := models.User{}

	// Define query string.
	query := `SELECT * FROM t_user WHERE name = $1`

	// Send query to database.
	err := q.Get(&user, query, name)
	if err != nil {
		// Return empty object and error.
		return user, err
	}

	// Return query result.
	return user, nil
}

// CreateUser query for creating a new user by given name and password hash.
func (q *UserQueries) CreateUser(u *models.User) error {
	// Define query string.
	query := `INSERT INTO t_user VALUES ($1, $2, $3, $4, $5, $6, $7)`

	// Send query to database.
	_, err := q.Exec(
		query,
		u.ID, u.CreatedAt, u.UpdatedAt, u.Name, u.PasswordHash, u.UserStatus, u.UserRole,
	)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}
