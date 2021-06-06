package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hellodoge/delivery-manager/dm"
	"github.com/hellodoge/delivery-manager/internal/cache"
	"github.com/hellodoge/delivery-manager/internal/repository"
	"github.com/hellodoge/delivery-manager/pkg/response"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	// migrations/000001_init.up.sql
	pwSaltBytes = 32 / 2
	pwHashBytes = 64 / 2

	// pbkdf2
	hashIterations = 2 << 4

	// jwt
	claimsHashPartLen = 16
)

type AuthService struct {
	repo   repository.Authorization
	cache  cache.RefreshTokens
	config AuthServiceConfig
}

func NewAuthService(repo repository.Authorization, cache cache.RefreshTokens, config AuthServiceConfig) *AuthService {
	return &AuthService{
		repo:   repo,
		cache:  cache,
		config: config,
	}
}

func (s *AuthService) CreateUser(user dm.User) (int, error) {

	salt := make([]byte, pwSaltBytes)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return 0, err
	}
	user.PasswordSalt = hex.EncodeToString(salt)

	hashed := pbkdf2.Key([]byte(user.Password), salt, hashIterations, pwHashBytes, sha256.New)
	user.PasswordHash = hex.EncodeToString(hashed)

	id, err := s.repo.CreateUser(user)
	if err != nil {
		return -1, response.ErrorResponse{
			Internal:   err,
			Message:    "User already exists",
			StatusCode: http.StatusUnauthorized,
		}
	}
	return id, nil
}

type TokenClaims struct {
	jwt.StandardClaims
	UserId       int    `json:"user_id"`
	UserHashPart string `json:"user_hash"`
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {

	user, err := s.repo.GetUser(username)
	if err != nil {
		return "", response.ErrorResponse{
			Internal:   err,
			Message:    "User not found",
			StatusCode: http.StatusUnauthorized,
		}
	}

	salt, err2 := hex.DecodeString(user.PasswordSalt)
	if err2 != nil {
		return "", err
	}

	hashed := pbkdf2.Key([]byte(password), salt, hashIterations, pwHashBytes, sha256.New)

	if user.PasswordHash != hex.EncodeToString(hashed) {
		return "", response.ErrorResponse{
			Internal:   fmt.Errorf("user %s failed login", user.Username),
			Message:    "Invalid password",
			StatusCode: http.StatusUnauthorized,
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.config.TokenLifetime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId:       user.Id,
		UserHashPart: user.PasswordHash[:claimsHashPartLen],
	})

	return token.SignedString([]byte(os.Getenv("SIGNING_KEY")))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{},
		func(token *jwt.Token) (interface{}, error) { return []byte(os.Getenv("SIGNING_KEY")), nil })

	if err != nil {
		return -1, response.ErrorResponse{
			Internal:   err,
			Message:    "Invalid Credentials",
			StatusCode: http.StatusUnauthorized,
		}
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return -1, response.ErrorResponse{
			Internal:   fmt.Errorf("credentials of invalid type: %T", token.Claims),
			Message:    "Unexpected Credentials Type",
			StatusCode: http.StatusUnauthorized,
		}
	}

	if s.config.CheckHash {
		user, err := s.repo.GetUserByID(claims.UserId)
		if err != nil {
			return -1, response.ErrorResponse{
				Internal:   err,
				IsInternal: true,
				Message:    "User not found",
				StatusCode: http.StatusUnauthorized,
			}
		}

		if claims.UserHashPart != user.PasswordHash[:claimsHashPartLen] {
			return -1, response.ErrorResponse{
				Internal:   fmt.Errorf("jwt cache sign is valid, but user %s hash don't match", user.Username),
				IsInternal: true,
				Message:    "User not found",
				StatusCode: http.StatusUnauthorized,
			}
		}
	}

	return claims.UserId, nil
}
