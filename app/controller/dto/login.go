package dto

type LoginRequest struct {
	Username string `json:"username"`
	Secret   string `json:"secret"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
