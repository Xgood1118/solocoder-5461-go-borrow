package store

import (
	"errors"
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

type DeleteReaderResult struct {
	Success          bool
	Message          string
	UnreturnedBooks  []*models.Book
	UnpaidFineAmount float64
	CancelledReserves int
}

func (s *Store) DeleteReaderWithCascade(id string) (*DeleteReaderResult, error) {
	reader, ok := s.GetReader(id)
	if !ok {
		return nil, errors.New("读者不存在")
	}

	result := &DeleteReaderResult{}

	activeBorrows := s.GetReaderActiveBorrowRecords(id)
	if len(activeBorrows) > 0 {
		unreturnedBooks := make([]*models.Book, 0, len(activeBorrows))
		for _, br := range activeBorrows {
			if book, bookOk := s.GetBook(br.BookID); bookOk {
				unreturnedBooks = append(unreturnedBooks, book)
			}
		}
		result.UnreturnedBooks = unreturnedBooks
		result.Message = "存在未归还的图书，无法删除"
		result.Success = false
		return result, errors.New(result.Message)
	}

	unpaidFine := s.GetReaderUnpaidFine(id)
	if unpaidFine > 0 {
		result.UnpaidFineAmount = unpaidFine
		result.Message = "存在未结清的罚款，无法删除"
		result.Success = false
		return result, errors.New(result.Message)
	}

	activeReserves := s.GetReaderReserveRecords(id)
	cancelledCount := 0
	for _, r := range activeReserves {
		if r.Status == models.ReserveStatusWaiting || r.Status == models.ReserveStatusAvailable {
			wasAvailable := r.Status == models.ReserveStatusAvailable
			r.Status = models.ReserveStatusCancelled
			s.UpdateReserveRecord(r)
			cancelledCount++

			bookID := r.BookID
			if wasAvailable {
				s.NotifyNextReserve(bookID)
			}
			s.updateReserveQueuePositions(bookID)
		}
	}
	result.CancelledReserves = cancelledCount

	if s.IsBlacklisted(id) {
		s.RemoveFromBlacklist(id)
	}

	borrowRecords := s.GetReaderBorrowRecords(id)
	for _, br := range borrowRecords {
		s.BorrowRecords.Delete(br.ID)
		removeFromSliceMap(&s.borrowByBookID, br.BookID, br.ID, func(v interface{}) string {
			return v.(*models.BorrowRecord).ID
		})
	}
	deleteKeyFromSliceMap(&s.borrowByReaderID, id)

	for _, r := range activeReserves {
		s.ReserveRecords.Delete(r.ID)
		removeFromSliceMap(&s.reserveByBookID, r.BookID, r.ID, func(v interface{}) string {
			return v.(*models.ReserveRecord).ID
		})
	}
	deleteKeyFromSliceMap(&s.reserveByReaderID, id)

	fineRecords := s.GetReaderFineRecords(id)
	for _, fr := range fineRecords {
		s.FineRecords.Delete(fr.ID)
	}
	deleteKeyFromSliceMap(&s.fineByReaderID, id)

	s.Blacklist.Delete("bl_" + id)

	s.Readers.Delete(id)
	s.readerByCardNo.Delete(reader.CardNo)

	result.Success = true
	result.Message = "删除成功"
	return result, nil
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
