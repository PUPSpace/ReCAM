package database

import (
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/kaleemubarok/recam/backend/app/queries"
)

// Queries struct for collect all app queries.
type Queries struct {
	*queries.UserQueries  // load queries from User model
	*queries.BookQueries  // load queries from Book model
	*queries.RouteQueries // load queries from Route model
	*queries.LogQueries   // load queries from RouteLog model
}

// OpenDBConnection func for opening database connection.
func OpenDBConnection() (*Queries, error) {
	// Define Database connection variables.
	var (
		db  *sqlx.DB
		err error
	)

	// Get DB_TYPE value from .env file.
	dbType := os.Getenv("DB_TYPE")

	// Define a new Database connection with right DB type.
	switch dbType {
	case "pgx":
		db, err = PostgreSQLConnection()
	case "mysql":
		db, err = MysqlConnection()
	}

	if err != nil {
		return nil, err
	}

	return &Queries{
		// Set queries from models:
		UserQueries:  &queries.UserQueries{DB: db},  // from User model
		BookQueries:  &queries.BookQueries{DB: db},  // from Book model
		RouteQueries: &queries.RouteQueries{DB: db}, // from Route model
		LogQueries:   &queries.LogQueries{DB: db},   // from RouteLog model
	}, nil
}
