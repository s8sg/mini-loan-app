package app_errors

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// AppError custom error to wrap application code and message
type AppError struct {
	Code    int
	Message string
}

func (err *AppError) Error() string {
	return fmt.Sprintf("%s", err.Message)
}

type ErrorResponse struct {
	Error string
}

// Generic error object
var (
	BadRequest          = &AppError{Code: 400, Message: "bad request"}
	Unauthorised        = &AppError{Code: 401, Message: "unauthorised"}
	InternalServerError = &AppError{Code: 500, Message: "internal server error"}
)

func RespondWithError(c *gin.Context, err error) {

	appError, ok := err.(*AppError)
	if !ok {
		c.AbortWithStatusJSON(500, &ErrorResponse{Error: err.Error()})
	}
	c.AbortWithStatusJSON(appError.Code, &ErrorResponse{Error: appError.Message})
}
