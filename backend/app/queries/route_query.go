package queries

import (
	// "github.com/google/uuid"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kaleemubarok/recam/backend/app/models"
)

// RouteQueries struct for queries from Route model.
type RouteQueries struct {
	*sqlx.DB
}

// GetRoutes method for getting all routes.
func (q *RouteQueries) GetRoutes() ([]models.Route, error) {
	// Define routes variable.
	routes := []models.Route{}

	// Define query string.
	query := `SELECT * FROM t_route ORDER BY created_at DESC`

	// Send query to database.
	err := q.Select(&routes, query)
	if err != nil {
		// Return empty object and error.
		return routes, err
	}

	// Return query result.
	return routes, nil
}

/*
// GetRoutesByAuthor method for getting all routes by given author.
func (q *RouteQueries) GetRoutesByAuthor(author string) ([]models.Route, error) {
	// Define routes variable.
	routes := []models.Route{}

	// Define query string.
	query := `SELECT * FROM routes WHERE author = $1`

	// Send query to database.
	err := q.Get(&routes, query, author)
	if err != nil {
		// Return empty object and error.
		return routes, err
	}

	// Return query result.
	return routes, nil
}
*/

// GetRoute method for getting one route by given ID.
func (q *RouteQueries) GetRoute(id uuid.UUID) (models.Route, error) {
	// Define route variable.
	route := models.Route{}

	// Define query string.
	query := `SELECT * FROM t_route WHERE id = $1`

	// Send query to database.
	err := q.Get(&route, query, id)
	if err != nil {
		// Return empty object and error.
		return route, err
	}

	// Return query result.
	return route, nil
}

// GetRoute method for getting one route by given slug.
func (q *RouteQueries) GetRouteSlug(slug string) (models.Route, error) {
	// Define route variable.
	route := models.Route{}

	// Define query string.
	query := `SELECT * FROM t_route WHERE slug = $1`

	// Send query to database.
	err := q.Get(&route, query, slug)
	if err != nil {
		// Return empty object and error.
		return route, err
	}

	// Return query result.
	return route, nil
}

// CreateRoute method for creating route by given Route object.
func (q *RouteQueries) CreateRoute(r *models.Route) error {
	// Define query string.
	query := `INSERT INTO t_route VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	// Send query to database.
	_, err := q.Exec(query, r.ID, r.IsRetryable, r.RetryPeriod, r.MaxRetry, r.Slug, r.Token, r.CreatedAt, r.UpdatedAt, r.HostAddr, r.Name, r.Description, r.CommType, r.UserID, r.Query)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// UpdateRoute method for updating route by given Route object.
func (q *RouteQueries) UpdateRoute(id uuid.UUID, r *models.Route) error {
	// Define query string.
	query := `UPDATE t_route SET is_retryable = $1, retry_period = $2, max_retry = $3, updated_at = $4, host_addr = $5, name = $6, description = $7, comm_type = $8, query = $9 WHERE id = $10`

	// Send query to database.
	_, err := q.Exec(query, r.IsRetryable, r.RetryPeriod, r.MaxRetry, r.UpdatedAt, r.HostAddr, r.Name, r.Description, r.CommType, r.Query, id)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// DeleteRoute method for delete route by given ID.
func (q *RouteQueries) DeleteRoute(id uuid.UUID) error {
	// Define query string.
	query := `DELETE FROM t_route WHERE id = $1`

	// Send query to database.
	_, err := q.Exec(query, id)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// CountSlug method for preventing use of a same slug
func (q *RouteQueries) CountSlug(slug string) (int, error) {
	// Define query string.
	query := `SELECT count(slug) FROM t_route WHERE slug = $1`

	total := 0
	// Send query to database.
	err := q.Get(&total, query, slug)
	if err != nil {
		// Return empty object and error.
		return total, err
	}

	// Return query result.
	return total, nil
}

// UpdateToken method for generating new token to a route
func (q *RouteQueries) UpdateToken(id uuid.UUID, token string) error {
	// Define query string.
	query := `UPDATE t_route token = $1 where id = $2`

	// Send query to database.
	_, err := q.Exec(query, token, id)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// CreateLog method for creating route logging by given RouteLog object.
func (q *RouteQueries) CreateLog(r *models.RouteLog) error {
	// Define query string.
	query := `INSERT INTO t_log VALUES ($1, $2, $3, $4, $5, $6, $7)`

	// Send query to database.
	_, err := q.Exec(query, r.ID, r.Data, r.Type, r.CreatedAt, r.ResponseCode, r.RouteID, r.TrialAttempt)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}
