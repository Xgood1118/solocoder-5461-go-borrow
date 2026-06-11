package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"library-borrow-system/models"
	"library-borrow-system/store"
)

type RenewHandler struct {
	store *store.Store
}

func NewRenewHandler(s *store.Store) *RenewHandler {
	return &RenewHandler{store: s}
}

func (h *RenewHandler) RenewBook(c *gin.Context) {
	var req models.RenewRequest
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "读者在黑名单中，无法续借"})
		return
	}

	unpaidFine := h.store.GetReaderUnpaidFine(reader.ID)
	if unpaidFine > 30 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "有未结清罚款超过30元，无法续借"})
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

	if record.ReaderID != reader.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该书不是当前读者借阅"})
		return
	}

	overdueDays, _ := h.store.CalculateOverdue(record)
	if overdueDays > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该书已逾期，无法续借"})
		return
	}

	privilege := models.ReaderPrivileges[reader.Type]
	if record.RenewCount >= privilege.MaxRenewCount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已达续借次数上限"})
		return
	}

	now := time.Now()
	daysUntilDue := int(record.DueDate.Sub(now).Hours() / 24)
	if daysUntilDue > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "距到期日超过3天，暂不可续借"})
		return
	}

	reserveQueue := h.store.GetActiveReserveQueue(book.ID)
	if len(reserveQueue) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该书有预约队列，无法续借"})
		return
	}

	newDueDate := record.DueDate.AddDate(0, 0, 15)
	record.RenewCount++
	record.DueDate = newDueDate
	h.store.UpdateBorrowRecord(record)

	c.JSON(http.StatusOK, models.RenewResponse{
		Record:     record,
		Reader:     reader,
		Book:       book,
		NewDueDate: newDueDate,
		RenewCount: record.RenewCount,
	})
}

func (h *RenewHandler) CanRenew(c *gin.Context) {
	readerCardNo := c.Query("reader_card_no")
	bookRFID := c.Query("book_rfid")

	reader, ok := h.store.GetReaderByCardNo(readerCardNo)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "读者不存在"})
		return
	}

	book, ok := h.store.GetBookByRFID(bookRFID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图书不存在"})
		return
	}

	record, ok := h.store.GetActiveBorrowByBookID(book.ID)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"can_renew": false, "reason": "该书未被借出"})
		return
	}

	if record.ReaderID != reader.ID {
		c.JSON(http.StatusOK, gin.H{"can_renew": false, "reason": "该书不是当前读者借阅"})
		return
	}

	overdueDays, _ := h.store.CalculateOverdue(record)
	if overdueDays > 0 {
		c.JSON(http.StatusOK, gin.H{"can_renew": false, "reason": "该书已逾期"})
		return
	}

	privilege := models.ReaderPrivileges[reader.Type]
	if record.RenewCount >= privilege.MaxRenewCount {
		c.JSON(http.StatusOK, gin.H{"can_renew": false, "reason": "已达续借次数上限"})
		return
	}

	now := time.Now()
	daysUntilDue := int(record.DueDate.Sub(now).Hours() / 24)
	if daysUntilDue > 3 {
		c.JSON(http.StatusOK, gin.H{"can_renew": false, "reason": "距到期日超过3天"})
		return
	}

	reserveQueue := h.store.GetActiveReserveQueue(book.ID)
	if len(reserveQueue) > 0 {
		c.JSON(http.StatusOK, gin.H{"can_renew": false, "reason": "该书有预约队列"})
		return
	}

	if reader.IsBlacklisted {
		c.JSON(http.StatusOK, gin.H{"can_renew": false, "reason": "读者在黑名单中"})
		return
	}

	unpaidFine := h.store.GetReaderUnpaidFine(reader.ID)
	if unpaidFine > 30 {
		c.JSON(http.StatusOK, gin.H{"can_renew": false, "reason": "未结清罚款超过30元"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"can_renew":    true,
		"renew_count":  record.RenewCount,
		"max_renew":    privilege.MaxRenewCount,
		"current_due":  record.DueDate,
		"new_due":      record.DueDate.AddDate(0, 0, 15),
	})
}
