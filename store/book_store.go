package store

import (
	"time"

	"library-borrow-system/models"
)

func (s *Store) CreateBook(book *models.Book) {
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()
	if book.Status == "" {
		book.Status = models.BookStatusAvailable
	}
	s.Books.Store(book.ID, book)
	s.bookByRFID.Store(book.RFID, book)
}

func (s *Store) GetBook(id string) (*models.Book, bool) {
	v, ok := s.Books.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*models.Book), true
}

func (s *Store) GetBookByRFID(rfid string) (*models.Book, bool) {
	v, ok := s.bookByRFID.Load(rfid)
	if !ok {
		return nil, false
	}
	return v.(*models.Book), true
}

func (s *Store) UpdateBook(book *models.Book) bool {
	_, ok := s.Books.Load(book.ID)
	if !ok {
		return false
	}
	oldBook, _ := s.GetBook(book.ID)
	if oldBook.RFID != book.RFID {
		s.bookByRFID.Delete(oldBook.RFID)
		s.bookByRFID.Store(book.RFID, book)
	}
	book.UpdatedAt = time.Now()
	s.Books.Store(book.ID, book)
	return true
}

func (s *Store) DeleteBook(id string) bool {
	book, ok := s.GetBook(id)
	if !ok {
		return false
	}
	s.Books.Delete(id)
	s.bookByRFID.Delete(book.RFID)
	return true
}

func (s *Store) UpdateBookStatus(bookID string, status models.BookStatus) bool {
	book, ok := s.GetBook(bookID)
	if !ok {
		return false
	}
	book.Status = status
	book.UpdatedAt = time.Now()
	s.Books.Store(bookID, book)
	s.bookByRFID.Store(book.RFID, book)
	return true
}
