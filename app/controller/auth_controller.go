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

// LoginAsCustomer Login user as a Customer
// @Summary      Login user as a Customer
// @Description  Responds with the bearer token with customer role
// @Tags         Login
// @accept       json
// @Param        data body dto.LoginRequest true "username is mandatory"
// @Produce      json
// @Success      200 {object} dto.LoginResponse
// @Failure      400 {object} app_errors.ErrorResponse
// @Failure      500 {object} app_errors.ErrorResponse
// @Router       /auth/customer/login [post]
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

// LoginAsAdmin  Login user as an Admin
// @Summary      Login user as an Admin
// @Description  Responds with the bearer token with admin role
// @Tags         Login
// @accept       json
// @Param        data body dto.LoginRequest true "username is mandatory"
// @Produce      json
// @Success      200 {object} dto.LoginResponse
// @Failure      400 {object} app_errors.ErrorResponse
// @Failure      500 {object} app_errors.ErrorResponse
// @Router       /auth/admin/login [post]
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
