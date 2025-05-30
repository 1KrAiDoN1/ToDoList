package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	todo "helloapp"
	"helloapp/pkg/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	salt       = "abuyzrfrfznj12344321"
	signingKey = "abuyzrfrfznj12344321"
	tokenTTL   = 12 * time.Hour
)

type AuthService struct {
	repos repository.Authorization
}

func NewAuthService(repository repository.Authorization) *AuthService {
	return &AuthService{
		repos: repository,
	}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repos.CreateUser(user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repos.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}
