package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/s8sg/mini-loan-app/app/config"
	"github.com/s8sg/mini-loan-app/app/dto"
	"github.com/s8sg/mini-loan-app/app/service"
	"github.com/s8sg/mini-loan-app/app/util"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

const (
	LoanAmount1 = 100000
	Term1       = 2

	LoanAmount2 = 10000
	Term2       = 3
)

var (
	ValidUser1 = "user1-" + uuid.New().String()
	ValidUser2 = "user2-" + uuid.New().String()
	ValidAdmin = "admin-" + uuid.New().String()

	CustomerToken1        = ""
	CustomerToken2        = ""
	AdminToken            = ""
	User1LoanId           = ""
	User2LoanId           = ""
	User1LoanRepaymentIds = []string{}
	User2LoanRepaymentIds = []string{}
)

func Init() {
	server, err := config.InitializeServer()
	if err != nil {
		log.Fatalf("Failed to initialize server, error: %v", err)
	}
	go func() {
		log.Println("Starting the server")
		err = server.Start()
		if err != nil {
			log.Println("stopped listening")
			os.Exit(1)
		}
	}()
}

func callAPI(t *testing.T, method, url string, body []byte, token string) (int, []byte) {
	t.Helper()

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err := res.Body.Close(); err != nil {
		t.Fatal(err)
	}

	return res.StatusCode, resBody
}

