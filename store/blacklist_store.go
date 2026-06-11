package store

import (
	"time"

	"library-borrow-system/models"
)

func (s *Store) AddToBlacklist(readerID string, reason string, fineAmount float64) *models.BlacklistRecord {
	record := &models.BlacklistRecord{
		ID:         "bl_" + readerID,
		ReaderID:   readerID,
		Reason:     reason,
		FineAmount: fineAmount,
		AddedAt:    time.Now(),
		IsActive:   true,
	}
	s.Blacklist.Store(record.ID, record)

	reader, ok := s.GetReader(readerID)
	if ok {
		reader.IsBlacklisted = true
		s.UpdateReader(reader)
	}
	return record
}

func (s *Store) RemoveFromBlacklist(readerID string) bool {
	record, ok := s.GetBlacklistRecord(readerID)
	if !ok {
		return false
	}
	now := time.Now()
	record.RemovedAt = &now
	record.IsActive = false
	s.Blacklist.Store(record.ID, record)

	reader, ok := s.GetReader(readerID)
	if ok {
		reader.IsBlacklisted = false
		s.UpdateReader(reader)
	}
	return true
}

func (s *Store) GetBlacklistRecord(readerID string) (*models.BlacklistRecord, bool) {
	v, ok := s.Blacklist.Load("bl_" + readerID)
	if !ok {
		return nil, false
	}
	return v.(*models.BlacklistRecord), true
}

func (s *Store) IsBlacklisted(readerID string) bool {
	record, ok := s.GetBlacklistRecord(readerID)
	if !ok {
		return false
	}
	return record.IsActive
}
