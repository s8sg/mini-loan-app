package server

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/s8sg/mini-loan-app/app/controller"
	_ "github.com/s8sg/mini-loan-app/app/docs"
	"github.com/s8sg/mini-loan-app/app/middleware"
	"github.com/s8sg/mini-loan-app/app/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server : implements Server and controller.RestServer
type Server struct {
	router *gin.Engine
	port   string
}

func GetServer(port string) *Server {
	return &Server{
		router: gin.Default(),
		port:   port,
	}
}

// InitRoute : takes a list of controller and initialize the routes for the server
func (server *Server) InitRoute(
	authService service.AuthService,
	loanController *controller.LoanController,
	authController *controller.AuthController,
	repaymentController *controller.RepaymentController) {

	router := server.router
	// Host swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// all /v1/user is authenticated and authorized for customer
	userRoute := router.Group("/api/v1/user",
		middleware.AuthMiddleware(authService, service.USER_TYPE_CUSTOMER))

	userRoute.POST("/loan", loanController.CreateLoanHandler)
	userRoute.GET("/loans", loanController.GetLoansHandler)
	userRoute.POST("/loan/repayment", repaymentController.RepayLoanHandler)

	// all /v1/admin is authenticated and authorized for admin
	adminRoute := router.Group("/api/v1/admin",
		middleware.AuthMiddleware(authService, service.USER_TYPE_ADMIN))

	adminRoute.POST("/loan/approve", loanController.ApproveLoanHandler)

	// all /v1/auth is open
	authRoute := router.Group("/api/v1/auth")

	authRoute.POST("/customer/login", authController.LoginAsCustomer)
	authRoute.POST("/admin/login", authController.LoginAsAdmin)
}

// Start : starts the server on the provided port (listen to signals)
func (server *Server) Start() error {
	err := endless.ListenAndServe(":"+server.port, server.router)
	return err
}
