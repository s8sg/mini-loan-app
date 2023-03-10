package repository

import (
	"context"
	"database/sql"
	"github.com/s8sg/mini-loan-app/app/dto"
)

type LoanRepository interface {
	CreateLoan(loanDetails *dto.LoanDetails) (*dto.LoanDetails, error)

	GetAllLoansByCustomerId(customerId string) ([]*dto.LoanDetails, error)

	GetLoanById(loanId string, transactionalContext *Transaction) (*dto.LoanDetails, error)

	UpdateLoanStatus(loanId string, status string, transactionalContext *Transaction) error

	GetRepaymentsByLoanId(loanId string, transactionalContext *Transaction) ([]*dto.RepaymentDetails, error)

	GetRepaymentById(repaymentId string, transactionalContext *Transaction) (*dto.RepaymentDetails, error)

	UpdateRepaymentStatus(id string, status string, tx *Transaction) error

	CreateTransaction(ctx context.Context, opts *sql.TxOptions) (*Transaction, error)
}

type Transaction struct {
	ctx context.Context
	tx  *sql.Tx
}

func (t *Transaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *Transaction) Commit() error {
	return t.tx.Commit()
}
