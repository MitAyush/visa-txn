package model

import "time"

type Transaction struct {
	ID              int64     `json:"id"`
	IdempotencyKey  string    `json:"idempotency_key"`
	AccountID       int64     `json:"account_id"`
	OperationTypeID int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventDate       time.Time `json:"event_date"`
	CreatedAt       time.Time `json:"created_at"`
}
