package dto

import "github.com/s8sg/mini-loan-app/app/dto"

type LoanCreateRequest struct {
	Amount float64 `json:"amount" example:"300000"`
	Term   int     `json:"term" example:"1"`
}

type LoanApproveRequest struct {
	LoanId string `json:"loan-id" example:"b9348325-d798-4f81-85fc-336220380d4f"`
}

type LoanRepaymentRequest struct {
	RepaymentID string  `json:"repayment-id" example:"393be183-ecc3-4a52-a035-f2e8a70d3711"`
	Amount      float64 `json:"amount" example:"300000"`
}

type GetAllLoansResponse struct {
	Loans []*dto.LoanDetails `json:"loans"`
}

type GenericSuccessResponse struct {
	Message string `json:"message" example:"successfully completed"`
}
