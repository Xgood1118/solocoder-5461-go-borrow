package models

import "time"

type BlacklistRecord struct {
	ID         string    `json:"id"`
	ReaderID   string    `json:"reader_id"`
	Reason     string    `json:"reason"`
	FineAmount float64   `json:"fine_amount"`
	AddedAt    time.Time `json:"added_at"`
	RemovedAt  *time.Time `json:"removed_at"`
	IsActive   bool      `json:"is_active"`
}
