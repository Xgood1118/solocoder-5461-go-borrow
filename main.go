package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"library-borrow-system/handlers"
	"library-borrow-system/store"
)

func setupRouter(s *store.Store) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	readerHandler := handlers.NewReaderHandler(s)
	bookHandler := handlers.NewBookHandler(s)
	borrowHandler := handlers.NewBorrowHandler(s)
	renewHandler := handlers.NewRenewHandler(s)
	reserveHandler := handlers.NewReserveHandler(s)
	fineHandler := handlers.NewFineHandler(s)
	blacklistHandler := handlers.NewBlacklistHandler(s)
	statsHandler := handlers.NewStatsHandler(s)

	api := r.Group("/api")
	{
		reader := api.Group("/readers")
		{
			reader.POST("", readerHandler.CreateReader)
			reader.GET("", readerHandler.ListReaders)
			reader.GET("/privileges", readerHandler.GetReaderPrivileges)
			reader.GET("/:id", readerHandler.GetReader)
			reader.GET("/card/:cardNo", readerHandler.GetReaderByCardNo)
			reader.GET("/:id/status", readerHandler.GetReaderStatus)
			reader.PUT("/:id", readerHandler.UpdateReader)
			reader.DELETE("/:id", readerHandler.DeleteReader)
		}

		book := api.Group("/books")
		{
			book.POST("", bookHandler.CreateBook)
			book.GET("", bookHandler.ListBooks)
			book.GET("/:id", bookHandler.GetBook)
			book.GET("/rfid/:rfid", bookHandler.GetBookByRFID)
			book.PUT("/:id", bookHandler.UpdateBook)
			book.DELETE("/:id", bookHandler.DeleteBook)
			book.POST("/:id/lost", bookHandler.MarkBookLost)
		}

		borrow := api.Group("/borrow")
		{
			borrow.POST("", borrowHandler.BorrowBook)
			borrow.GET("", borrowHandler.ListBorrowRecords)
			borrow.GET("/:id", borrowHandler.GetBorrowRecord)
			borrow.GET("/reader/:readerId", borrowHandler.GetReaderBorrowRecords)
			borrow.GET("/book/:bookId", borrowHandler.GetBookBorrowRecords)
		}

		returnGroup := api.Group("/return")
		{
			returnGroup.POST("", borrowHandler.ReturnBook)
		}

		renew := api.Group("/renew")
		{
			renew.POST("", renewHandler.RenewBook)
			renew.GET("/check", renewHandler.CanRenew)
		}

		reserve := api.Group("/reserves")
		{
			reserve.POST("", reserveHandler.CreateReserve)
			reserve.GET("", reserveHandler.ListReserves)
			reserve.GET("/:id", reserveHandler.GetReserve)
			reserve.GET("/book/:bookId/queue", reserveHandler.GetReserveQueue)
			reserve.GET("/reader/:readerId", reserveHandler.GetReaderReserves)
			reserve.DELETE("/:id", reserveHandler.CancelReserve)
			reserve.POST("/:id/expire", reserveHandler.ExpireReserve)
			reserve.POST("/check-expired", reserveHandler.CheckExpiredReserves)
		}

		fine := api.Group("/fines")
		{
			fine.POST("/pay", fineHandler.PayFine)
			fine.GET("", fineHandler.ListAllFines)
			fine.GET("/reader/:readerId", fineHandler.GetReaderFines)
			fine.GET("/reader/:readerId/unpaid", fineHandler.GetReaderUnpaidFines)
		}

		blacklist := api.Group("/blacklist")
		{
			blacklist.GET("", blacklistHandler.ListBlacklist)
			blacklist.GET("/:readerId", blacklistHandler.GetBlacklistStatus)
			blacklist.DELETE("/:readerId", blacklistHandler.RemoveFromBlacklist)
		}

		stats := api.Group("/stats")
		{
			stats.GET("/overview", statsHandler.GetOverview)
			stats.GET("/monthly", statsHandler.GetMonthlyStats)
			stats.GET("/monthly/by-month", statsHandler.GetStatsByMonth)
		}
	}

	return r
}

func main() {
	s := store.GlobalStore

	InitSampleData(s)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := setupRouter(s)

	log.Printf("图书馆借阅管理系统启动成功，端口: %s", port)
	log.Printf("API 文档: http://localhost:%s/api", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
