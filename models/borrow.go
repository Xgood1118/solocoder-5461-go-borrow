package models

import "time"

type BorrowRecord struct {
	ID           string    `json:"id"`
	ReaderID     string    `json:"reader_id"`
	BookID       string    `json:"book_id"`
	BorrowDate   time.Time `json:"borrow_date"`
	DueDate      time.Time `json:"due_date"`
	ReturnDate   *time.Time `json:"return_date"`
	RenewCount   int       `json:"renew_count"`
	IsOverdue    bool      `json:"is_overdue"`
	OverdueDays  int       `json:"overdue_days"`
	FineAmount   float64   `json:"fine_amount"`
	IsFinePaid   bool      `json:"is_fine_paid"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type BorrowRequest struct {
	ReaderCardNo string `json:"reader_card_no" binding:"required"`
	BookRFID     string `json:"book_rfid" binding:"required"`
}

type ReturnRequest struct {
	BookRFID string `json:"book_rfid" binding:"required"`
}

type BorrowResponse struct {
	Record   *BorrowRecord `json:"record"`
	Reader   *Reader       `json:"reader"`
	Book     *Book         `json:"book"`
	DueDate  time.Time     `json:"due_date"`
}

type ReturnResponse struct {
	Record      *BorrowRecord `json:"record"`
	Reader      *Reader       `json:"reader"`
	Book        *Book         `json:"book"`
	OverdueDays int           `json:"overdue_days"`
	FineAmount  float64       `json:"fine_amount"`
}
