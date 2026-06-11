package store

import (
	"sort"
	"time"

	"library-borrow-system/models"
)

func (s *Store) CreateReserveRecord(record *models.ReserveRecord) {
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()
	if record.Status == "" {
		record.Status = models.ReserveStatusWaiting
	}
	s.ReserveRecords.Store(record.ID, record)
	appendToSliceMap(&s.reserveByBookID, record.BookID, record)
	appendToSliceMap(&s.reserveByReaderID, record.ReaderID, record)
	s.updateReserveQueuePositions(record.BookID)
}

func (s *Store) GetReserveRecord(id string) (*models.ReserveRecord, bool) {
	v, ok := s.ReserveRecords.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*models.ReserveRecord), true
}

func (s *Store) UpdateReserveRecord(record *models.ReserveRecord) bool {
	_, ok := s.ReserveRecords.Load(record.ID)
	if !ok {
		return false
	}
	record.UpdatedAt = time.Now()
	s.ReserveRecords.Store(record.ID, record)
	return true
}

func (s *Store) GetReserveQueue(bookID string) []*models.ReserveRecord {
	v, ok := s.reserveByBookID.Load(bookID)
	if !ok {
		return nil
	}
	raw := v.([]interface{})
	records := make([]*models.ReserveRecord, 0, len(raw))
	for _, r := range raw {
		records = append(records, r.(*models.ReserveRecord))
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].ReserveDate.Before(records[j].ReserveDate)
	})
	return records
}

func (s *Store) GetActiveReserveQueue(bookID string) []*models.ReserveRecord {
	all := s.GetReserveQueue(bookID)
	active := make([]*models.ReserveRecord, 0)
	for _, r := range all {
		if r.Status == models.ReserveStatusWaiting || r.Status == models.ReserveStatusAvailable {
			active = append(active, r)
		}
	}
	return active
}

func (s *Store) updateReserveQueuePositions(bookID string) {
	queue := s.GetActiveReserveQueue(bookID)
	for i, r := range queue {
		r.QueuePosition = i + 1
		s.UpdateReserveRecord(r)
	}
}

func (s *Store) GetFirstWaitingReserve(bookID string) (*models.ReserveRecord, bool) {
	queue := s.GetActiveReserveQueue(bookID)
	for _, r := range queue {
		if r.Status == models.ReserveStatusWaiting {
			return r, true
		}
	}
	return nil, false
}

func (s *Store) GetAvailableReserve(bookID string) (*models.ReserveRecord, bool) {
	queue := s.GetActiveReserveQueue(bookID)
	for _, r := range queue {
		if r.Status == models.ReserveStatusAvailable {
			return r, true
		}
	}
	return nil, false
}

func (s *Store) GetReaderReserveRecords(readerID string) []*models.ReserveRecord {
	v, ok := s.reserveByReaderID.Load(readerID)
	if !ok {
		return nil
	}
	raw := v.([]interface{})
	records := make([]*models.ReserveRecord, 0, len(raw))
	for _, r := range raw {
		records = append(records, r.(*models.ReserveRecord))
	}
	return records
}

func (s *Store) HasActiveReserve(bookID string, readerID string) bool {
	records := s.GetReaderReserveRecords(readerID)
	for _, r := range records {
		if r.BookID == bookID && (r.Status == models.ReserveStatusWaiting || r.Status == models.ReserveStatusAvailable) {
			return true
		}
	}
	return false
}

func (s *Store) ExpireReserve(reserveID string) {
	reserve, ok := s.GetReserveRecord(reserveID)
	if !ok {
		return
	}
	reserve.Status = models.ReserveStatusExpired
	s.UpdateReserveRecord(reserve)
	s.updateReserveQueuePositions(reserve.BookID)
}

func (s *Store) NotifyNextReserve(bookID string) (*models.ReserveRecord, bool) {
	next, ok := s.GetFirstWaitingReserve(bookID)
	if !ok {
		return nil, false
	}
	now := time.Now()
	expireDate := now.AddDate(0, 0, 3)
	next.Status = models.ReserveStatusAvailable
	next.AvailableDate = &now
	next.ExpireDate = &expireDate
	next.IsNotified = true
	s.UpdateReserveRecord(next)
	return next, true
}