// TestMainAPI expects that postgres and the service is already running
func TestMainAPI(t *testing.T) {
	Init()

	t.Run("Customer Login", func(t *testing.T) {
		// invalid request with empty body
		t.Run("POST /api/v1/auth/customer/login 400", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{}`))
			status, _ := callAPI(t, "POST", "http://localhost:8085/api/v1/auth/customer/login", body, "")
			if status != 400 {
				t.Errorf("expected status 400 but got %d", status)
			}
		})

		// valid request for user1
		t.Run("POST /api/v1/auth/customer/login 200", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{"username":"%s","secret":"secret"}`, ValidUser1))
			status, body := callAPI(t, "POST", "http://localhost:8085/api/v1/auth/customer/login", body, "")
			if status != 200 {
				t.Errorf("expected status 200 but got %d %v", status, string(body))
			}

			tokenStruct := struct {
				Token string `json:"token"`
			}{}
			if err := json.Unmarshal(body, &tokenStruct); err != nil {
				t.Fatal(err)
			}
			if tokenStruct.Token == "" {
				t.Fatal(`response body doesn't contain "token" field'`)
			}

			CustomerToken1 = tokenStruct.Token
		})

		// valid request for user2
		t.Run("POST /api/v1/auth/customer/login 200", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{"username":"%s","secret":"secret"}`, ValidUser2))
			status, body := callAPI(t, "POST", "http://localhost:8085/api/v1/auth/customer/login", body, "")
			if status != 200 {
				t.Errorf("expected status 200 but got %d %v", status, string(body))
			}

			tokenStruct := struct {
				Token string `json:"token"`
			}{}
			if err := json.Unmarshal(body, &tokenStruct); err != nil {
				t.Fatal(err)
			}
			if tokenStruct.Token == "" {
				t.Fatal(`response body doesn't contain "token" field'`)
			}

			CustomerToken2 = tokenStruct.Token
		})
	})

	t.Run("Admin Login", func(t *testing.T) {
		// invalid request for admin login
		t.Run("POST /api/v1/auth/admin/login 400", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{}`))
			status, _ := callAPI(t, "POST", "http://localhost:8085/api/v1/auth/admin/login", body, "")
			if status != 400 {
				t.Errorf("expected status 400 but got %d", status)
			}
		})

		// valid request for admin1
		t.Run("POST /api/v1/auth/admin/login 200", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{"username":"%s","secret":"secret"}`, ValidAdmin))
			status, body := callAPI(t, "POST", "http://localhost:8085/api/v1/auth/admin/login", body, "")
			if status != 200 {
				t.Errorf("expected status 200 but got %d %v", status, string(body))
			}

			tokenStruct := struct {
				Token string `json:"token"`
			}{}
			if err := json.Unmarshal(body, &tokenStruct); err != nil {
				t.Fatal(err)
			}
			if tokenStruct.Token == "" {
				t.Fatal(`response body doesn't contain "token" field'`)
			}

			AdminToken = tokenStruct.Token
		})
	})

	t.Run("Loan Create", func(t *testing.T) {
		// request with no token set
		t.Run("POST /api/v1/user/loan 401", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{}`))
			status, _ := callAPI(t, "POST", "http://localhost:8085/api/v1/user/loan", body, "")
			if status != 401 {
				t.Errorf("expected status 401 but got %d", status)
			}
		})

		// request with admin token set
		t.Run("POST /api/v1/user/loan 401", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{}`))
			status, _ := callAPI(t, "POST", "http://localhost:8085/api/v1/user/loan", body, AdminToken)
			if status != 401 {
				t.Errorf("expected status 401 but got %d", status)
			}
		})

		// request with user token set, but invalid body
		t.Run("POST /api/v1/user/loan 401", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{}`))
			status, _ := callAPI(t, "POST", "http://localhost:8085/api/v1/user/loan", body, CustomerToken1)
			if status != 400 {
				t.Errorf("expected status 400 but got %d", status)
			}
		})

		// request with user token set, but invalid amount 0
		t.Run("POST /api/v1/user/loan 401", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{"amount": %d, "term": 2}`, 0))
			status, body := callAPI(t, "POST", "http://localhost:8085/api/v1/user/loan", body, CustomerToken1)
			if status != 400 {
				t.Errorf("expected status 400 but got %d, %v", status, string(body))
			}
		})

		// request with user token set, but invalid amount 0
		t.Run("POST /api/v1/user/loan 401", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{"amount": %d, "term": 2}`, 0))
			status, body := callAPI(t, "POST", "http://localhost:8085/api/v1/user/loan", body, CustomerToken1)
			if status != 400 {
				t.Errorf("expected status 400 but got %d, %v", status, string(body))
			}
		})

		// request with user token set, and valid amount 10000 for customer 1
		t.Run("POST /api/v1/user/loan 200", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{"amount": %d, "term": %d}`, LoanAmount1, Term1))
			status, body := callAPI(t, "POST", "http://localhost:8085/api/v1/user/loan", body, CustomerToken1)
			if status != 201 {
				t.Errorf("expected status 201 but got %d, %v", status, string(body))
			}

			loanDetails := &dto.LoanDetails{}
			if err := json.Unmarshal(body, &loanDetails); err != nil {
				t.Fatal(err)
			}

			if loanDetails.LoanId == "" {
				t.Errorf("loan id must not be blank, %v", string(body))
			}

			User1LoanId = loanDetails.LoanId

			if loanDetails.CustomerId != ValidUser1 {
				t.Errorf("loan created for wrong customer, %v", string(body))
			}

			if loanDetails.Status != "PENDING" {
				t.Errorf("loan created with wrong status, %v", string(body))
			}

			if loanDetails.Term != Term1 {
				t.Errorf("loan created with wrong terms, %v", string(body))
			}

			if !loanDetails.TotalAmount.Equal(decimal.NewFromInt(LoanAmount1)) {
				t.Errorf("loan created with wrong amount, %v", string(body))
			}

			year, month, day := loanDetails.StartDate.Date()
			expectedYear, expectedMonth, expectedDay := util.GetCurrentTimeInUtc().Date()

			if year != expectedYear || month != expectedMonth || day != expectedDay {
				t.Errorf("loan created with wrong date, %v", string(body))
			}

			if len(loanDetails.Repayments) != Term1 {
				t.Errorf("loan created with invalid no of repayments, %v", string(body))
			}

			for _, repayment := range loanDetails.Repayments {
				User1LoanRepaymentIds = append(User1LoanRepaymentIds, repayment.RepaymentId)
			}

			repayment1 := loanDetails.Repayments[0]

			if repayment1.RepaymentId == "" {
				t.Errorf("rerpayment id must not be blank, %v", string(body))
			}

			if repayment1.Status != "PENDING" {
				t.Errorf("rerpayment created with wrong status, %v", string(body))
			}

			if !repayment1.Amount.Equal(decimal.NewFromInt(LoanAmount1).Div(decimal.NewFromInt(Term1))) {
				t.Errorf("rerpayment created with wrong status, %v", string(body))
			}

			if repayment1.Number != 1 {
				t.Errorf("rerpayment created with wrong no, %v", string(body))
			}

			if repayment1.DueDate.Compare(loanDetails.StartDate.Add(service.RepaymentFrequency)) != 0 {
				t.Errorf("rerpayment created with wrong no, %v", string(body))
			}

			repayment2 := loanDetails.Repayments[1]

			if repayment2.RepaymentId == "" {
				t.Errorf("rerpayment id must not be blank, %v", string(body))
			}

			if repayment2.Status != "PENDING" {
				t.Errorf("rerpayment created with wrong status, %v", string(body))
			}

			if !repayment2.Amount.Equal(decimal.NewFromInt(LoanAmount1).Div(decimal.NewFromInt(Term1))) {
				t.Errorf("rerpayment created with wrong status, %v", string(body))
			}

			if repayment2.Number != 2 {
				t.Errorf("rerpayment created with wrong no, %v", string(body))
			}

			if repayment2.DueDate.Compare(loanDetails.StartDate.Add(service.RepaymentFrequency*2)) != 0 {
				t.Errorf("rerpayment created with wrong no, %v", string(body))
			}
		})

		// request with user token set, and valid amount for customer 2
		t.Run("POST /api/v1/user/loan 200", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{"amount": %d, "term": %d}`, LoanAmount2, Term2))
			status, body := callAPI(t, "POST", "http://localhost:8085/api/v1/user/loan", body, CustomerToken2)
			if status != 201 {
				t.Errorf("expected status 201 but got %d, %v", status, string(body))
			}

			loanDetails := &dto.LoanDetails{}
			if err := json.Unmarshal(body, &loanDetails); err != nil {
				t.Fatal(err)
			}

			if loanDetails.LoanId == "" {
				t.Errorf("loan id must not be blank, %v", string(body))
			}

			User2LoanId = loanDetails.LoanId

			if loanDetails.CustomerId != ValidUser2 {
				t.Errorf("loan created for wrong customer, %v", string(body))
			}

			if loanDetails.Status != "PENDING" {
				t.Errorf("loan created with wrong status, %v", string(body))
			}

			if loanDetails.Term != Term2 {
				t.Errorf("loan created with wrong terms, %v", string(body))
			}

			if !loanDetails.TotalAmount.Equal(decimal.NewFromInt(LoanAmount2)) {
				t.Errorf("loan created with wrong amount, %v", string(body))
			}

			year, month, day := loanDetails.StartDate.Date()
			expectedYear, expectedMonth, expectedDay := util.GetCurrentTimeInUtc().Date()

			if year != expectedYear || month != expectedMonth || day != expectedDay {
				t.Errorf("loan created with wrong date, %v", string(body))
			}

			if len(loanDetails.Repayments) != Term2 {
				t.Errorf("loan created with invalid no of repayments, %v", string(body))
			}

			for _, repayment := range loanDetails.Repayments {
				User2LoanRepaymentIds = append(User2LoanRepaymentIds, repayment.RepaymentId)
			}

			repayment1 := loanDetails.Repayments[0]

			if repayment1.RepaymentId == "" {
				t.Errorf("rerpayment id must not be blank, %v", string(body))
			}

			if repayment1.Status != "PENDING" {
				t.Errorf("rerpayment created with wrong status, %v", string(body))
			}

			if !repayment1.Amount.Equal(decimal.NewFromInt(LoanAmount2).Div(decimal.NewFromInt(Term2))) {
				t.Errorf("rerpayment created with wrong status, %v", string(body))
			}

			if repayment1.Number != 1 {
				t.Errorf("rerpayment created with wrong no, %v", string(body))
			}

			if repayment1.DueDate.Compare(loanDetails.StartDate.Add(service.RepaymentFrequency)) != 0 {
				t.Errorf("rerpayment created with wrong no, %v", string(body))
			}

			repayment2 := loanDetails.Repayments[1]

			if repayment2.RepaymentId == "" {
				t.Errorf("rerpayment id must not be blank, %v", string(body))
			}

			if repayment2.Status != "PENDING" {
				t.Errorf("rerpayment created with wrong status, %v", string(body))
			}

			if !repayment2.Amount.Equal(decimal.NewFromInt(LoanAmount2).Div(decimal.NewFromInt(Term2))) {
				t.Errorf("rerpayment created with wrong status, %v", string(body))
			}

			if repayment2.Number != 2 {
				t.Errorf("rerpayment created with wrong no, %v", string(body))
			}

			if repayment2.DueDate.Compare(loanDetails.StartDate.Add(service.RepaymentFrequency*2)) != 0 {
				t.Errorf("rerpayment created with wrong no, %v", string(body))
			}
		})
	})

	t.Run("Approve Loan", func(t *testing.T) {
		// request with no token set
		t.Run("POST /api/v1/admin/loan/approve 401", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{}`))
			status, _ := callAPI(t, "POST", "http://localhost:8085/api/v1/admin/loan/approve", body, "")
			if status != 401 {
				t.Errorf("expected status 401 but got %d", status)
			}
		})

		// request with customer token set
		t.Run("POST /api/v1/admin/loan/approve 401", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{}`))
			status, _ := callAPI(t, "POST", "http://localhost:8085/api/v1/admin/loan/approve", body, CustomerToken1)
			if status != 401 {
				t.Errorf("expected status 401 but got %d", status)
			}
		})

		// request with admin token set but invalid body
		t.Run("POST /api/v1/admin/loan/approve 400", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{}`))
			status, _ := callAPI(t, "POST", "http://localhost:8085/api/v1/admin/loan/approve", body, AdminToken)
			if status != 400 {
				t.Errorf("expected status 400 but got %d", status)
			}
		})

		// request with admin token set but invalid loan id
		t.Run("POST /api/v1/admin/loan/approve 400", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{"loan-id": "invalid"}`))
			status, _ := callAPI(t, "POST", "http://localhost:8085/api/v1/admin/loan/approve", body, AdminToken)
			if status != 404 {
				t.Errorf("expected status 400 but got %d", status)
			}
		})

		// request with admin token set and valid body
		t.Run("POST /api/v1/admin/loan/approve 200", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{"loan-id": "%s"}`, User1LoanId))
			status, body := callAPI(t, "POST", "http://localhost:8085/api/v1/admin/loan/approve", body, AdminToken)
			if status != 200 {
				t.Errorf("expected status 400 but got %d", status)
			}
			response := struct {
				Message string `json:"message"`
			}{}
			if err := json.Unmarshal(body, &response); err != nil {
				t.Fatal(err)
			}

			if response.Message != "successfully completed" {
				t.Errorf("expected 'successfully completed' but got %s", response.Message)
			}
		})
	})

	t.Run("Repay Loan", func(t *testing.T) {
		// request with no token set
		t.Run("POST /api/v1/user/loan/repayment 401", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{}`))
			status, _ := callAPI(t, "POST", "http://localhost:8085/api/v1/user/loan/repayment", body, "")
			if status != 401 {
				t.Errorf("expected status 401 but got %d", status)
			}
		})

		// request with invalid repayment id
		t.Run("POST /api/v1/user/loan/repayment 404", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{"repayment-id": "invalidRepaymentId", "amount": 10000}`))
			status, _ := callAPI(t, "POST", "http://localhost:8085/api/v1/user/loan/repayment", body, CustomerToken1)
			if status != 404 {
				t.Errorf("expected status 404 but got %d", status)
			}
		})

		// request with other customers repayment id
		t.Run("POST /api/v1/user/loan/repayment 404", func(t *testing.T) {
			body := []byte(fmt.Sprintf(`{"repayment-id": "%s", "amount": 10000}`, User2LoanRepaymentIds[0]))
			status, _ := callAPI(t, "POST", "http://localhost:8085/api/v1/user/loan/repayment", body, CustomerToken1)
			if status != 404 {
				t.Errorf("expected status 404 but got %d", status)
			}
		})

		// repay all repayments for customer 1
		t.Run("POST /api/v1/user/loan/repayment 200", func(t *testing.T) {
			for _, repaymentId := range User1LoanRepaymentIds {
				body := []byte(fmt.Sprintf(`{"repayment-id": "%s", "amount": %d}`, repaymentId, LoanAmount1/Term1))
				status, body := callAPI(t, "POST", "http://localhost:8085/api/v1/user/loan/repayment", body, CustomerToken1)
				if status != 200 {
					t.Errorf("expected status 404 but got %d", status)
				}
				response := struct {
					Message string `json:"message"`
				}{}
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatal(err)
				}

				if response.Message != "successfully completed" {
					t.Errorf("expected 'successfully completed' but got %s", response.Message)
				}
			}
		})
	})

	t.Run("Get Loans", func(t *testing.T) {
		// request with no token set
		t.Run("GET /api/v1/user/loans 401", func(t *testing.T) {
			status, _ := callAPI(t, "GET", "http://localhost:8085/api/v1/user/loans", nil, "")
			if status != 401 {
				t.Errorf("expected status 401 but got %d", status)
			}
		})

		// request with admin token set
		t.Run("GET /api/v1/user/loans 401", func(t *testing.T) {
			status, _ := callAPI(t, "GET", "http://localhost:8085/api/v1/user/loans", nil, AdminToken)
			if status != 401 {
				t.Errorf("expected status 401 but got %d", status)
			}
		})

		// request with customer1 token set (validate paid loan and repayments)
		t.Run("GET /api/v1/user/loans 401", func(t *testing.T) {
			status, body := callAPI(t, "GET", "http://localhost:8085/api/v1/user/loans", nil, CustomerToken1)
			if status != 200 {
				t.Errorf("expected status 200 but got %d", status)
			}
			response := struct {
				Loans []*dto.LoanDetails `json:"loans"`
			}{}
			if err := json.Unmarshal(body, &response); err != nil {
				t.Fatal(err)
			}

			if len(response.Loans) != 1 {
				t.Errorf("expected loan count 1 but got %d", len(response.Loans))
			}

			loanDetails := response.Loans[0]

			if loanDetails.LoanId != User1LoanId {
				t.Errorf("loan id is invalid, %v", string(body))
			}

			if loanDetails.CustomerId != ValidUser1 {
				t.Errorf("loan created for wrong customer, %v", string(body))
			}

			if loanDetails.Status != "PAID" {
				t.Errorf("loan isn't paid, %v", string(body))
			}

			for _, repayments := range loanDetails.Repayments {
				if repayments.Status != "PAID" {
					t.Errorf("repayment isn't paid, %v", string(body))
				}
			}
		})

		// request with customer2 token set (validate pending loan and repayments)
		t.Run("GET /api/v1/user/loans 401", func(t *testing.T) {
			status, body := callAPI(t, "GET", "http://localhost:8085/api/v1/user/loans", nil, CustomerToken2)
			if status != 200 {
				t.Errorf("expected status 200 but got %d", status)
			}
			response := struct {
				Loans []*dto.LoanDetails `json:"loans"`
			}{}
			if err := json.Unmarshal(body, &response); err != nil {
				t.Fatal(err)
			}

			if len(response.Loans) != 1 {
				t.Errorf("expected loan count 1 but got %d", len(response.Loans))
			}

			loanDetails := response.Loans[0]

			if loanDetails.LoanId != User2LoanId {
				t.Errorf("loan id is invalid, %v", string(body))
			}

			if loanDetails.CustomerId != ValidUser2 {
				t.Errorf("loan created for wrong customer, %v", string(body))
			}

			if loanDetails.Status != "PENDING" {
				t.Errorf("loan isn't paid, %v", string(body))
			}

			for _, repayments := range loanDetails.Repayments {
				if repayments.Status != "PENDING" {
					t.Errorf("repayment isn't paid, %v", string(body))
				}
			}
		})
	})
}
