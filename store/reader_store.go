package store

import (
	"time"

	"library-borrow-system/models"
)

func (s *Store) CreateReader(reader *models.Reader) {
	reader.CreatedAt = time.Now()
	reader.UpdatedAt = time.Now()
	s.Readers.Store(reader.ID, reader)
	s.readerByCardNo.Store(reader.CardNo, reader)
}

func (s *Store) GetReader(id string) (*models.Reader, bool) {
	v, ok := s.Readers.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*models.Reader), true
}

func (s *Store) GetReaderByCardNo(cardNo string) (*models.Reader, bool) {
	v, ok := s.readerByCardNo.Load(cardNo)
	if !ok {
		return nil, false
	}
	return v.(*models.Reader), true
}

func (s *Store) UpdateReader(reader *models.Reader) bool {
	_, ok := s.Readers.Load(reader.ID)
	if !ok {
		return false
	}
	oldReader, _ := s.GetReader(reader.ID)
	if oldReader.CardNo != reader.CardNo {
		s.readerByCardNo.Delete(oldReader.CardNo)
		s.readerByCardNo.Store(reader.CardNo, reader)
	}
	reader.UpdatedAt = time.Now()
	s.Readers.Store(reader.ID, reader)
	return true
}

func (s *Store) DeleteReader(id string) bool {
	reader, ok := s.GetReader(id)
	if !ok {
		return false
	}
	s.Readers.Delete(id)
	s.readerByCardNo.Delete(reader.CardNo)
	return true
}

func (s *Store) GetReaderBorrowCount(readerID string) int {
	v, ok := s.borrowByReaderID.Load(readerID)
	if !ok {
		return 0
	}
	records := v.([]interface{})
	count := 0
	for _, r := range records {
		record := r.(*models.BorrowRecord)
		if record.ReturnDate == nil {
			count++
		}
	}
	return count
}

func (s *Store) GetReaderUnpaidFine(readerID string) float64 {
	v, ok := s.fineByReaderID.Load(readerID)
	if !ok {
		return 0
	}
	records := v.([]interface{})
	var total float64
	for _, r := range records {
		record := r.(*models.FineRecord)
		if record.Status == models.FineStatusUnpaid {
			total += record.Amount
		}
	}
	return total
}

func (s *Store) GetReaderBorrowRecords(readerID string) []*models.BorrowRecord {
	v, ok := s.borrowByReaderID.Load(readerID)
	if !ok {
		return nil
	}
	rawRecords := v.([]interface{})
	records := make([]*models.BorrowRecord, 0, len(rawRecords))
	for _, r := range rawRecords {
		records = append(records, r.(*models.BorrowRecord))
	}
	return records
}

func (s *Store) GetReaderActiveBorrowRecords(readerID string) []*models.BorrowRecord {
	allRecords := s.GetReaderBorrowRecords(readerID)
	active := make([]*models.BorrowRecord, 0)
	for _, r := range allRecords {
		if r.ReturnDate == nil {
			active = append(active, r)
		}
	}
	return active
}
