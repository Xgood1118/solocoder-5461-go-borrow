package models

type MonthlyStats struct {
	Year                 int                      `json:"year"`
	Month                int                      `json:"month"`
	TotalBorrowCount     int                      `json:"total_borrow_count"`
	TotalReturnCount     int                      `json:"total_return_count"`
	TotalReserveCount    int                      `json:"total_reserve_count"`
	TotalFineAmount      float64                  `json:"total_fine_amount"`
	ReserveConvertRate   float64                  `json:"reserve_convert_rate"`
	ReaderTypeStats      map[ReaderType]int       `json:"reader_type_stats"`
	TopBooks             []*BookRanking           `json:"top_books"`
}

type BookRanking struct {
	BookID      string  `json:"book_id"`
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	BorrowCount int     `json:"borrow_count"`
}

type StatsResponse struct {
	Monthly *MonthlyStats `json:"monthly"`
}
