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
	IsActive    string    `db:"is_active"`
}

// log section
type RouteLog struct {
	ID           uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	Data         string    `db:"data" json:"data"`
	Type         string    `db:"type" json:"type"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	ResponseCode string    `db:"response_code" json:"response_code"`
	RouteID      uuid.UUID `db:"route_id" json:"route_id"`
	TrialAttempt int       `db:"trial_attempt" json:"trial_attempt"`
	IsResolved   string    `db:"is_resolved" json:"is_resolved"`
	Others       string    `db:"others" json:"others"`
}

type RouteLogView struct {
	ID           uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	Data         string    `db:"data" json:"data"`
	Type         string    `db:"type" json:"type"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	ResponseCode string    `db:"response_code" json:"response_code"`
	RouteID      uuid.UUID `db:"route_id" json:"route_id"`
	HostAddr     string    `db:"host_addr" json:"host_addr"`
	Name         string    `db:"name" json:"name"`
	Slug         string    `db:"slug" json:"slug"`
	TrialAttempt int       `db:"trial_attempt" json:"trial_attempt"`
	Others       string    `db:"others" json:"others"`
	IsResolved   string    `db:"is_resolved" json:"is_resolved"`
}

type LogData struct {
	No        int    `json:"no"`
	ReqHeader string `json:"req_header"`
	ReqBody   string `json:"req_body"`
	ResCode   int    `json:"res_code"`
	ResBody   string `json:"res_body"`
	IPAddr    string `json:"ip_address"`
}

type RetryLog struct {
	Total int       `json:"total"`
	ID    string    `json:"id"`
	Data  []LogData `json:"data"`
}

type Chart struct {
	Label  string `json:"label" db:"date"`
	Sukses int    `json:"sukses" db:"sukses"`
	Gagal  int    `json:"gagal" db:"gagal"`
}
