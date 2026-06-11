package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"library-borrow-system/store"
)

type BlacklistHandler struct {
	store *store.Store
}

func NewBlacklistHandler(s *store.Store) *BlacklistHandler {
	return &BlacklistHandler{store: s}
}

func (h *BlacklistHandler) ListBlacklist(c *gin.Context) {
	records := h.store.GetAllBlacklistRecords()
	active := make([]interface{}, 0)
	for _, r := range records {
		if r.IsActive {
			active = append(active, r)
		}
	}
	c.JSON(http.StatusOK, active)
}

func (h *BlacklistHandler) GetBlacklistStatus(c *gin.Context) {
	readerID := c.Param("readerId")
	record, ok := h.store.GetBlacklistRecord(readerID)
	if !ok || !record.IsActive {
		c.JSON(http.StatusOK, gin.H{"is_blacklisted": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"is_blacklisted": true,
		"record":         record,
	})
}

func (h *BlacklistHandler) RemoveFromBlacklist(c *gin.Context) {
	readerID := c.Param("readerId")

	reader, ok := h.store.GetReader(readerID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "读者不存在"})
		return
	}

	unpaidFine := h.store.GetReaderUnpaidFine(readerID)
	if unpaidFine > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "读者仍有未结清罚款，无法移出黑名单"})
		return
	}

	h.store.RemoveFromBlacklist(readerID)

	c.JSON(http.StatusOK, gin.H{
		"message": "已移出黑名单",
		"reader":  reader,
	})
}
