package model

import "time"

type IdempotencyKey struct {
	ID             int64     `json:"id"`
	IdempotencyKey string    `json:"idempotency_key"`
	Request        string    `json:"request"`
	Response       string    `json:"response"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}
