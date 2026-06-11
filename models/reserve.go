package models

import "time"

type ReserveStatus string

const (
	ReserveStatusWaiting   ReserveStatus = "waiting"
	ReserveStatusAvailable ReserveStatus = "available"
	ReserveStatusCompleted ReserveStatus = "completed"
	ReserveStatusExpired   ReserveStatus = "expired"
	ReserveStatusCancelled ReserveStatus = "cancelled"
)

type ReserveRecord struct {
	ID              string        `json:"id"`
	ReaderID        string        `json:"reader_id"`
	BookID          string        `json:"book_id"`
	ReserveDate     time.Time     `json:"reserve_date"`
	Status          ReserveStatus `json:"status"`
	QueuePosition   int           `json:"queue_position"`
	AvailableDate   *time.Time    `json:"available_date"`
	ExpireDate      *time.Time    `json:"expire_date"`
	IsNotified      bool          `json:"is_notified"`
	ConvertedBorrowID *string     `json:"converted_borrow_id"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

type ReserveRequest struct {
	ReaderCardNo string `json:"reader_card_no" binding:"required"`
	BookRFID     string `json:"book_rfid" binding:"required"`
}

type ReserveResponse struct {
	Reserve       *ReserveRecord `json:"reserve"`
	Reader        *Reader        `json:"reader"`
	Book          *Book          `json:"book"`
	QueuePosition int            `json:"queue_position"`
	WaitCount     int            `json:"wait_count"`
}

type ReserveQueueInfo struct {
	BookID  string           `json:"book_id"`
	Book    *Book            `json:"book"`
	Queue   []*ReserveRecord `json:"queue"`
	Count   int              `json:"count"`
}
