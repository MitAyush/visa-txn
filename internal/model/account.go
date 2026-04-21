package model

import "time"

type Account struct {
	ID             int64     `json:"id"`
	AccountID      int64     `json:"account_id"`
	DocumentNumber string    `json:"document_number"`
	CreatedAt      time.Time `json:"created_at"`
}
