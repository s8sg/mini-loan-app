package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/s8sg/mini-loan-app/app/app_errors"
	"github.com/s8sg/mini-loan-app/app/controller/dto"
	"github.com/s8sg/mini-loan-app/app/service"
	"log"
	"net/http"
)

type AuthController struct {
	authService service.AuthService
}

func InitAuthController(authService service.AuthService) *AuthController {
	loginController := &AuthController{
		authService: authService,
	}
	return loginController
}

func (h *AuthController) LoginAsCustomer(c *gin.Context) {
	loginRequest := &dto.LoginRequest{}
	err := c.BindJSON(loginRequest)
	if err != nil {
		log.Printf("LoginAsCustomer: failed to parse request, error: %v\n", err)
		app_errors.RespondWithError(c, app_errors.BadRequest)
		return
	}
	token, err := h.authService.Login(loginRequest.Username, service.USER_TYPE_CUSTOMER, loginRequest.Secret)
	if err != nil {
		log.Printf("LoginAsCustomer: failed to login, error: %v\n", err)
		app_errors.RespondWithError(c, err)
		return
	}
	loginResponse := &dto.LoginResponse{
		Token: token,
	}
	c.JSON(http.StatusOK, loginResponse)
}

func (h *AuthController) LoginAsAdmin(c *gin.Context) {
	loginRequest := &dto.LoginRequest{}
	err := c.BindJSON(loginRequest)
	if err != nil {
		log.Printf("LoginAsAdmin: failed to parse request, error: %v\n", err)
		app_errors.RespondWithError(c, app_errors.BadRequest)
		return
	}
	token, err := h.authService.Login(loginRequest.Username, service.USER_TYPE_ADMIN, loginRequest.Secret)
	if err != nil {
		log.Printf("LoginAsAdmin: failed to login, error: %v\n", err)
		app_errors.RespondWithError(c, err)
		return
	}

	loginResponse := &dto.LoginResponse{
		Token: token,
	}
	c.JSON(http.StatusOK, loginResponse)
}
