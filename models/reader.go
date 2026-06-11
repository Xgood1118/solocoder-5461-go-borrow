package models

import "time"

type ReaderType string

const (
	ReaderTypePrimaryStudent   ReaderType = "primary_student"
	ReaderTypeJuniorStudent    ReaderType = "junior_student"
	ReaderTypeSeniorStudent    ReaderType = "senior_student"
	ReaderTypeCollegeStudent   ReaderType = "college_student"
	ReaderTypeTeacher          ReaderType = "teacher"
	ReaderTypeSocial           ReaderType = "social"
	ReaderTypeVIP              ReaderType = "vip"
)

type ReaderPrivilege struct {
	MaxBorrowCount int
	MaxBorrowDays  int
	MaxRenewCount  int
}

var ReaderPrivileges = map[ReaderType]ReaderPrivilege{
	ReaderTypePrimaryStudent: {MaxBorrowCount: 3, MaxBorrowDays: 14, MaxRenewCount: 2},
	ReaderTypeJuniorStudent:  {MaxBorrowCount: 5, MaxBorrowDays: 21, MaxRenewCount: 2},
	ReaderTypeSeniorStudent:  {MaxBorrowCount: 8, MaxBorrowDays: 30, MaxRenewCount: 2},
	ReaderTypeCollegeStudent: {MaxBorrowCount: 10, MaxBorrowDays: 30, MaxRenewCount: 2},
	ReaderTypeTeacher:        {MaxBorrowCount: 30, MaxBorrowDays: 90, MaxRenewCount: 2},
	ReaderTypeSocial:         {MaxBorrowCount: 5, MaxBorrowDays: 30, MaxRenewCount: 2},
	ReaderTypeVIP:            {MaxBorrowCount: 20, MaxBorrowDays: 60, MaxRenewCount: 2},
}

type Reader struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	CardNo     string     `json:"card_no"`
	Type       ReaderType `json:"type"`
	Phone      string     `json:"phone"`
	Email      string     `json:"email"`
	IsBlacklisted bool    `json:"is_blacklisted"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type ReaderCreateRequest struct {
	Name   string     `json:"name" binding:"required"`
	CardNo string     `json:"card_no" binding:"required"`
	Type   ReaderType `json:"type" binding:"required"`
	Phone  string     `json:"phone"`
	Email  string     `json:"email"`
}

type ReaderUpdateRequest struct {
	Name  string     `json:"name"`
	Type  ReaderType `json:"type"`
	Phone string     `json:"phone"`
	Email string     `json:"email"`
}
