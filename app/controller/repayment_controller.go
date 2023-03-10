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

type RepaymentController struct {
	repaymentService service.RepaymentService
}

func InitRepaymentController(repaymentService service.RepaymentService) *RepaymentController {
	repaymentController := &RepaymentController{
		repaymentService: repaymentService,
	}
	return repaymentController
}

// RepayLoanHandler Repay a repayment of loan
// @Summary      Repay a repayment of loan
// @Description  repay a repayment, mark loan as paid when all repayment paid
// @Tags         Loans
// @accept       json
// @Param        token header  string true "Bearer admin-token"
// @Param        data body dto.LoanRepaymentRequest true "loan repayment request"
// @Produce      json
// @Success      200 {object} dto.GenericSuccessResponse
// @Failure      400 {object} app_errors.ErrorResponse
// @Failure      404 {object} app_errors.ErrorResponse
// @Failure      500 {object} app_errors.ErrorResponse
// @Router       /user/loan/repayment [post]
func (h *RepaymentController) RepayLoanHandler(c *gin.Context) {
	loanRepaymentRequest := &dto.LoanRepaymentRequest{}
	err := c.BindJSON(loanRepaymentRequest)
	if err != nil {
		log.Printf("RepayLoanHandler: failed to parse request, error %v\n", err)
		serverError.RespondWithError(c, serverError.BadRequest)
		return
	}

	userIdContext, ok := c.Get("id")
	if !ok {
		log.Printf("GetLoansHandler: user context not initialized\n")
		serverError.RespondWithError(c, serverError.BadRequest)
		return
	}

	customerId := fmt.Sprint(userIdContext)

	err = h.repaymentService.Repay(customerId, loanRepaymentRequest)
	if err != nil {
		log.Printf("GetLoansHandler: failed to get loans %v\n", err)
		serverError.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, &dto.GenericSuccessResponse{Message: "successfully completed"})
}
