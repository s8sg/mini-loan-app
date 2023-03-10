package dto

type LoanCreateRequest struct {
	Amount float64 `json:"amount"`
	Term   int     `json:"term"`
}

type LoanApproveRequest struct {
	LoanId string `json:"loan-id"`
}

type LoanRepaymentRequest struct {
	RepaymentID string  `json:"repayment-id"`
	Amount      float64 `json:"amount"`
}
