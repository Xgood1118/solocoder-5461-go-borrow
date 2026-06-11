package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"library-borrow-system/models"
	"library-borrow-system/store"
)

type FineHandler struct {
	store *store.Store
}

func NewFineHandler(s *store.Store) *FineHandler {
	return &FineHandler{store: s}
}

func (h *FineHandler) PayFine(c *gin.Context) {
	var req models.FinePayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reader, ok := h.store.GetReaderByCardNo(req.ReaderCardNo)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "读者不存在"})
		return
	}

	paidAmount, paidRecords := h.store.PayReaderFines(reader.ID, req.Amount)

	remainingFine := h.store.GetReaderUnpaidFine(reader.ID)

	if remainingFine <= 0 && reader.IsBlacklisted {
		h.store.RemoveFromBlacklist(reader.ID)
	}

	c.JSON(http.StatusOK, models.FinePayResponse{
		PaidAmount:    paidAmount,
		RemainingFine: remainingFine,
		Records:       paidRecords,
	})
}

func (h *FineHandler) GetReaderFines(c *gin.Context) {
	readerID := c.Param("readerId")
	records := h.store.GetReaderFineRecords(readerID)
	c.JSON(http.StatusOK, records)
}

func (h *FineHandler) GetReaderUnpaidFines(c *gin.Context) {
	readerID := c.Param("readerId")
	records := h.store.GetReaderUnpaidFines(readerID)

	var total float64
	for _, r := range records {
		total += r.Amount
	}

	c.JSON(http.StatusOK, gin.H{
		"records":     records,
		"total_unpaid": total,
		"count":       len(records),
	})
}

func (h *FineHandler) ListAllFines(c *gin.Context) {
	records := h.store.GetAllFineRecords()
	c.JSON(http.StatusOK, records)
}
