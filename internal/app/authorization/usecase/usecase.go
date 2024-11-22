package usecase

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/allnightmarel0Ng/albums/internal/app/authorization/repository"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthorizationUseCase interface {
	Authorize(b64 string) (string, int, error)
}

type authorizationUseCase struct {
	repo         repository.AuthorizationRepository
	jwtSecretKey string
}

func NewAuthorizationUseCase(repo repository.AuthorizationRepository, jwtSecretKey string) AuthorizationUseCase {
	return &authorizationUseCase{
		repo:         repo,
		jwtSecretKey: jwtSecretKey,
	}
}

func (a *authorizationUseCase) Authorize(b64 string) (string, int, error) {
	unsplittedCredentials, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", http.StatusBadRequest, errors.New("error parsing base64")
	}

	credentials := strings.Split(string(unsplittedCredentials), ":")
	if len(credentials) != 2 {
		return "", http.StatusBadRequest, errors.New("wrong authorization format")
	}

	encrypted, err := bcrypt.GenerateFromPassword([]byte(credentials[1]), bcrypt.DefaultCost)
	if err != nil {
		return "", http.StatusInternalServerError, errors.New("encryption error")
	}

	user, err := a.repo.Authorize(credentials[0], string(encrypted))
	if err != nil {
		return "", http.StatusNotFound, err
	}

	result, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}).SignedString(a.jwtSecretKey)
	if err != nil {
		return "", http.StatusInternalServerError, errors.New("unable to create jwt key")
	}

	return result, http.StatusOK, nil
}
