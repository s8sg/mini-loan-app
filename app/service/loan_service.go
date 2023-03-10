package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/s8sg/mini-loan-app/app/app_errors"
	"github.com/s8sg/mini-loan-app/app/controller/dto"
	responseDto "github.com/s8sg/mini-loan-app/app/dto"
	repository "github.com/s8sg/mini-loan-app/app/repostory"
	"github.com/s8sg/mini-loan-app/app/util"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

var (
	// RepaymentFrequency is not fixed to 1 week
	RepaymentFrequency = (time.Hour * 24) * 7
)

var (
	loanInvalidStatus    = &app_errors.AppError{Code: 400, Message: "loan invalid status"}
	loanNotPresent       = &app_errors.AppError{Code: 404, Message: "loan not found"}
	loanAmountNotPresent = &app_errors.AppError{Code: 400, Message: "loan amount must be provided"}
	loanTermInvalid      = &app_errors.AppError{Code: 400, Message: "loan term can;t be less than 1"}
	invalidLoanId        = &app_errors.AppError{Code: 400, Message: "invalid loan id"}
)

type LoanService interface {
	CreateLoan(customerId string, loanCreateRequest *dto.LoanCreateRequest) (*responseDto.LoanDetails, error)
	GetAllLoansForCustomer(customerId string) ([]*responseDto.LoanDetails, error)
	ApproveLoan(loanApproveRequest *dto.LoanApproveRequest) error
}

type LoanServiceImplementation struct {
	repo repository.LoanRepository
}

// GetLoanService : Initialise loan-service, uses dependency loanRepository
func GetLoanService(loanRepository repository.LoanRepository) LoanService {
	loanServiceImpl := &LoanServiceImplementation{
		repo: loanRepository,
	}
	return loanServiceImpl
}

func (l LoanServiceImplementation) CreateLoan(customerId string,
	loanCreateRequest *dto.LoanCreateRequest) (*responseDto.LoanDetails, error) {

	// validate amount
	if loanCreateRequest.Amount == 0 {
		log.Printf("loan must be present")
		return nil, loanAmountNotPresent
	}

	// validate term
	if loanCreateRequest.Term < 1 {
		log.Printf("term must be provide and should be greater than 1")
		return nil, loanTermInvalid
	}

	// create loan details
	loanDetails := &responseDto.LoanDetails{
		LoanId:           util.GenerateLoanID(),
		TotalAmount:      decimal.NewFromFloat(loanCreateRequest.Amount),
		CustomerId:       customerId,
		Term:             loanCreateRequest.Term,
		StartDate:        util.GetCurrentTimeInUtc(),
		Repayments:       make([]*responseDto.RepaymentDetails, loanCreateRequest.Term),
		Status:           responseDto.LoanStatusPending,
		CreatedTimestamp: util.GetCurrentTimeInUtc(),
		UpdatedTimestamp: util.GetCurrentTimeInUtc(),
	}

	// generate repayment details
	repaymentAmountPerTenure := loanDetails.TotalAmount.Div(decimal.NewFromInt32(int32(loanDetails.Term)))
	nextDueDate := loanDetails.StartDate
	for i := 0; i < loanDetails.Term; i++ {
		nextDueDate = nextDueDate.Add(RepaymentFrequency)
		repayment := &responseDto.RepaymentDetails{
			RepaymentId:      util.GenerateRepaymentID(),
			Number:           i + 1,
			Amount:           repaymentAmountPerTenure,
			DueDate:          nextDueDate,
			Status:           responseDto.RepaymentStatusPending,
			CreatedTimestamp: util.GetCurrentTimeInUtc(),
			UpdatedTimestamp: util.GetCurrentTimeInUtc(),
		}

		loanDetails.Repayments[i] = repayment
	}

	loanDetails, err := l.repo.CreateLoan(loanDetails)
	if err != nil {
		log.Printf("failed to create loan, error %v\n", err)
		return nil, app_errors.InternalServerError
	}

	return loanDetails, nil
}

func (l LoanServiceImplementation) GetAllLoansForCustomer(customerId string) ([]*responseDto.LoanDetails, error) {
	loanDetails, err := l.repo.GetAllLoansByCustomerId(customerId)
	if err != nil {
		log.Printf("failed to get loans for customer %s, error %v\n", customerId, err)
		return nil, app_errors.InternalServerError
	}
	return loanDetails, nil
}

func (l LoanServiceImplementation) ApproveLoan(loanApproveRequest *dto.LoanApproveRequest) error {
	loanId := loanApproveRequest.LoanId

	// validate loanId
	if loanId == "" {
		log.Println("loan id not specified")
		return invalidLoanId
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), TimeoutInSecond*time.Second)
	defer cancelFunc()

	txOption := &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	}

	tx, err := l.repo.CreateTransaction(ctx, txOption)
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

	loanDetails, err := l.repo.GetLoanById(loanId, tx)
	if err != nil {
		log.Println("loan can not be fetched")
		return loanNotPresent
	}

	if loanDetails.Status != responseDto.LoanStatusPending {
		log.Println("loan can not be approved, invalid status")
		err = fmt.Errorf("loan can not be approved, invalid status")
		return loanInvalidStatus
	}

	err = l.repo.UpdateLoanStatus(loanId, responseDto.LoanStatusApproved, tx)
	if err != nil {
		log.Printf("failed to approve loan for loanId %s, error %v\n", loanId, err)
		return app_errors.InternalServerError
	}
	return nil
}
