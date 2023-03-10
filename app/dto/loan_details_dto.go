package dto

import (
	"github.com/shopspring/decimal"
	"time"
)

const (
	LoanStatusPending  = "PENDING"
	LoanStatusApproved = "APPROVED"
	LOAN_STATUS_PAID   = "PAID"
)

const (
	RepaymentStatusPending = "PENDING"
	RepaymentStatusPaid    = "PAID"
)

type LoanDetails struct {
	LoanId           string              `json:"id"`
	CustomerId       string              `json:"customer-id"`
	TotalAmount      decimal.Decimal     `json:"total-amount"`
	Status           string              `json:"status"`
	Term             int                 `json:"term""`
	Repayments       []*RepaymentDetails `json:"repayments"`
	StartDate        time.Time           `json:"start-date"`
	CreatedTimestamp time.Time           `json:"created-timestamp"`
	UpdatedTimestamp time.Time           `json:"updated-timestamp"`
}

type RepaymentDetails struct {
	RepaymentId      string          `json:"id"`
	Number           int             `json:"number"`
	LoanId           string          `json:"loan-id,omitempty"`
	Amount           decimal.Decimal `json:"due-amount"`
	Status           string          `json:"status"`
	DueDate          time.Time       `json:"due-date"`
	CreatedTimestamp time.Time       `json:"created-timestamp"`
	UpdatedTimestamp time.Time       `json:"updated-timestamp"`
}
