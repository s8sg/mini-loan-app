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

type RepaymentController struct {
	repaymentService service.RepaymentService
}

func InitRepaymentController(repaymentService service.RepaymentService) *RepaymentController {
	repaymentController := &RepaymentController{
		repaymentService: repaymentService,
	}
	return repaymentController
}

func (h *RepaymentController) RepayLoanHandler(c *gin.Context) {
	loanRepaymentRequest := &dto.LoanRepaymentRequest{}
	err := c.BindJSON(loanRepaymentRequest)
	if err != nil {
		log.Printf("RepayLoanHandler: failed to parse request, error %v\n", err)
		serverError.RespondWithGenericError(c, serverError.BadRequest)
	}

	userIdContext, ok := c.Get("id")
	if !ok {
		log.Printf("GetLoansHandler: user context not initialized\n")
		serverError.RespondWithGenericError(c, serverError.InternalServerError)
	}

	customerId := fmt.Sprint(userIdContext)

	err = h.repaymentService.Repay(customerId, loanRepaymentRequest)
	if err != nil {
		log.Printf("GetLoansHandler: failed to get loans %v\n", err)
		serverError.RespondWithGenericError(c, serverError.InternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"message": "loan repayment successful"})
}
