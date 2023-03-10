package dto

// LoginRequest login request
// @Description login request (Secret is optional)
type LoginRequest struct {
	Username string `json:"username" example:"user1"`
	Secret   string `json:"secret" example:"dummy-value"`
}

// LoginResponse login response body
// @Description login response with token
type LoginResponse struct {
	Token string `json:"token" example:"<bearer token>"`
}
