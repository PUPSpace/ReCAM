package queries

import (
	// "github.com/google/uuid"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kaleemubarok/recam/backend/app/models"
)

// LogQueries struct for queries from Route model.
type LogQueries struct {
	*sqlx.DB
}

// GetLogs method for getting all routes.
func (q *LogQueries) GetLogs() ([]models.RouteLogView, error) {
	// Define routes variable.
	logs := []models.RouteLogView{}

	// Define query string.
	query := `SELECT l.*, r.slug, r.name, r.host_addr FROM t_log l
	INNER JOIN t_route r ON l.route_id=r.id
	ORDER BY l.created_at DESC
	LIMIT 50`

	// Send query to database.
	err := q.Select(&logs, query)
	if err != nil {
		// Return empty object and error.
		return logs, err
	}

	// Return query result.
	return logs, nil
}

/*
// GetLogsByAuthor method for getting all routes by given author.
func (q *LogQueries) GetLogsByAuthor(author string) ([]models.Route, error) {
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
func (q *LogQueries) GetLog(id uuid.UUID) (models.RouteLogView, error) {
	// Define route variable.
	rlog := models.RouteLogView{}

	// Define query string.
	query := `SELECT l.*, r.slug, r.name, r.host_addr FROM t_log l
	INNER JOIN t_route r ON l.route_id=r.id
	WHERE l.id = $1`

	// Send query to database.
	err := q.Get(&rlog, query, id)
	if err != nil {
		// Return empty object and error.
		return rlog, err
	}

	// Return query result.
	return rlog, nil
}

// // GetRoute method for getting one route by given slug.
// func (q *LogQueries) GetLogslug(slug string) (models.Route, error) {
// 	// Define route variable.
// 	route := models.Route{}

// 	// Define query string.
// 	query := `SELECT * FROM t_route WHERE slug = $1`

// 	// Send query to database.
// 	err := q.Get(&route, query, slug)
// 	if err != nil {
// 		// Return empty object and error.
// 		return route, err
// 	}

// 	// Return query result.
// 	return route, nil
// }

// // CreateRoute method for creating route by given Route object.
// func (q *LogQueries) CreateRoute(r *models.Route) error {
// 	// Define query string.
// 	query := `INSERT INTO t_route VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

// 	// Send query to database.
// 	_, err := q.Exec(query, r.ID, r.IsRetryable, r.RetryPeriod, r.MaxRetry, r.Slug, r.Token, r.CreatedAt, r.UpdatedAt, r.HostAddr, r.Name, r.Description, r.CommType, r.UserID, r.Query)
// 	if err != nil {
// 		// Return only error.
// 		return err
// 	}

// 	// This query returns nothing.
// 	return nil
// }

// // UpdateRoute method for updating route by given Route object.
// func (q *LogQueries) UpdateRoute(id uuid.UUID, r *models.Route) error {
// 	// Define query string.
// 	query := `UPDATE t_route SET is_retryable = $1, retry_period = $2, max_retry = $3, updated_at = $4, host_addr = $5, name = $6, description = $7, comm_type = $8, query = $9 WHERE id = $10`

// 	// Send query to database.
// 	_, err := q.Exec(query, r.IsRetryable, r.RetryPeriod, r.MaxRetry, r.UpdatedAt, r.HostAddr, r.Name, r.Description, r.CommType, r.Query, id)
// 	if err != nil {
// 		// Return only error.
// 		return err
// 	}

// 	// This query returns nothing.
// 	return nil
// }

// // DeleteRoute method for delete route by given ID.
// func (q *LogQueries) DeleteRoute(id uuid.UUID) error {
// 	// Define query string.
// 	query := `DELETE FROM t_route WHERE id = $1`

// 	// Send query to database.
// 	_, err := q.Exec(query, id)
// 	if err != nil {
// 		// Return only error.
// 		return err
// 	}

// 	// This query returns nothing.
// 	return nil
// }

// // CountSlug method for preventing use of a same slug
// func (q *LogQueries) CountSlug(slug string) (int, error) {
// 	// Define query string.
// 	query := `SELECT count(slug) FROM t_route WHERE slug = $1`

// 	total := 0
// 	// Send query to database.
// 	err := q.Get(&total, query, slug)
// 	if err != nil {
// 		// Return empty object and error.
// 		return total, err
// 	}

// 	// Return query result.
// 	return total, nil
// }

// // UpdateToken method for generating new token to a route
// func (q *LogQueries) UpdateToken(id uuid.UUID, token string) error {
// 	// Define query string.
// 	query := `UPDATE t_route token = $1 where id = $2`

// 	// Send query to database.
// 	_, err := q.Exec(query, token, id)
// 	if err != nil {
// 		// Return only error.
// 		return err
// 	}

// 	// This query returns nothing.
// 	return nil
// }

// // CreateLog method for creating route logging by given RouteLog object.
// func (q *LogQueries) CreateLog(r *models.RouteLog) error {
// 	// Define query string.
// 	query := `INSERT INTO t_log VALUES ($1, $2, $3, $4, $5, $6)`

// 	// Send query to database.
// 	_, err := q.Exec(query, r.ID, r.Data, r.Type, r.CreatedAt, r.ResponseCode, r.RouteID)
// 	if err != nil {
// 		// Return only error.
// 		return err
// 	}

// 	// This query returns nothing.
// 	return nil
// }

// GetLogs method for getting all routes.
func (q *LogQueries) GetLogsSpecial(ec string) ([]models.RouteLogView, error) {
	// Define routes variable.
	logs := []models.RouteLogView{}

	// Define query string.
	query := ""
	if ec == "4XX" { //error;
		query = `SELECT l.*, r.slug, r.name, r.host_addr FROM t_log l
		INNER JOIN t_route r ON l.route_id=r.id
		WHERE l.response_code >=399 AND response_code < 500
		AND DATE(l.created_at) = CURRENT_DATE
		ORDER BY l.created_at DESC
		LIMIT 50;`
	} else { //gagal
		query = `SELECT l.*, r.slug, r.name, r.host_addr FROM t_log l
		INNER JOIN t_route r ON l.route_id=r.id
		WHERE l.response_code >=500
		ORDER BY l.created_at DESC
		LIMIT 50`
	}

	// Send query to database.
	err := q.Select(&logs, query)
	if err != nil {
		// Return empty object and error.
		return logs, err
	}

	// Return query result.
	return logs, nil
}

// GetLogs method for getting all routes.
func (q *LogQueries) CountUnresolved5XX() (int, error) {
	// Define res variable.
	res := 0

	// Define query string.
	query := `SELECT count(id) FROM t_log
	WHERE response_code >= 500
	AND is_resolved = 'N'
	GROUP BY is_resolved
	LIMIT 50;`

	// Send query to database.
	err := q.Select(&res, query)
	if err != nil {
		// Return empty object and error.
		return res, err
	}

	// Return query result.
	return res, nil
}
