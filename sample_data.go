package main

import (
	"library-borrow-system/models"
	"library-borrow-system/store"
	"library-borrow-system/utils"
)

func InitSampleData(s *store.Store) {
	readers := []*models.Reader{
		{ID: utils.GenerateReaderID(), Name: "张小明", CardNo: "STU001", Type: models.ReaderTypePrimaryStudent, Phone: "13800138001"},
		{ID: utils.GenerateReaderID(), Name: "李小红", CardNo: "STU002", Type: models.ReaderTypeJuniorStudent, Phone: "13800138002"},
		{ID: utils.GenerateReaderID(), Name: "王小华", CardNo: "STU003", Type: models.ReaderTypeSeniorStudent, Phone: "13800138003"},
		{ID: utils.GenerateReaderID(), Name: "赵小刚", CardNo: "STU004", Type: models.ReaderTypeCollegeStudent, Phone: "13800138004"},
		{ID: utils.GenerateReaderID(), Name: "陈老师", CardNo: "TEA001", Type: models.ReaderTypeTeacher, Phone: "13900139001"},
		{ID: utils.GenerateReaderID(), Name: "刘大爷", CardNo: "SOC001", Type: models.ReaderTypeSocial, Phone: "13700137001"},
		{ID: utils.GenerateReaderID(), Name: "孙女士", CardNo: "VIP001", Type: models.ReaderTypeVIP, Phone: "13600136001"},
	}

	for _, r := range readers {
		s.CreateReader(r)
	}

	books := []*models.Book{
		{ID: utils.GenerateBookID(), Title: "活着", Author: "余华", ISBN: "9787506365437", RFID: "RFID001", Price: 35.00, Category: "文学", Location: "A区-01架"},
		{ID: utils.GenerateBookID(), Title: "三体", Author: "刘慈欣", ISBN: "9787536692930", RFID: "RFID002", Price: 23.00, Category: "科幻", Location: "B区-03架"},
		{ID: utils.GenerateBookID(), Title: "百年孤独", Author: "加西亚·马尔克斯", ISBN: "9787544253994", RFID: "RFID003", Price: 39.50, Category: "文学", Location: "A区-02架"},
		{ID: utils.GenerateBookID(), Title: "平凡的世界", Author: "路遥", ISBN: "9787530209555", RFID: "RFID004", Price: 68.00, Category: "文学", Location: "A区-03架"},
		{ID: utils.GenerateBookID(), Title: "小王子", Author: "圣埃克苏佩里", ISBN: "9787020042494", RFID: "RFID005", Price: 22.00, Category: "童话", Location: "C区-01架"},
		{ID: utils.GenerateBookID(), Title: "红楼梦", Author: "曹雪芹", ISBN: "9787020002207", RFID: "RFID006", Price: 59.70, Category: "古典", Location: "D区-01架"},
		{ID: utils.GenerateBookID(), Title: "人类简史", Author: "尤瓦尔·赫拉利", ISBN: "9787508647357", RFID: "RFID007", Price: 68.00, Category: "历史", Location: "E区-02架"},
		{ID: utils.GenerateBookID(), Title: "围城", Author: "钱钟书", ISBN: "9787020024759", RFID: "RFID008", Price: 19.00, Category: "文学", Location: "A区-05架"},
		{ID: utils.GenerateBookID(), Title: "明朝那些事儿", Author: "当年明月", ISBN: "9787213046438", RFID: "RFID009", Price: 358.00, Category: "历史", Location: "E区-03架"},
		{ID: utils.GenerateBookID(), Title: "白夜行", Author: "东野圭吾", ISBN: "9787544258609", RFID: "RFID010", Price: 39.50, Category: "推理", Location: "F区-01架"},
	}

	for _, b := range books {
		b.Status = models.BookStatusAvailable
		s.CreateBook(b)
	}
}
