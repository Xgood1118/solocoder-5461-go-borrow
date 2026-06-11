package models

import "time"

type FineStatus string

const (
	FineStatusUnpaid FineStatus = "unpaid"
	FineStatusPaid   FineStatus = "paid"
)

type FineRecord struct {
	ID          string     `json:"id"`
	ReaderID    string     `json:"reader_id"`
	BorrowID    string     `json:"borrow_id"`
	BookID      string     `json:"book_id"`
	Amount      float64    `json:"amount"`
	Reason      string     `json:"reason"`
	Status      FineStatus `json:"status"`
	OverdueDays int        `json:"overdue_days"`
	PaidDate    *time.Time `json:"paid_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type FinePayRequest struct {
	ReaderCardNo string  `json:"reader_card_no" binding:"required"`
	Amount       float64 `json:"amount" binding:"required,gt=0"`
}

type FinePayResponse struct {
	PaidAmount    float64       `json:"paid_amount"`
	RemainingFine float64       `json:"remaining_fine"`
	Records       []*FineRecord `json:"paid_records"`
}
