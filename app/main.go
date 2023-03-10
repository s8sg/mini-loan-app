package main

import (
	"github.com/s8sg/mini-loan-app/app/config"
	"log"
	"os"
)

// @title           Mini Loan APP
// @version         1.0
// @description     A loan management service API in Go using Gin framework
// @host      localhost:8085
// @BasePath  /api/v1
func main() {
	server, err := config.InitializeServer()
	if err != nil {
		log.Fatalf("Failed to initialize server, error: %v", err)
	}
	log.Println("Starting the server")
	err = server.Start()
	if err != nil {
		log.Println("stopped listening")
		os.Exit(1)
	}
}
