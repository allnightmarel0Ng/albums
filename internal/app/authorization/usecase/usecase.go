package usecase

import (
	"encoding/base64"
	"errors"
	"log"
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
	unsplittedCredentials, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", http.StatusBadRequest, errors.New("error parsing base64")
	}

	credentials := strings.Split(string(unsplittedCredentials), ":")
	if len(credentials) != 2 {
		return "", http.StatusBadRequest, errors.New("wrong authorization format")
	}

	user, hash, err := a.repo.Authorize(credentials[0])
	log.Printf("user: %v, hash: %s, password: %s", user, hash, credentials[1])
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(credentials[1])) != nil {
		return "", http.StatusNotFound, errors.New("email or password mismatch")
	}

	result, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}).SignedString(a.jwtSecretKey)
	if err != nil {
		return "", http.StatusInternalServerError, errors.New("unable to create jwt key")
	}

	return result, http.StatusOK, nil
}
