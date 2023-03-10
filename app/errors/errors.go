package errors

import (
	"github.com/gin-gonic/gin"
)

type Error struct {
	code    int
	message string
}

/*
type (err *Error) error() {

}*/

var (
	BadRequest          = Error{code: 400, message: "bad request"}
	Unauthorised        = Error{code: 401, message: "unauthorised"}
	InternalServerError = Error{code: 500, message: "internal server error"}
)

func RespondWithGenericError(c *gin.Context, err Error) {
	c.AbortWithStatusJSON(err.code, gin.H{"error": err.message})
}

func RespondWithError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}
