package util

import "github.com/google/uuid"

func GenerateLoanID() string {
	return uuid.New().String()
}

func GenerateRepaymentID() string {
	return uuid.New().String()
}
