package main

import (
	"github.com/s8sg/mini-loan-app/app/config"
	"github.com/s8sg/mini-loan-app/app/controller"
	repository "github.com/s8sg/mini-loan-app/app/repostory"
	"github.com/s8sg/mini-loan-app/app/server"
	"github.com/s8sg/mini-loan-app/app/service"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	port        = "8081"
	dbUser      = "root"
	dbPassword  = "aspire123"
	dbHost      = "localhost"
	dbName      = "mini_loan_app"
	authHmacKey = "secretkey"
)

func main() {
	// init db connection
	db, err := config.InitialiseDB(config.DbConfig{
		User:     dbUser,
		Password: dbPassword,
		Host:     dbHost,
		DBName:   dbName,
	})
	if err != nil {
		log.Println("cannot initialize db", err)
		os.Exit(1)
	}

	// init repository with db
	loanRepository := repository.GetLoanRepository(db)

	authService := service.GetAuthService(authHmacKey)
	// init service with repository
	loanService := service.GetLoanService(loanRepository)
	repaymentService := service.GetRepaymentService(loanRepository)

	// init controllers with service
	authController := controller.InitAuthController(authService)
	loanController := controller.InitLoanController(loanService)
	repaymentController := controller.InitRepaymentController(repaymentService)

	// create server and configure with controller specific route configuration
	server := server.GetServer(port)
	// Initialize routes
	server.InitRoute(authService, loanController, authController, repaymentController)

	log.Println("Starting the server on " + port)
	err = server.Start()
	if err != nil {
		log.Println("stopped listen on " + port)
		os.Exit(1)
	}
}
