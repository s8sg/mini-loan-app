# Mini Loan App
This is a implementation for mini loan app

## Requirements
* This project needs [docker](https://www.docker.com) to run the app
* To build the project outside of docker container `go1.20.2` is required.


## Build
This projects builds the app within docker container (recommended)
```bash
docker-compose build
```

#### Build locally

In order to build it outside Go `go1.20.2` is required.  
Below commands build the app
```bash
cd app
./build.sh
```

## Run 
To run the whole stack in docker
```bash
docker-compose up -d
```

#### Run locally

In order to run outside of the docker container
First run postgres
```bash
docker-compose up -d postgres
```
The run the app locally (this will require you to build locally) 
```bash
./app
```

## Swagger
Once you run the stack you should be able to see the swagger at 
[http://localhost:8081/docs/index.html](http://localhost:8081/docs/index.html)

## Generate Swagger
Swagger is generated with `github.com/swaggo/swag/cmd/swag`  
To install
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```
In order to re-generate swagger file with latest changes run
```bash
swag init
```

## Run Unit Test
Due to limited time I opted out for writing unit test for the core business module.  
Tests are written for `loan_service`, `auth_service` and `repayment_service`.  


Similar approach can be used for testing other components

## Run Integration Test
The integration test tests the primary business logic
This runs the server in port 8080 and execute the happy flow
* User Login 
* User creates a loan 
* User lists all loan
* Admin Login
* Admin approves the loan created by user 
* User repay all repayments
* User list all loan (PAID)
```bash
cd app
./integration_test.sh
```

## DB
We are using **postgres** as DB.  
The schema are present at `db/schema/schema.sql`

## Design Choice
The project has the below modules
```
controller :  Prvides all http handler endpoint 
service:      Implements the core business logic
repository:   Provide the storage functionality
server:       Initate the route for controller
middleware:   Provides middleware for server to use 
```

This project is written keeping **SOLID** principle in mind.  
All layers are abstracted with `interfaces`.  
This is also allows to mock dependencies (written as interface) to test specific module.

The dependency of project module can be clearly seen at `main.go`
```
1. init db connection                                  | DB
2. init repository with db                             | Repository(DB)
3. init service (with repository for which its needed) | Service(Repository(DB))
4. init controller with service                        | Controller(Service(Repository(DB)))
5. create server                                       | Server
6. configure server route configuration                | Server(Controller, Middleware)
7. start server
```

This project uses [Go-GIN](https://github.com/gin-gonic/gin) as the web framework
Frame-work specific implementation can only be seen at 
```
server
controller
middleware
```
If needed in future iteration the framework dependencies can also be abstracted 
