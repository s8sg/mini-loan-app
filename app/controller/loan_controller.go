package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	serverError "github.com/s8sg/mini-loan-app/app/app_errors"
	"github.com/s8sg/mini-loan-app/app/controller/dto"
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

// CreateLoanHandler Create a loan for a customer
// @Summary      Create a loan for a customer
// @Description  Create a loan for a customer, responds with the newly created loan details
// @Tags         Loans
// @accept       json
// @Param        token header  string true "Bearer customer-token"
// @Param        data body dto.LoanCreateRequest true "loan creation request"
// @Produce      json
// @Success      200 {object} dto.LoanDetails
// @Failure      400 {object} app_errors.ErrorResponse
// @Failure      500 {object} app_errors.ErrorResponse
// @Router       /user/loan [post]
func (h *LoanController) CreateLoanHandler(c *gin.Context) {
	loanCreateRequest := &dto.LoanCreateRequest{}
	err := c.BindJSON(loanCreateRequest)
	if err != nil {
		log.Printf("CreateLoanHandler: failed to parse request, error %v\n", err)
		serverError.RespondWithError(c, serverError.BadRequest)
		return
	}

	userIdContext, ok := c.Get("id")
	if !ok {
		log.Printf("CreateLoanHandler: user context not initialized\n")
		serverError.RespondWithError(c, serverError.BadRequest)
		return
	}

	customerId := fmt.Sprint(userIdContext)

	loanDetails, err := h.loanService.CreateLoan(customerId, loanCreateRequest)
	if err != nil {
		log.Printf("CreateLoanHandler: failed to create loan %v\n", err)
		serverError.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, loanDetails)
}

// GetLoansHandler Get all loans for a customer
// @Summary      Get all loans for a customer
// @Description  Responds with the all loan details belongs to customer
// @Tags         Loans
// @accept       json
// @Param        token header  string true "Bearer customer-token"
// @Produce      json
// @Success      200 {object} dto.GetAllLoansResponse
// @Failure      400 {object} app_errors.ErrorResponse
// @Failure      500 {object} app_errors.ErrorResponse
// @Router       /user/loans [get]
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func (h *LoanController) GetLoansHandler(c *gin.Context) {
	userIdContext, ok := c.Get("id")
	if !ok {
		log.Printf("GetLoansHandler: user context not initialized\n")
		serverError.RespondWithError(c, serverError.BadRequest)
		return
	}

	customerId := fmt.Sprint(userIdContext)

	loanDetails, err := h.loanService.GetAllLoansForCustomer(customerId)
	if err != nil {
		log.Printf("GetLoansHandler: failed to get loans %v\n", err)
		serverError.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.GetAllLoansResponse{Loans: loanDetails})
}

// ApproveLoanHandler Approve a loan
// @Summary      Approve a loan
// @Description  approve a loan
// @Tags         Loan Approval
// @accept       json
// @Param        token header  string true "Bearer admin-token"
// @Param        data body dto.LoanApproveRequest true "loan approval request"
// @Produce      json
// @Success      200 {object} dto.GenericSuccessResponse
// @Failure      400 {object} app_errors.ErrorResponse
// @Failure      404 {object} app_errors.ErrorResponse
// @Failure      500 {object} app_errors.ErrorResponse
// @Router       /loan/approve [post]
func (h *LoanController) ApproveLoanHandler(c *gin.Context) {
	loanApproveRequest := &dto.LoanApproveRequest{}
	err := c.BindJSON(loanApproveRequest)
	if err != nil {
		log.Printf("ApproveLoanHandler: failed to parse request, error %v\n", err)
		serverError.RespondWithError(c, serverError.BadRequest)
		return
	}

	err = h.loanService.ApproveLoan(loanApproveRequest)
	if err != nil {
		log.Printf("GetLoansHandler: failed to get loans %v\n", err)
		serverError.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, &dto.GenericSuccessResponse{Message: "successfully completed"})
}
