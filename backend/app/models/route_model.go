package models

import (
	"time"

	"github.com/google/uuid"
)

// Route struct to describe Route object.
type Route struct {
	ID          uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	IsRetryable string    `db:"is_retryable" json:"is_retryable" validate:"required,lte=1"`
	RetryPeriod int       `db:"retry_period" json:"retry_period" `
	MaxRetry    int       `db:"max_retry" json:"max_retry"`
	Slug        string    `db:"slug" json:"slug"`
	Token       string    `db:"token" json:"token"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	HostAddr    string    `db:"host_addr" json:"host_addr"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CommType    string    `db:"comm_type" json:"comm_type"`
	UserID      uuid.UUID `db:"user_id" json:"user_id"`
	Query       string    `db:"query" json:"query"`
}

type RouteLog struct {
	ID           uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	Data         string    `db:"data" json:"data"`
	Type         string    `db:"type" json:"type"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	ResponseCode string    `db:"response_code" json:"response_code"`
	RouteID      uuid.UUID `db:"route_id" json:"route_id"`
}
