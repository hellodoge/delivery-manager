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
	"regexp"
	"strings"
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

	usernameRegex = "^[a-z]+[a-z1-9_]*$"
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

	user.Username = strings.ToLower(user.Username)
	matched, err := regexp.MatchString(usernameRegex, user.Username)
	if err != nil {
		return -1, err
	} else if !matched {
		return -1, response.ErrorResponseParameters{
			Message:    "Username must start with a letter and consist of alphanumeric characters/underscores",
			StatusCode: 0,
		}
	}

	salt := make([]byte, pwSaltBytes)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return 0, err
	}
	user.PasswordSalt = hex.EncodeToString(salt)

	hashed := pbkdf2.Key([]byte(user.Password), salt, hashIterations, pwHashBytes, sha256.New)
	user.PasswordHash = hex.EncodeToString(hashed)

	id, err := s.repo.CreateUser(user)
	if err == repository.ErrUserExists {
		return -1, response.ErrorResponseParameters{
			Message:    "User already exists",
			StatusCode: http.StatusUnauthorized,
		}
	} else if err != nil {
		return -1, err
	}
	return id, nil
}

type TokenClaims struct {
	jwt.StandardClaims
	UserId       int    `json:"user_id"`
	UserHashPart string `json:"user_hash"`
}

func (s *AuthService) GenerateToken(refreshToken string) (string, error) {

	saved, err := s.cache.GetSavedFields(refreshToken)
	if err != nil {
		return "", err
	} else if saved == nil {
		return "", response.ErrorResponseParameters{
			Internal:   err,
			Message:    "Refresh token not found",
			StatusCode: http.StatusUnauthorized,
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.config.TokenLifetime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId:       saved.UserID,
		UserHashPart: saved.UserHashPart,
	})

	return token.SignedString([]byte(os.Getenv("SIGNING_KEY")))
}

func (s *AuthService) GenerateRefreshToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return "", response.ErrorResponseParameters{
			Internal:   err,
			Message:    "User not found",
			StatusCode: http.StatusUnauthorized,
		}
	}

	salt, err := hex.DecodeString(user.PasswordSalt)
	if err != nil {
		return "", err
	}

	hashed := pbkdf2.Key([]byte(password), salt, hashIterations, pwHashBytes, sha256.New)

	if user.PasswordHash != hex.EncodeToString(hashed) {
		return "", response.ErrorResponseParameters{
			Internal:   fmt.Errorf("user %s failed login", user.Username),
			Message:    "Invalid password",
			StatusCode: http.StatusUnauthorized,
		}
	}

	token, err := s.cache.NewRefreshToken(&cache.RefreshTokenSavedFields{
		UserID:       user.Id,
		UserHashPart: user.PasswordHash[:claimsHashPartLen],
	})
	return token, err
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{},
		func(token *jwt.Token) (interface{}, error) { return []byte(os.Getenv("SIGNING_KEY")), nil })

	if err != nil {
		return -1, response.ErrorResponseParameters{
			Internal:   err,
			Message:    "Invalid Credentials",
			StatusCode: http.StatusUnauthorized,
		}
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return -1, response.ErrorResponseParameters{
			Internal:   fmt.Errorf("credentials of invalid type: %T", token.Claims),
			Message:    "Unexpected Credentials Type",
			StatusCode: http.StatusUnauthorized,
		}
	}

	if s.config.CheckHash {
		user, err := s.repo.GetUserByID(claims.UserId)
		if err != nil {
			return -1, response.ErrorResponseParameters{
				Internal:   err,
				IsInternal: true,
				Message:    "User not found",
				StatusCode: http.StatusUnauthorized,
			}
		}

		if claims.UserHashPart != user.PasswordHash[:claimsHashPartLen] {
			return -1, response.ErrorResponseParameters{
				Internal:   fmt.Errorf("jwt cache sign is valid, but user %s hash don't match", user.Username),
				IsInternal: true,
				Message:    "User not found",
				StatusCode: http.StatusUnauthorized,
			}
		}
	}

	return claims.UserId, nil
}
