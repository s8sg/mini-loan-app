package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/s8sg/mini-loan-app/app/controller/dto"
	serverError "github.com/s8sg/mini-loan-app/app/errors"
	"github.com/s8sg/mini-loan-app/app/service"
	"log"
	"net/http"
)

type LoanController struct {
	loanService service.LoanService
}

func InitLoanController(loanService service.LoanService) *LoanController {
	loanController := &LoanController{
		loanService: loanService,
	}
	return loanController
}

func (h *LoanController) CreateLoanHandler(c *gin.Context) {
	loanCreateRequest := &dto.LoanCreateRequest{}
	err := c.BindJSON(loanCreateRequest)
	if err != nil {
		log.Printf("CreateLoanHandler: failed to parse request, error %v\n", err)
		serverError.RespondWithGenericError(c, serverError.BadRequest)
		return
	}

	userIdContext, ok := c.Get("id")
	if !ok {
		log.Printf("CreateLoanHandler: user context not initialized\n")
		serverError.RespondWithGenericError(c, serverError.InternalServerError)
		return
	}

	customerId := fmt.Sprint(userIdContext)

	loanDetails, err := h.loanService.CreateLoan(customerId, loanCreateRequest)
	if err != nil {
		log.Printf("CreateLoanHandler: failed to create loan %v\n", err)
		serverError.RespondWithGenericError(c, serverError.InternalServerError)
		return
	}

	c.JSON(http.StatusOK, loanDetails)
}

func (h *LoanController) GetLoansHandler(c *gin.Context) {
	userIdContext, ok := c.Get("id")
	if !ok {
		log.Printf("GetLoansHandler: user context not initialized\n")
		serverError.RespondWithGenericError(c, serverError.InternalServerError)
		return
	}

	customerId := fmt.Sprint(userIdContext)

	loanDetails, err := h.loanService.GetAllLoansForCustomer(customerId)
	if err != nil {
		log.Printf("GetLoansHandler: failed to get loans %v\n", err)
		serverError.RespondWithGenericError(c, serverError.InternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"loans": loanDetails})
}

func (h *LoanController) ApproveLoanHandler(c *gin.Context) {
	loanApproveRequest := &dto.LoanApproveRequest{}
	err := c.BindJSON(loanApproveRequest)
	if err != nil {
		log.Printf("ApproveLoanHandler: failed to parse request, error %v\n", err)
		serverError.RespondWithGenericError(c, serverError.BadRequest)
		return
	}

	err = h.loanService.ApproveLoan(loanApproveRequest)
	if err != nil {
		log.Printf("GetLoansHandler: failed to get loans %v\n", err)
		serverError.RespondWithGenericError(c, serverError.InternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully approved loan"})
}
