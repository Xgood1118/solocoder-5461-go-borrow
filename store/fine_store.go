package store

import (
	"time"

	"library-borrow-system/models"
)

func (s *Store) CreateFineRecord(record *models.FineRecord) {
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()
	if record.Status == "" {
		record.Status = models.FineStatusUnpaid
	}
	s.FineRecords.Store(record.ID, record)
	appendToSliceMap(&s.fineByReaderID, record.ReaderID, record)
}

func (s *Store) GetFineRecord(id string) (*models.FineRecord, bool) {
	v, ok := s.FineRecords.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*models.FineRecord), true
}

func (s *Store) UpdateFineRecord(record *models.FineRecord) bool {
	_, ok := s.FineRecords.Load(record.ID)
	if !ok {
		return false
	}
	record.UpdatedAt = time.Now()
	s.FineRecords.Store(record.ID, record)
	return true
}

func (s *Store) GetReaderFineRecords(readerID string) []*models.FineRecord {
	v, ok := s.fineByReaderID.Load(readerID)
	if !ok {
		return nil
	}
	raw := v.([]interface{})
	records := make([]*models.FineRecord, 0, len(raw))
	for _, r := range raw {
		records = append(records, r.(*models.FineRecord))
	}
	return records
}

func (s *Store) GetReaderUnpaidFines(readerID string) []*models.FineRecord {
	all := s.GetReaderFineRecords(readerID)
	unpaid := make([]*models.FineRecord, 0)
	for _, r := range all {
		if r.Status == models.FineStatusUnpaid {
			unpaid = append(unpaid, r)
		}
	}
	return unpaid
}

func (s *Store) PayFine(fineID string) bool {
	fine, ok := s.GetFineRecord(fineID)
	if !ok {
		return false
	}
	now := time.Now()
	fine.Status = models.FineStatusPaid
	fine.PaidDate = &now
	return s.UpdateFineRecord(fine)
}

func (s *Store) PayReaderFines(readerID string, amount float64) (float64, []*models.FineRecord) {
	unpaid := s.GetReaderUnpaidFines(readerID)
	if len(unpaid) == 0 || amount <= 0 {
		return 0, nil
	}

	paid := make([]*models.FineRecord, 0)
	remaining := amount

	for _, fine := range unpaid {
		if remaining <= 0 {
			break
		}
		if remaining >= fine.Amount {
			now := time.Now()
			fine.Status = models.FineStatusPaid
			fine.PaidDate = &now
			s.UpdateFineRecord(fine)
			remaining -= fine.Amount
			paid = append(paid, fine)
		} else {
			break
		}
	}

	paidAmount := amount - remaining
	return paidAmount, paid
}
