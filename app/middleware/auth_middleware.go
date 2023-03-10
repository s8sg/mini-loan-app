package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/s8sg/mini-loan-app/app/app_errors"
	"github.com/s8sg/mini-loan-app/app/service"
	"log"
	"strings"
)

const (
	ROLE_KEY    = "role"
	USER_ID_KEY = "id"
)

func AuthMiddleware(service service.AuthService, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			log.Println("No token provided")
			app_errors.RespondWithError(c, app_errors.Unauthorised)
			return
		}

		splitToken := strings.Split(token, "Bearer ")
		if len(splitToken) != 2 {
			log.Println("Invalid token provided")
			app_errors.RespondWithError(c, app_errors.Unauthorised)
			return
		}

		authContext, err := service.ValidateToken(splitToken[1], role)
		if err != nil {
			log.Println(err)
			app_errors.RespondWithError(c, app_errors.Unauthorised)
			return
		}

		// set role in request context
		c.Set(ROLE_KEY, authContext.Role)
		// set userId in request context
		c.Set(USER_ID_KEY, authContext.UserId)

		c.Next()
	}
}
