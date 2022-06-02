package services

import (
	"fmt"
	"log"
	"time"

	"github.com/croisade/chimichanga/pkg/conf"
	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	CreateToken()
	CreateRefreshToken()
}

type JWTAuthService struct{}

func NewJWTAuthService() JWTAuthService {
	return JWTAuthService{}
}
func (j JWTAuthService) CreateToken() (string, error) {
	config, err := conf.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	type MyCustomClaims struct {
		Foo string `json:"foo"`
		jwt.StandardClaims
	}

	// Create the Claims
	claims := MyCustomClaims{
		"bar",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			Issuer:    "CorroYouRun",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(config.JWTSecret))
	return ss, err
}

func (j JWTAuthService) CreateRefreshToken() (string, error) {
	config, err := conf.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	type MyCustomClaims struct {
		Foo string `json:"foo"`
		jwt.StandardClaims
	}

	claims := MyCustomClaims{
		"bar",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 168).Unix(),
			Issuer:    "CorroYouRun",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(config.JWTSecret))
	fmt.Printf("%v %v", ss, err)
	return ss, err
}

func (j JWTAuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	config, err := conf.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
