package models

import "time"

type RenewRequest struct {
	ReaderCardNo string `json:"reader_card_no" binding:"required"`
	BookRFID     string `json:"book_rfid" binding:"required"`
}

type RenewResponse struct {
	Record      *BorrowRecord `json:"record"`
	Reader      *Reader       `json:"reader"`
	Book        *Book         `json:"book"`
	NewDueDate  time.Time     `json:"new_due_date"`
	RenewCount  int           `json:"renew_count"`
}
