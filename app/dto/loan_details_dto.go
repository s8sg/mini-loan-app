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
	LoanId           string              `json:"id" example:"b9348325-d798-4f81-85fc-336220380d4f"`
	CustomerId       string              `json:"customer-id" example:"user1"`
	TotalAmount      decimal.Decimal     `json:"total-amount" example:"100000"`
	Status           string              `json:"status" example:"PENDING"`
	Term             int                 `json:"term" example:"1"`
	Repayments       []*RepaymentDetails `json:"repayments"`
	StartDate        time.Time           `json:"start-date" example:"2023-03-10T09:58:40.009375Z"`
	CreatedTimestamp time.Time           `json:"created-timestamp" example:"2023-03-10T09:58:40.011177Z"`
	UpdatedTimestamp time.Time           `json:"updated-timestamp" example:"2023-03-10T09:58:40.011177Z"`
}

type RepaymentDetails struct {
	RepaymentId      string          `json:"id" example:"9b02d974-2b09-4e42-8006-5e94ee93659a"`
	Number           int             `json:"number" example:"1"`
	LoanId           string          `json:"loan-id,omitempty"`
	Amount           decimal.Decimal `json:"due-amount" example:"100000"`
	Status           string          `json:"status" example:"PENDING"`
	DueDate          time.Time       `json:"due-date" example:"2023-03-17T10:36:48.430739Z"`
	CreatedTimestamp time.Time       `json:"created-timestamp" example:"2023-03-10T10:36:48.431463Z"`
	UpdatedTimestamp time.Time       `json:"updated-timestamp" example:"2023-03-10T10:36:48.431463Z"`
}
