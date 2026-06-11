package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateID(prefix string) string {
	return fmt.Sprintf("%s_%d_%06d", prefix, time.Now().UnixNano(), rand.Intn(1000000))
}

func GenerateReaderID() string {
	return GenerateID("R")
}

func GenerateBookID() string {
	return GenerateID("B")
}

func GenerateBorrowID() string {
	return GenerateID("BR")
}

func GenerateReserveID() string {
	return GenerateID("RS")
}

func GenerateFineID() string {
	return GenerateID("F")
}

func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
