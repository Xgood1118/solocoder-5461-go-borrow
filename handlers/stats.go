package handlers

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"library-borrow-system/models"
	"library-borrow-system/store"
)

type StatsHandler struct {
	store *store.Store
}

func NewStatsHandler(s *store.Store) *StatsHandler {
	return &StatsHandler{store: s}
}

func (h *StatsHandler) GetMonthlyStats(c *gin.Context) {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	stats := h.calculateMonthlyStats(year, month)
	c.JSON(http.StatusOK, stats)
}

func (h *StatsHandler) GetStatsByMonth(c *gin.Context) {
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	var year, month int
	now := time.Now()

	if yearStr != "" {
		year = parseInt(yearStr)
	} else {
		year = now.Year()
	}
	if monthStr != "" {
		month = parseInt(monthStr)
	} else {
		month = int(now.Month())
	}

	if month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的月份"})
		return
	}

	stats := h.calculateMonthlyStats(year, month)
	c.JSON(http.StatusOK, stats)
}

func (h *StatsHandler) calculateMonthlyStats(year, month int) *models.MonthlyStats {
	stats := &models.MonthlyStats{
		Year:            year,
		Month:           month,
		ReaderTypeStats: make(map[models.ReaderType]int),
	}

	borrowRecords := h.store.GetAllBorrowRecords()
	returnRecords := h.store.GetAllBorrowRecords()
	reserveRecords := h.store.GetAllReserveRecords()
	fineRecords := h.store.GetAllFineRecords()

	monthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	monthEnd := monthStart.AddDate(0, 1, 0)

	bookBorrowCount := make(map[string]int)

	for _, r := range borrowRecords {
		if r.BorrowDate.After(monthStart) && r.BorrowDate.Before(monthEnd) {
			stats.TotalBorrowCount++
			bookBorrowCount[r.BookID]++

			reader, ok := h.store.GetReader(r.ReaderID)
			if ok {
				stats.ReaderTypeStats[reader.Type]++
			}
		}
	}

	for _, r := range returnRecords {
		if r.ReturnDate != nil && r.ReturnDate.After(monthStart) && r.ReturnDate.Before(monthEnd) {
			stats.TotalReturnCount++
		}
	}

	for _, r := range reserveRecords {
		if r.ReserveDate.After(monthStart) && r.ReserveDate.Before(monthEnd) {
			stats.TotalReserveCount++
		}
	}

	for _, r := range fineRecords {
		if r.CreatedAt.After(monthStart) && r.CreatedAt.Before(monthEnd) {
			stats.TotalFineAmount += r.Amount
		}
	}

	if stats.TotalReserveCount > 0 {
		convertedCount := 0
		for _, r := range reserveRecords {
			if r.ConvertedBorrowID != nil {
				convertedCount++
			}
		}
		stats.ReserveConvertRate = float64(convertedCount) / float64(stats.TotalReserveCount)
	}

	topBooks := h.getTopBooks(bookBorrowCount, 20)
	stats.TopBooks = topBooks

	return stats
}

func (h *StatsHandler) getTopBooks(bookCount map[string]int, limit int) []*models.BookRanking {
	type bookCountPair struct {
		bookID string
		count  int
	}

	pairs := make([]bookCountPair, 0, len(bookCount))
	for id, count := range bookCount {
		pairs = append(pairs, bookCountPair{bookID: id, count: count})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	if len(pairs) > limit {
		pairs = pairs[:limit]
	}

	rankings := make([]*models.BookRanking, 0, len(pairs))
	for _, p := range pairs {
		book, ok := h.store.GetBook(p.bookID)
		if ok {
			rankings = append(rankings, &models.BookRanking{
				BookID:      book.ID,
				Title:       book.Title,
				Author:      book.Author,
				BorrowCount: p.count,
			})
		}
	}

	return rankings
}

func (h *StatsHandler) GetOverview(c *gin.Context) {
	readers := h.store.GetAllReaders()
	books := h.store.GetAllBooks()
	borrowRecords := h.store.GetAllBorrowRecords()
	reserveRecords := h.store.GetAllReserveRecords()
	fineRecords := h.store.GetAllFineRecords()

	activeBorrowCount := 0
	for _, r := range borrowRecords {
		if r.ReturnDate == nil {
			activeBorrowCount++
		}
	}

	totalFineAmount := 0.0
	unpaidFineAmount := 0.0
	for _, r := range fineRecords {
		totalFineAmount += r.Amount
		if r.Status == models.FineStatusUnpaid {
			unpaidFineAmount += r.Amount
		}
	}

	activeReserveCount := 0
	for _, r := range reserveRecords {
		if r.Status == models.ReserveStatusWaiting || r.Status == models.ReserveStatusAvailable {
			activeReserveCount++
		}
	}

	availableBookCount := 0
	for _, b := range books {
		if b.Status == models.BookStatusAvailable {
			availableBookCount++
		}
	}

	blacklistCount := 0
	blacklistRecords := h.store.GetAllBlacklistRecords()
	for _, r := range blacklistRecords {
		if r.IsActive {
			blacklistCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_readers":         len(readers),
		"total_books":           len(books),
		"available_books":       availableBookCount,
		"active_borrow_count":   activeBorrowCount,
		"total_borrow_records":  len(borrowRecords),
		"active_reserve_count":  activeReserveCount,
		"total_reserve_records": len(reserveRecords),
		"total_fine_amount":     totalFineAmount,
		"unpaid_fine_amount":    unpaidFineAmount,
		"blacklist_count":       blacklistCount,
	})
}

func parseInt(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}
