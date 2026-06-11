package store

import (
	"sync"

	"library-borrow-system/models"
)

type Store struct {
	Readers       sync.Map
	Books         sync.Map
	BorrowRecords sync.Map
	ReserveRecords sync.Map
	FineRecords   sync.Map
	Blacklist     sync.Map

	readerByCardNo sync.Map
	bookByRFID     sync.Map
	borrowByBookID sync.Map
	borrowByReaderID sync.Map
	reserveByBookID  sync.Map
	reserveByReaderID sync.Map
	fineByReaderID   sync.Map
}

var GlobalStore = &Store{}

func (s *Store) ClearAll() {
	s.Readers = sync.Map{}
	s.Books = sync.Map{}
	s.BorrowRecords = sync.Map{}
	s.ReserveRecords = sync.Map{}
	s.FineRecords = sync.Map{}
	s.Blacklist = sync.Map{}
	s.readerByCardNo = sync.Map{}
	s.bookByRFID = sync.Map{}
	s.borrowByBookID = sync.Map{}
	s.borrowByReaderID = sync.Map{}
	s.reserveByBookID = sync.Map{}
	s.reserveByReaderID = sync.Map{}
	s.fineByReaderID = sync.Map{}
}

func appendToSliceMap(m *sync.Map, key string, value interface{}) {
	var slice []interface{}
	if v, ok := m.Load(key); ok {
		slice = v.([]interface{})
	}
	slice = append(slice, value)
	m.Store(key, slice)
}

func removeFromSliceMap(m *sync.Map, key string, id string, matchFunc func(interface{}) string) {
	if v, ok := m.Load(key); ok {
		slice := v.([]interface{})
		newSlice := make([]interface{}, 0, len(slice))
		for _, item := range slice {
			if matchFunc(item) != id {
				newSlice = append(newSlice, item)
			}
		}
		if len(newSlice) > 0 {
			m.Store(key, newSlice)
		} else {
			m.Delete(key)
		}
	}
}

func getSliceFromMap(m *sync.Map, key string) []interface{} {
	if v, ok := m.Load(key); ok {
		return v.([]interface{})
	}
	return nil
}

func deleteKeyFromSliceMap(m *sync.Map, key string) {
	m.Delete(key)
}

func (s *Store) GetAllReaders() []*models.Reader {
	var readers []*models.Reader
	s.Readers.Range(func(key, value interface{}) bool {
		readers = append(readers, value.(*models.Reader))
		return true
	})
	return readers
}

func (s *Store) GetAllBooks() []*models.Book {
	var books []*models.Book
	s.Books.Range(func(key, value interface{}) bool {
		books = append(books, value.(*models.Book))
		return true
	})
	return books
}

func (s *Store) GetAllBorrowRecords() []*models.BorrowRecord {
	var records []*models.BorrowRecord
	s.BorrowRecords.Range(func(key, value interface{}) bool {
		records = append(records, value.(*models.BorrowRecord))
		return true
	})
	return records
}

func (s *Store) GetAllReserveRecords() []*models.ReserveRecord {
	var records []*models.ReserveRecord
	s.ReserveRecords.Range(func(key, value interface{}) bool {
		records = append(records, value.(*models.ReserveRecord))
		return true
	})
	return records
}

func (s *Store) GetAllFineRecords() []*models.FineRecord {
	var records []*models.FineRecord
	s.FineRecords.Range(func(key, value interface{}) bool {
		records = append(records, value.(*models.FineRecord))
		return true
	})
	return records
}

func (s *Store) GetAllBlacklistRecords() []*models.BlacklistRecord {
	var records []*models.BlacklistRecord
	s.Blacklist.Range(func(key, value interface{}) bool {
		records = append(records, value.(*models.BlacklistRecord))
		return true
	})
	return records
}
