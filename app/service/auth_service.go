package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/s8sg/mini-loan-app/app/app_errors"
	"log"
	"time"
)

const (
	USER_TYPE_CUSTOMER = "customer"
	USER_TYPE_ADMIN    = "admin"
)

var (
	InvalidToken = fmt.Errorf("token is not valid")
)

var (
	userIdMustBeProvided = &app_errors.AppError{Code: 400, Message: "username must be provided"}
)

type AuthContext struct {
	UserId string
	Role   string
}

type AuthService interface {
	Login(userid string, userType string, secret string) (string, error)
	ValidateToken(token string, role string) (*AuthContext, error)
}

type AuthServiceImplementation struct {
	secretKey string
}

type Claims struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

func GetAuthService(secretKey string) AuthService {
	return &AuthServiceImplementation{
		secretKey: secretKey,
	}
}

func (service *AuthServiceImplementation) Login(userid string, userType string, secret string) (string, error) {

	// check if userID not provided
	if userid == "" {
		return "", userIdMustBeProvided
	}

	// TODO: Login doesn't do any validation now,
	// it can be enhanced to add validation of the secret passed

	// Declare the token with the algorithm used for signing
	// TODO: We are using HMAC secure, this can be changed to RSA
	token := jwt.New(jwt.SigningMethodHS256)

	// Use jwt.MapClaims
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["id"] = userid
	claims["role"] = userType
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	// Create the JWT string
	tokenString, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		log.Println("failed to generate token string, error ", err)
		return "", app_errors.InternalServerError
	}

	return tokenString, nil
}

func (service *AuthServiceImplementation) ValidateToken(token string, role string) (*AuthContext, error) {

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing")
		}
		return []byte(service.secretKey), nil
	})

	if err != nil {
		log.Printf("token parse failed %v\n", err)
		return nil, app_errors.Unauthorised
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		if claims["role"] == nil {
			log.Printf("token parse failed: role is empty")
			return nil, app_errors.Unauthorised
		}
		if claims["id"] == nil {
			log.Printf("token parse failed: id is empty")
			return nil, app_errors.Unauthorised
		}

		// validate role
		if claims["role"] != role {
			log.Printf("User role doesn't match, provided %v, required %v\n", claims["role"], role)
			return nil, app_errors.Unauthorised
		}

		return &AuthContext{UserId: fmt.Sprint(claims["id"]), Role: role}, nil
	}

	return nil, InvalidToken
}
