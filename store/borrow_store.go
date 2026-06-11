package store

import (
	"math"
	"time"

	"library-borrow-system/models"
)

func (s *Store) CreateBorrowRecord(record *models.BorrowRecord) {
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()
	s.BorrowRecords.Store(record.ID, record)
	appendToSliceMap(&s.borrowByBookID, record.BookID, record)
	appendToSliceMap(&s.borrowByReaderID, record.ReaderID, record)
}

func (s *Store) GetBorrowRecord(id string) (*models.BorrowRecord, bool) {
	v, ok := s.BorrowRecords.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*models.BorrowRecord), true
}

func (s *Store) GetActiveBorrowByBookID(bookID string) (*models.BorrowRecord, bool) {
	v, ok := s.borrowByBookID.Load(bookID)
	if !ok {
		return nil, false
	}
	records := v.([]interface{})
	for _, r := range records {
		record := r.(*models.BorrowRecord)
		if record.ReturnDate == nil {
			return record, true
		}
	}
	return nil, false
}

func (s *Store) UpdateBorrowRecord(record *models.BorrowRecord) bool {
	_, ok := s.BorrowRecords.Load(record.ID)
	if !ok {
		return false
	}
	record.UpdatedAt = time.Now()
	s.BorrowRecords.Store(record.ID, record)
	return true
}

func (s *Store) GetBorrowRecordsByBookID(bookID string) []*models.BorrowRecord {
	v, ok := s.borrowByBookID.Load(bookID)
	if !ok {
		return nil
	}
	raw := v.([]interface{})
	records := make([]*models.BorrowRecord, 0, len(raw))
	for _, r := range raw {
		records = append(records, r.(*models.BorrowRecord))
	}
	return records
}

func (s *Store) CalculateOverdue(record *models.BorrowRecord) (int, float64) {
	if record.ReturnDate != nil {
		return 0, 0
	}
	now := time.Now()
	if now.Before(record.DueDate) {
		return 0, 0
	}
	days := int(math.Ceil(now.Sub(record.DueDate).Hours() / 24))
	if days <= 0 {
		return 0, 0
	}
	fine := float64(days) * 0.5
	if fine > 30 {
		fine = 30
	}
	return days, fine
}
