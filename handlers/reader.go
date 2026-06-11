package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"library-borrow-system/models"
	"library-borrow-system/store"
	"library-borrow-system/utils"
)

type ReaderHandler struct {
	store *store.Store
}

func NewReaderHandler(s *store.Store) *ReaderHandler {
	return &ReaderHandler{store: s}
}

func (h *ReaderHandler) CreateReader(c *gin.Context) {
	var req models.ReaderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, ok := h.store.GetReaderByCardNo(req.CardNo); ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "读者证号已存在"})
		return
	}

	reader := &models.Reader{
		ID:       utils.GenerateReaderID(),
		Name:     req.Name,
		CardNo:   req.CardNo,
		Type:     req.Type,
		Phone:    req.Phone,
		Email:    req.Email,
	}

	h.store.CreateReader(reader)
	c.JSON(http.StatusCreated, reader)
}

func (h *ReaderHandler) GetReader(c *gin.Context) {
	id := c.Param("id")
	reader, ok := h.store.GetReader(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "读者不存在"})
		return
	}
	c.JSON(http.StatusOK, reader)
}

func (h *ReaderHandler) GetReaderByCardNo(c *gin.Context) {
	cardNo := c.Param("cardNo")
	reader, ok := h.store.GetReaderByCardNo(cardNo)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "读者不存在"})
		return
	}
	c.JSON(http.StatusOK, reader)
}

func (h *ReaderHandler) ListReaders(c *gin.Context) {
	readers := h.store.GetAllReaders()
	c.JSON(http.StatusOK, readers)
}

func (h *ReaderHandler) UpdateReader(c *gin.Context) {
	id := c.Param("id")
	reader, ok := h.store.GetReader(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "读者不存在"})
		return
	}

	var req models.ReaderUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		reader.Name = req.Name
	}
	if req.Type != "" {
		reader.Type = req.Type
	}
	if req.Phone != "" {
		reader.Phone = req.Phone
	}
	if req.Email != "" {
		reader.Email = req.Email
	}

	h.store.UpdateReader(reader)
	c.JSON(http.StatusOK, reader)
}

func (h *ReaderHandler) DeleteReader(c *gin.Context) {
	id := c.Param("id")
	if !h.store.DeleteReader(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "读者不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (h *ReaderHandler) GetReaderPrivileges(c *gin.Context) {
	privileges := make(map[string]models.ReaderPrivilege)
	for t, p := range models.ReaderPrivileges {
		privileges[string(t)] = p
	}
	c.JSON(http.StatusOK, privileges)
}

func (h *ReaderHandler) GetReaderStatus(c *gin.Context) {
	id := c.Param("id")
	reader, ok := h.store.GetReader(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "读者不存在"})
		return
	}

	privilege := models.ReaderPrivileges[reader.Type]
	borrowCount := h.store.GetReaderBorrowCount(id)
	unpaidFine := h.store.GetReaderUnpaidFine(id)
	isBlacklisted := h.store.IsBlacklisted(id)

	c.JSON(http.StatusOK, gin.H{
		"reader":              reader,
		"max_borrow_count":    privilege.MaxBorrowCount,
		"current_borrow_count": borrowCount,
		"remaining_quota":     privilege.MaxBorrowCount - borrowCount,
		"unpaid_fine":         unpaidFine,
		"is_blacklisted":      isBlacklisted,
		"max_borrow_days":     privilege.MaxBorrowDays,
		"max_renew_count":     privilege.MaxRenewCount,
	})
}
