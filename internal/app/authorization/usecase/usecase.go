package usecase

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/allnightmarel0Ng/albums/internal/app/authorization/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthorizationUseCase interface {
	Authorize(b64 string) (string, int, error)
}

type authorizationUseCase struct {
	repo         repository.AuthorizationRepository
	jwtSecretKey []byte
}

func NewAuthorizationUseCase(repo repository.AuthorizationRepository, jwtSecretKey []byte) AuthorizationUseCase {
	return &authorizationUseCase{
		repo:         repo,
		jwtSecretKey: jwtSecretKey,
	}
}

func (a *authorizationUseCase) Authorize(b64 string) (string, int, error) {
	rawCredentials, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", http.StatusBadRequest, errors.New("error parsing base64")
	}

	credentials := strings.Split(string(rawCredentials), ":")
	if len(credentials) != 2 {
		return "", http.StatusBadRequest, errors.New("wrong authorization format")
	}

	id, hash, isAdmin, err := a.repo.GetIDPasswordHash(credentials[0])
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(credentials[1])) != nil {
		return "", http.StatusUnauthorized, errors.New("email or password mismatch")
	}

	result, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      id,
		"isAdmin": isAdmin,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}).SignedString(a.jwtSecretKey)
	if err != nil {
		return "", http.StatusInternalServerError, errors.New("unable to create jwt key")
	}

	return result, http.StatusOK, nil
}
