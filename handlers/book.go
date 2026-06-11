package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"library-borrow-system/models"
	"library-borrow-system/store"
	"library-borrow-system/utils"
)

type BookHandler struct {
	store *store.Store
}

func NewBookHandler(s *store.Store) *BookHandler {
	return &BookHandler{store: s}
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var req models.BookCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, ok := h.store.GetBookByRFID(req.RFID); ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图书RFID已存在"})
		return
	}

	book := &models.Book{
		ID:        utils.GenerateBookID(),
		ISBN:      req.ISBN,
		Title:     req.Title,
		Author:    req.Author,
		Publisher: req.Publisher,
		RFID:      req.RFID,
		Price:     req.Price,
		Location:  req.Location,
		Category:  req.Category,
		Status:    models.BookStatusAvailable,
	}

	h.store.CreateBook(book)
	c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) GetBook(c *gin.Context) {
	id := c.Param("id")
	book, ok := h.store.GetBook(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "图书不存在"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) GetBookByRFID(c *gin.Context) {
	rfid := c.Param("rfid")
	book, ok := h.store.GetBookByRFID(rfid)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "图书不存在"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) ListBooks(c *gin.Context) {
	books := h.store.GetAllBooks()
	c.JSON(http.StatusOK, books)
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	id := c.Param("id")
	book, ok := h.store.GetBook(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "图书不存在"})
		return
	}

	var req models.BookUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ISBN != "" {
		book.ISBN = req.ISBN
	}
	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Author != "" {
		book.Author = req.Author
	}
	if req.Publisher != "" {
		book.Publisher = req.Publisher
	}
	if req.Price > 0 {
		book.Price = req.Price
	}
	if req.Location != "" {
		book.Location = req.Location
	}
	if req.Category != "" {
		book.Category = req.Category
	}

	h.store.UpdateBook(book)
	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	id := c.Param("id")
	if !h.store.DeleteBook(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "图书不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (h *BookHandler) MarkBookLost(c *gin.Context) {
	id := c.Param("id")
	book, ok := h.store.GetBook(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "图书不存在"})
		return
	}

	compensation := book.Price * 3
	h.store.UpdateBookStatus(id, models.BookStatusLost)

	if borrowRecord, ok := h.store.GetActiveBorrowByBookID(id); ok {
		now := time.Now()
		borrowRecord.ReturnDate = &now
		borrowRecord.IsOverdue = false
		borrowRecord.OverdueDays = 0
		borrowRecord.FineAmount = compensation
		borrowRecord.IsFinePaid = false
		h.store.UpdateBorrowRecord(borrowRecord)

		fineRecord := &models.FineRecord{
			ID:          utils.GenerateFineID(),
			ReaderID:    borrowRecord.ReaderID,
			BorrowID:    borrowRecord.ID,
			BookID:      id,
			Amount:      compensation,
			Reason:      "图书丢失赔偿",
			OverdueDays: 0,
		}
		h.store.CreateFineRecord(fineRecord)

		totalFine := h.store.GetReaderUnpaidFine(borrowRecord.ReaderID)
		if totalFine > 30 {
			h.store.AddToBlacklist(borrowRecord.ReaderID, "图书丢失未赔偿", totalFine)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"book":          book,
		"compensation":  compensation,
		"message":       "图书已标记为丢失，需按原价3倍赔偿",
	})
}
