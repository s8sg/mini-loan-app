package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/s8sg/mini-loan-app/app/app_errors"
	"github.com/s8sg/mini-loan-app/app/controller/dto"
	repoDto "github.com/s8sg/mini-loan-app/app/dto"
	repository "github.com/s8sg/mini-loan-app/app/repostory"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

const (
	TimeoutInSecond = 5
)

var (
	repaymentNotFound      = &app_errors.AppError{Code: 404, Message: "repayment not found"}
	amountNotProvided      = &app_errors.AppError{Code: 400, Message: "amount must be provided"}
	repaymentIdNotProvided = &app_errors.AppError{Code: 400, Message: "repaymentId must be provided"}
	invalidLoanStatus      = &app_errors.AppError{Code: 400, Message: "invalid loan status"}
	invalidRepaymentStatus = &app_errors.AppError{Code: 400, Message: "invalid repayment status"}
	amountNotSufficient    = &app_errors.AppError{Code: 400, Message: "amount not sufficient"}
)

type RepaymentService interface {
	Repay(customerId string, request *dto.LoanRepaymentRequest) error
}

type RepaymentServiceImplementation struct {
	repo repository.LoanRepository
}

// GetRepaymentService :  Initialise repayment-service, uses dependency loanRepository
func GetRepaymentService(loanRepository repository.LoanRepository) RepaymentService {
	repaymentService := &RepaymentServiceImplementation{
		repo: loanRepository,
	}
	return repaymentService
}

func (r RepaymentServiceImplementation) Repay(customerId string, request *dto.LoanRepaymentRequest) error {

	if request.Amount == 0 {
		log.Println("amount must be provided")
		return amountNotProvided
	}

	if request.RepaymentID == "" {
		log.Println("repaymentId must be provided")
		return repaymentIdNotProvided
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), TimeoutInSecond*time.Second)
	defer cancelFunc()

	txOption := &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	}

	tx, err := r.repo.CreateTransaction(ctx, txOption)
	if err != nil {
		log.Println("failed to initiate transaction")
		return app_errors.InternalServerError
	}

	defer func() {
		if err != nil {
			log.Println("calling rollback for error " + err.Error())
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	repaymentDetails, err := r.repo.GetRepaymentById(request.RepaymentID, tx)
	if err != nil {
		log.Println("failed to fetch repayment, err: " + err.Error())
		err = fmt.Errorf("failed to fetch repayment, %v", err)
		return repaymentNotFound
	}

	loanID := repaymentDetails.LoanId
	loanDetails, err := r.repo.GetLoanById(loanID, tx)
	if err != nil {
		log.Println("failed to fetch loan, err: " + err.Error())
		err = fmt.Errorf("failed to fetch loam, %v", err)
		return app_errors.InternalServerError
	}

	// check if loan belongs to customer
	if loanDetails.CustomerId != customerId {
		log.Println("loan doesn't belongs to customer")
		err = fmt.Errorf("loan doesn't belongs to customer")
		return repaymentNotFound
	}

	// check the lean status
	if loanDetails.Status != repoDto.LoanStatusApproved {
		log.Println("loan status invalid")
		err = fmt.Errorf("loan has an invalid status %s", loanDetails.Status)
		return invalidLoanStatus
	}

	// check if the repayment status
	if repaymentDetails.Status != repoDto.RepaymentStatusPending {
		log.Println("repayment status invalid")
		err = fmt.Errorf("repaymentis already paid")
		return invalidRepaymentStatus
	}

	// check of the repayment amount >= due amount
	if decimal.NewFromFloat(request.Amount).LessThan(repaymentDetails.Amount) {
		log.Println("invalid amount paid")
		err = fmt.Errorf("repaymentis paid with invalid amount")
		return amountNotSufficient
	}

	err = r.repo.UpdateRepaymentStatus(repaymentDetails.RepaymentId, repoDto.RepaymentStatusPaid, tx)
	if err != nil {
		log.Println("failed to update repayment, error " + err.Error())
		return app_errors.InternalServerError
	}

	repaidRepayments := getRepaidRepaymentCount(loanDetails)

	// check if all repayments are being paid
	// mark the loan as paid
	if repaidRepayments+1 == loanDetails.Term {
		err = r.repo.UpdateLoanStatus(loanID, repoDto.LOAN_STATUS_PAID, tx)
		if err != nil {
			log.Println("failed tp update loan status")
			return app_errors.InternalServerError
		}
	}

	return nil
}

func getRepaidRepaymentCount(loanDetails *repoDto.LoanDetails) int {
	repaidRepayments := 0
	for _, repayment := range loanDetails.Repayments {
		if repayment.Status == repoDto.RepaymentStatusPaid {
			repaidRepayments = repaidRepayments + 1
		}
	}
	return repaidRepayments
}
