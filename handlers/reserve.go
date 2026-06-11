package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"library-borrow-system/models"
	"library-borrow-system/store"
	"library-borrow-system/utils"
)

type ReserveHandler struct {
	store *store.Store
}

func NewReserveHandler(s *store.Store) *ReserveHandler {
	return &ReserveHandler{store: s}
}

func (h *ReserveHandler) CreateReserve(c *gin.Context) {
	var req models.ReserveRequest
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "读者在黑名单中，无法预约"})
		return
	}

	book, ok := h.store.GetBookByRFID(req.BookRFID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图书不存在"})
		return
	}

	if book.Status == models.BookStatusLost {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图书已丢失，无法预约"})
		return
	}

	if h.store.HasActiveReserve(book.ID, reader.ID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您已预约过该书"})
		return
	}

	if book.Status == models.BookStatusAvailable {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该书在馆可借，无需预约"})
		return
	}

	queue := h.store.GetActiveReserveQueue(book.ID)
	queuePosition := len(queue) + 1

	reserve := &models.ReserveRecord{
		ID:            utils.GenerateReserveID(),
		ReaderID:      reader.ID,
		BookID:        book.ID,
		ReserveDate:   time.Now(),
		Status:        models.ReserveStatusWaiting,
		QueuePosition: queuePosition,
	}

	h.store.CreateReserveRecord(reserve)

	c.JSON(http.StatusCreated, models.ReserveResponse{
		Reserve:       reserve,
		Reader:        reader,
		Book:          book,
		QueuePosition: queuePosition,
		WaitCount:     queuePosition - 1,
	})
}

func (h *ReserveHandler) GetReserve(c *gin.Context) {
	id := c.Param("id")
	reserve, ok := h.store.GetReserveRecord(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "预约记录不存在"})
		return
	}
	c.JSON(http.StatusOK, reserve)
}

func (h *ReserveHandler) ListReserves(c *gin.Context) {
	records := h.store.GetAllReserveRecords()
	c.JSON(http.StatusOK, records)
}

func (h *ReserveHandler) GetReserveQueue(c *gin.Context) {
	bookID := c.Param("bookId")
	book, ok := h.store.GetBook(bookID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "图书不存在"})
		return
	}

	queue := h.store.GetActiveReserveQueue(bookID)
	c.JSON(http.StatusOK, models.ReserveQueueInfo{
		BookID: bookID,
		Book:   book,
		Queue:  queue,
		Count:  len(queue),
	})
}

func (h *ReserveHandler) GetReaderReserves(c *gin.Context) {
	readerID := c.Param("readerId")
	records := h.store.GetReaderReserveRecords(readerID)
	c.JSON(http.StatusOK, records)
}

func (h *ReserveHandler) CancelReserve(c *gin.Context) {
	id := c.Param("id")
	reserve, ok := h.store.GetReserveRecord(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "预约记录不存在"})
		return
	}

	if reserve.Status != models.ReserveStatusWaiting && reserve.Status != models.ReserveStatusAvailable {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该预约无法取消"})
		return
	}

	wasAvailable := reserve.Status == models.ReserveStatusAvailable
	reserve.Status = models.ReserveStatusCancelled
	h.store.UpdateReserveRecord(reserve)

	if wasAvailable {
		next, ok := h.store.GetFirstWaitingReserve(reserve.BookID)
		if ok {
			h.store.NotifyNextReserveDirect(reserve.BookID, next)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "预约已取消"})
}

func (h *ReserveHandler) ExpireReserve(c *gin.Context) {
	id := c.Param("id")
	reserve, ok := h.store.GetReserveRecord(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "预约记录不存在"})
		return
	}

	if reserve.Status != models.ReserveStatusAvailable {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只有待取书状态的预约才能过期"})
		return
	}

	h.store.ExpireReserve(id)

	h.store.NotifyNextReserve(reserve.BookID)

	c.JSON(http.StatusOK, gin.H{"message": "预约已过期，已通知下一位预约者"})
}

func (h *ReserveHandler) CheckExpiredReserves(c *gin.Context) {
	now := time.Now()
	records := h.store.GetAllReserveRecords()
	expiredCount := 0

	for _, r := range records {
		if r.Status == models.ReserveStatusAvailable && r.ExpireDate != nil && now.After(*r.ExpireDate) {
			h.store.ExpireReserve(r.ID)
			h.store.NotifyNextReserve(r.BookID)
			expiredCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"expired_count": expiredCount,
		"message":       "已检查并处理过期预约",
	})
}
