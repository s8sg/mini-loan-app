package service

import (
	"context"
	"database/sql"
	"fmt"
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
		return nil, fmt.Errorf("failed to create loan: %v", err)
	}

	return loanDetails, nil
}

func (l LoanServiceImplementation) GetAllLoansForCustomer(customerId string) ([]*responseDto.LoanDetails, error) {
	loanDetails, err := l.repo.GetAllLoansByCustomerId(customerId)
	if err != nil {
		log.Printf("failed to get loans for customer %s, error %v\n", customerId, err)
		return nil, fmt.Errorf("failed to get loans, error: %v", err)
	}
	return loanDetails, nil
}

func (l LoanServiceImplementation) ApproveLoan(loanApproveRequest *dto.LoanApproveRequest) error {
	loanId := loanApproveRequest.LoanId

	ctx, cancelFunc := context.WithTimeout(context.Background(), TimeoutInSecond*time.Second)
	defer cancelFunc()

	txOption := &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	}

	tx, err := l.repo.CreateTransaction(ctx, txOption)
	if err != nil {
		log.Println("failed to initiate transaction")
		return err
	}

	defer func() {
		if err != nil {
			log.Println("calling rollback for error " + err.Error())
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	loanDetails, err := l.repo.GetLoanById(loanId, tx)
	if err != nil {
		log.Println("loan can not be fetched")
		return err
	}

	if loanDetails.Status != responseDto.LoanStatusPending {
		log.Println("loan can not be approved, invalid status")
		err = fmt.Errorf("loan can not be approved, invalid status")
		return err
	}

	err = l.repo.UpdateLoanStatus(loanId, responseDto.LoanStatusApproved, tx)
	if err != nil {
		log.Printf("failed to approve loan for loanId %s, error %v\n", loanId, err)
		return fmt.Errorf("failed to approve loan, error: %v", err)
	}
	return nil
}
