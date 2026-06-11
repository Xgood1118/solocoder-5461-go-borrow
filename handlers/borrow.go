package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"library-borrow-system/models"
	"library-borrow-system/store"
	"library-borrow-system/utils"
)

type BorrowHandler struct {
	store *store.Store
}

func NewBorrowHandler(s *store.Store) *BorrowHandler {
	return &BorrowHandler{store: s}
}

func (h *BorrowHandler) BorrowBook(c *gin.Context) {
	var req models.BorrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reader, ok := h.store.GetReaderByCardNo(req.ReaderCardNo)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "读者不存在"})
		return
	}

	if reader.IsBlacklisted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "读者在黑名单中，无法借阅"})
		return
	}

	unpaidFine := h.store.GetReaderUnpaidFine(reader.ID)
	if unpaidFine > 30 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "读者有超过30元未结清罚款，无法借阅"})
		return
	}

	book, ok := h.store.GetBookByRFID(req.BookRFID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图书不存在"})
		return
	}

	if book.Status != models.BookStatusAvailable {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图书当前不可借"})
		return
	}

	privilege := models.ReaderPrivileges[reader.Type]
	currentBorrowCount := h.store.GetReaderBorrowCount(reader.ID)
	if currentBorrowCount >= privilege.MaxBorrowCount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已达借阅上限"})
		return
	}

	borrowDate := time.Now()
	dueDate := borrowDate.AddDate(0, 0, privilege.MaxBorrowDays)

	record := &models.BorrowRecord{
		ID:         utils.GenerateBorrowID(),
		ReaderID:   reader.ID,
		BookID:     book.ID,
		BorrowDate: borrowDate,
		DueDate:    dueDate,
		RenewCount: 0,
	}

	h.store.CreateBorrowRecord(record)
	h.store.UpdateBookStatus(book.ID, models.BookStatusBorrowed)

	reserveQueue := h.store.GetActiveReserveQueue(book.ID)
	if len(reserveQueue) > 0 {
		first := reserveQueue[0]
		if first.Status == models.ReserveStatusAvailable && first.ReaderID == reader.ID {
			first.Status = models.ReserveStatusCompleted
			first.ConvertedBorrowID = &record.ID
			h.store.UpdateReserveRecord(first)
		}
	}

	c.JSON(http.StatusOK, models.BorrowResponse{
		Record:  record,
		Reader:  reader,
		Book:    book,
		DueDate: dueDate,
	})
}

func (h *BorrowHandler) ReturnBook(c *gin.Context) {
	var req models.ReturnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, ok := h.store.GetBookByRFID(req.BookRFID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图书不存在"})
		return
	}

	record, ok := h.store.GetActiveBorrowByBookID(book.ID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该书未被借出"})
		return
	}

	reader, _ := h.store.GetReader(record.ReaderID)

	overdueDays, fineAmount := h.store.CalculateOverdue(record)
	returnDate := time.Now()

	record.ReturnDate = &returnDate
	record.IsOverdue = overdueDays > 0
	record.OverdueDays = overdueDays
	record.FineAmount = fineAmount

	if fineAmount > 0 {
		fineRecord := &models.FineRecord{
			ID:          utils.GenerateFineID(),
			ReaderID:    record.ReaderID,
			BorrowID:    record.ID,
			BookID:      book.ID,
			Amount:      fineAmount,
			Reason:      "逾期罚款",
			OverdueDays: overdueDays,
		}
		h.store.CreateFineRecord(fineRecord)
		record.IsFinePaid = false

		totalUnpaidFine := h.store.GetReaderUnpaidFine(record.ReaderID)
		if totalUnpaidFine > 30 {
			h.store.AddToBlacklist(record.ReaderID, "逾期罚款超过30元", totalUnpaidFine)
		}
	} else {
		record.IsFinePaid = true
	}

	h.store.UpdateBorrowRecord(record)
	h.store.UpdateBookStatus(book.ID, models.BookStatusAvailable)

	reserveQueue := h.store.GetActiveReserveQueue(book.ID)
	if len(reserveQueue) > 0 {
		nextReserve, notified := h.store.NotifyNextReserve(book.ID)
		if notified && nextReserve != nil {
			_ = nextReserve
		}
	}

	c.JSON(http.StatusOK, models.ReturnResponse{
		Record:      record,
		Reader:      reader,
		Book:        book,
		OverdueDays: overdueDays,
		FineAmount:  fineAmount,
	})
}

func (h *BorrowHandler) GetBorrowRecord(c *gin.Context) {
	id := c.Param("id")
	record, ok := h.store.GetBorrowRecord(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "借阅记录不存在"})
		return
	}
	c.JSON(http.StatusOK, record)
}

func (h *BorrowHandler) ListBorrowRecords(c *gin.Context) {
	records := h.store.GetAllBorrowRecords()
	c.JSON(http.StatusOK, records)
}

func (h *BorrowHandler) GetReaderBorrowRecords(c *gin.Context) {
	readerID := c.Param("readerId")
	records := h.store.GetReaderBorrowRecords(readerID)
	c.JSON(http.StatusOK, records)
}

func (h *BorrowHandler) GetBookBorrowRecords(c *gin.Context) {
	bookID := c.Param("bookId")
	records := h.store.GetBorrowRecordsByBookID(bookID)
	c.JSON(http.StatusOK, records)
}
