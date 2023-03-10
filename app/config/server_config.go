package config

import (
	"fmt"
	"github.com/s8sg/mini-loan-app/app/controller"
	repository "github.com/s8sg/mini-loan-app/app/repostory"
	"github.com/s8sg/mini-loan-app/app/server"
	"github.com/s8sg/mini-loan-app/app/service"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	Port        = "8085"
	DbUser      = "root"
	DbPassword  = "aspire123"
	DbHost      = "localhost"
	DbName      = "mini_loan_app"
	AuthHmacKey = "secretkey"
)

func InitializeServer() (*server.Server, error) {
	// init variable from env
	initializeConfigFromEnv()

	// init db connection
	db, err := InitialiseDB(DbConfig{
		User:     DbUser,
		Password: DbPassword,
		Host:     DbHost,
		DBName:   DbName,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot initialize db, err: %v", err)
	}

	// init repository with db
	loanRepository := repository.GetLoanRepository(db)

	authService := service.GetAuthService(AuthHmacKey)
	// init service with repository
	loanService := service.GetLoanService(loanRepository)
	repaymentService := service.GetRepaymentService(loanRepository)

	// init controllers with service
	authController := controller.InitAuthController(authService)
	loanController := controller.InitLoanController(loanService)
	repaymentController := controller.InitRepaymentController(repaymentService)

	// create server and configure with controller specific route configuration
	appServer := server.GetServer(Port)
	// Initialize routes
	appServer.InitRoute(authService, loanController, authController, repaymentController)

	return appServer, nil
}

func initializeConfigFromEnv() {
	env := os.Getenv("SERVER_PORT")
	if env != "" {
		log.Println("SERVER_PORT: ", env)
		Port = env
	}
	env = os.Getenv("DB_USER")
	if env != "" {
		log.Println("DB_USER: ", env)
		DbUser = env
	}
	env = os.Getenv("DB_PASSWORD")
	if env != "" {
		log.Println("DB_PASSWORD: ", env)
		DbPassword = env
	}
	env = os.Getenv("DB_HOST")
	if env != "" {
		log.Println("DB_HOST: ", env)
		DbHost = env
	}
	env = os.Getenv("DB_NAME")
	if env != "" {
		log.Println("DB_NAME: ", env)
		DbName = env
	}
	env = os.Getenv("AUTH_HMAC_SIGNING_KEY")
	if env != "" {
		log.Println("AUTH_HMAC_SIGNING_KEY: ", env)
		AuthHmacKey = env
	}
}
