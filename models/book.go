package models

import "time"

type BookStatus string

const (
	BookStatusAvailable BookStatus = "available"
	BookStatusBorrowed  BookStatus = "borrowed"
	BookStatusReserved  BookStatus = "reserved"
	BookStatusLost      BookStatus = "lost"
)

type Book struct {
	ID        string     `json:"id"`
	ISBN      string     `json:"isbn"`
	Title     string     `json:"title"`
	Author    string     `json:"author"`
	Publisher string     `json:"publisher"`
	RFID      string     `json:"rfid"`
	Price     float64    `json:"price"`
	Status    BookStatus `json:"status"`
	Location  string     `json:"location"`
	Category  string     `json:"category"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type BookCreateRequest struct {
	ISBN      string  `json:"isbn"`
	Title     string  `json:"title" binding:"required"`
	Author    string  `json:"author"`
	Publisher string  `json:"publisher"`
	RFID      string  `json:"rfid" binding:"required"`
	Price     float64 `json:"price"`
	Location  string  `json:"location"`
	Category  string  `json:"category"`
}

type BookUpdateRequest struct {
	ISBN      string  `json:"isbn"`
	Title     string  `json:"title"`
	Author    string  `json:"author"`
	Publisher string  `json:"publisher"`
	Price     float64 `json:"price"`
	Location  string  `json:"location"`
	Category  string  `json:"category"`
}
