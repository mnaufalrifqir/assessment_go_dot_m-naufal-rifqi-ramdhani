package auth

import (
	"api-dot/user"
	"api-dot/utils"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(payload user.User) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
}

func NewService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) GenerateToken(payload user.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()

	expirationTime := time.Now().Add(time.Duration(15) * time.Minute)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = payload.ID
	claims["role"] = payload.Role
	claims["exp"] = expirationTime.Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := token.SignedString([]byte(utils.GetConfig("TOKEN_SECRET")))

	if err != nil {
		return "", fmt.Errorf("generating JWT Token failed: %w", err)
	}

	return tokenString, nil
}

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(utils.GetConfig("TOKEN_SECRET")), nil
	})

	if err != nil {
		return token, err
	}

	return token, nil
}
