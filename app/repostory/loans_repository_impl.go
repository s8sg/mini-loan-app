package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/s8sg/mini-loan-app/app/dto"
	"github.com/s8sg/mini-loan-app/app/util"
	"log"
	"time"
)

const (
	TimeoutInSecond = 5
)

type SqlLoanRepository struct {
	*sql.DB
}

// GetLoanRepository : factory function initialize SqlLoanRepository
func GetLoanRepository(db *sql.DB) LoanRepository {
	loanRepository := &SqlLoanRepository{
		DB: db,
	}
	return loanRepository
}

func (db *SqlLoanRepository) CreateLoan(loanDetails *dto.LoanDetails) (*dto.LoanDetails, error) {

	ctx, cancelFunc := context.WithTimeout(context.Background(), TimeoutInSecond*time.Second)
	defer cancelFunc()

	option := &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	}
	tx, err := db.BeginTx(ctx, option)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactio: %w", err)
	}

	defer func() {
		if err != nil {
			log.Println("calling rollback for error " + err.Error())
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	query := "INSERT INTO loans (id, customer_id, amount, term, status, start_date) VALUES ($1, $2, $3, $4, $5, $6)"

	res, err := tx.ExecContext(ctx, query, loanDetails.LoanId, loanDetails.CustomerId, loanDetails.TotalAmount,
		loanDetails.Term, loanDetails.Status, loanDetails.StartDate)
	if err != nil {
		log.Printf("Error %s when inserting row into loans table", err)
		return nil, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if count != 1 {
		err = fmt.Errorf("no rows updated when inserting row into loans table")
		return nil, err
	}

	for _, repayment := range loanDetails.Repayments {
		query = "INSERT INTO repayments(id, num, loan_id, amount, status, due_date) VALUES ($1, $2, $3, $4, $5, $6)"

		_, err = tx.ExecContext(ctx, query, repayment.RepaymentId, repayment.Number, loanDetails.LoanId, repayment.Amount, repayment.Status, repayment.DueDate)
		if err != nil {
			log.Printf("Error %s when inserting row into repayments table", err)
			return nil, err
		}
		count, err = res.RowsAffected()
		if err != nil {
			return nil, err
		}

		if count != 1 {
			err = fmt.Errorf("no rows updated when inserting row into repayment table")
			return nil, err
		}
	}

	return loanDetails, nil
}

func (db *SqlLoanRepository) GetAllLoansByCustomerId(customer string) ([]*dto.LoanDetails, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), TimeoutInSecond*time.Second)
	defer cancelFunc()

	// TODO: This can later be done with a single query with join statement

	query := "SELECT id, customer_id, amount, term, status, start_date, created_at, updated_at FROM loans WHERE customer_id = $1"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, customer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loanDetailsList := make([]*dto.LoanDetails, 0)

	for rows.Next() {
		loanDetails := &dto.LoanDetails{}
		if err := rows.Scan(&loanDetails.LoanId, &loanDetails.CustomerId, &loanDetails.TotalAmount, &loanDetails.Term,
			&loanDetails.Status, &loanDetails.StartDate, &loanDetails.CreatedTimestamp, &loanDetails.UpdatedTimestamp); err != nil {
			return nil, err
		}

		query = "SELECT id, num, amount, status, due_date, created_at, updated_at FROM repayments WHERE loan_id = $1"
		stmt2, err := db.PrepareContext(ctx, query)
		if err != nil {
			log.Printf("Error %s when preparing SQL statement", err)
			return nil, err
		}
		defer stmt2.Close()
		rows2, err := stmt2.QueryContext(ctx, loanDetails.LoanId)
		if err != nil {
			return nil, err
		}
		defer rows2.Close()

		repaymentDetailsList := make([]*dto.RepaymentDetails, 0)

		for rows2.Next() {
			repaymentDetails := &dto.RepaymentDetails{}
			if err := rows2.Scan(&repaymentDetails.RepaymentId, &repaymentDetails.Number, &repaymentDetails.Amount,
				&repaymentDetails.Status, &repaymentDetails.DueDate, &repaymentDetails.CreatedTimestamp,
				&repaymentDetails.UpdatedTimestamp); err != nil {
				return nil, err
			}
			repaymentDetailsList = append(repaymentDetailsList, repaymentDetails)
		}
		if err := rows2.Err(); err != nil {
			return nil, err
		}

		loanDetails.Repayments = repaymentDetailsList

		loanDetailsList = append(loanDetailsList, loanDetails)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return loanDetailsList, nil
}

func (db *SqlLoanRepository) UpdateLoanStatus(loanId string, status string, transactionalContext *Transaction) error {

	query := "UPDATE loans set status = $1, updated_at = $2 WHERE id = $3"
	res, err := transactionalContext.tx.ExecContext(transactionalContext.ctx, query, status, util.GetCurrentTimeInUtc(), loanId)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count != 1 {
		err = fmt.Errorf("no rows updated")
		return err
	}

	return nil
}

func (db *SqlLoanRepository) GetLoanById(loanId string, transactionalContext *Transaction) (*dto.LoanDetails, error) {
	query := "SELECT id, customer_id, amount, term, status, start_date, created_at, updated_at FROM loans WHERE id = $1"
	row := transactionalContext.tx.QueryRowContext(transactionalContext.ctx, query, loanId)
	loanDetails := &dto.LoanDetails{}
	if err := row.Scan(&loanDetails.LoanId, &loanDetails.CustomerId, &loanDetails.TotalAmount, &loanDetails.Term,
		&loanDetails.Status, &loanDetails.StartDate, &loanDetails.CreatedTimestamp, &loanDetails.UpdatedTimestamp); err != nil {
		return nil, err
	}

	// TODO: This can later be done with a single query with join statement
	repaymentDetailsList, err := db.GetRepaymentsByLoanId(loanId, transactionalContext)
	if err != nil {
		return nil, err
	}
	loanDetails.Repayments = repaymentDetailsList

	return loanDetails, nil
}

func (db *SqlLoanRepository) GetRepaymentsByLoanId(loanId string, transactionalContext *Transaction) ([]*dto.RepaymentDetails, error) {
	query := "SELECT id, num, amount, status, due_date, created_at, updated_at FROM repayments WHERE loan_id = $1"
	stmt, err := transactionalContext.tx.PrepareContext(transactionalContext.ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(transactionalContext.ctx, loanId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	repaymentDetailsList := make([]*dto.RepaymentDetails, 0)
	for rows.Next() {
		repaymentDetails := &dto.RepaymentDetails{}
		if err := rows.Scan(&repaymentDetails.RepaymentId, &repaymentDetails.Number, &repaymentDetails.Amount, &repaymentDetails.Status,
			&repaymentDetails.DueDate, &repaymentDetails.CreatedTimestamp, &repaymentDetails.UpdatedTimestamp); err != nil {
			return nil, err
		}
		repaymentDetailsList = append(repaymentDetailsList, repaymentDetails)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return repaymentDetailsList, nil
}

func (db *SqlLoanRepository) GetRepaymentById(repaymentId string, transactionalContext *Transaction) (*dto.RepaymentDetails, error) {
	query := "SELECT id, num, loan_id, amount, status, due_date, created_at, updated_at FROM repayments WHERE id = $1"
	row := transactionalContext.tx.QueryRowContext(transactionalContext.ctx, query, repaymentId)
	repaymentDetails := &dto.RepaymentDetails{}
	if err := row.Scan(&repaymentDetails.RepaymentId, &repaymentDetails.Number, &repaymentDetails.LoanId, &repaymentDetails.Amount, &repaymentDetails.Status,
		&repaymentDetails.DueDate, &repaymentDetails.CreatedTimestamp, &repaymentDetails.UpdatedTimestamp); err != nil {
		return nil, err

	}
	return repaymentDetails, nil
}

func (db *SqlLoanRepository) UpdateRepaymentStatus(repaymentId string, status string, transactionalContext *Transaction) error {
	query := "UPDATE repayments set status = $1, updated_at = $2 WHERE id = $3"
	res, err := transactionalContext.tx.ExecContext(transactionalContext.ctx, query, status, util.GetCurrentTimeInUtc(), repaymentId)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count != 1 {
		err = fmt.Errorf("no rows updated")
		return err
	}

	return nil
}

func (db *SqlLoanRepository) CreateTransaction(ctx context.Context, opts *sql.TxOptions) (*Transaction, error) {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactio: %w", err)
	}
	return &Transaction{
		tx:  tx,
		ctx: ctx,
	}, nil
}
