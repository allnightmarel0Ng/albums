package usecase

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/allnightmarel0Ng/albums/internal/app/authorization/repository"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/utils"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthorizationUseCase interface {
	Authenticate(b64 string) api.Response
	Authorize(jsonWebToken string) api.Response
	Logout(jsonWebToken string) api.Response
	Register(request api.RegistrationRequest) api.Response
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

func (a *authorizationUseCase) Authenticate(b64 string) api.Response {
	rawCredentials, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return &api.AuthenticationResponse{
			Code:  http.StatusBadRequest,
			Error: "error parsing base64",
		}
	}

	credentials := strings.Split(string(rawCredentials), ":")
	if len(credentials) != 2 {
		return &api.AuthenticationResponse{
			Code:  http.StatusBadRequest,
			Error: "wrong authorization format",
		}
	}

	ctx, cancel := utils.DeadlineContext(5)
	defer cancel()

	id, hash, isAdmin, err := a.repo.GetIDPasswordHash(ctx, credentials[0])
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(credentials[1])) != nil {
		return &api.AuthenticationResponse{
			Code:  http.StatusUnauthorized,
			Error: "email or password mismatch",
		}
	}

	result, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      id,
		"isAdmin": isAdmin,
	}).SignedString(a.jwtSecretKey)
	if err != nil {
		return &api.AuthenticationResponse{
			Code:  http.StatusUnauthorized,
			Error: "unable to create jwt key",
		}
	}

	ctx, cancel = utils.DeadlineContext(5)
	defer cancel()
	err = a.repo.AddJWT(ctx, result, 3600)
	if err != nil {
		return &api.AuthenticationResponse{
			Code:  http.StatusInternalServerError,
			Error: "unable to create jwt key",
		}
	}

	return &api.AuthenticationResponse{
		Code:    http.StatusOK,
		Jwt:     result,
		IsAdmin: &isAdmin,
	}
}

func (a *authorizationUseCase) Authorize(jsonWebToken string) api.Response {
	ctx, cancel := utils.DeadlineContext(5)
	defer cancel()

	err := a.repo.FindJWT(ctx, jsonWebToken)
	if err != nil {
		switch err {
		case context.DeadlineExceeded:
		case repository.ErrUnexpected:
			return &api.AuthorizationResponse{
				Code:  http.StatusInternalServerError,
				Error: "jwt storage error",
			}
		default:
			return &api.AuthorizationResponse{
				Code:  http.StatusUnauthorized,
				Error: err.Error(),
			}
		}
	}

	return utils.GetJWTClaims(jsonWebToken, string(a.jwtSecretKey))
}

func (a *authorizationUseCase) Logout(jsonWebToken string) api.Response {
	ctx, cancel := utils.DeadlineContext(5)
	defer cancel()

	err := a.repo.DelJWT(ctx, jsonWebToken)
	if err != nil {
		switch err {
		case context.DeadlineExceeded:
		case repository.ErrUnexpected:
			return &api.ErrorResponse{
				Code:  http.StatusInternalServerError,
				Error: "jwt storage error",
			}
		default:
			return &api.ErrorResponse{
				Code:  http.StatusUnauthorized,
				Error: err.Error(),
			}
		}
	}

	return nil
}

func (a *authorizationUseCase) Register(request api.RegistrationRequest) api.Response {
	if len(request.Password) > 72 {
		return &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "password is too long",
		}
	}

	ctx, cancel := utils.DeadlineContext(10)
	defer cancel()

	found, err := a.repo.FindUserByEmail(ctx, request.Email)
	if err != nil {
		return &api.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "database communication error",
		}
	}

	if found {
		return &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "user with such email already exist",
		}
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), -1)
	if err != nil {
		return &api.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "unable to hash password",
		}
	}

	err = a.repo.AddNewUser(ctx, request.Email, string(hashed), *request.IsAdmin, request.Nickname, request.ImageURL)
	if err != nil {
		return &api.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "database communication error",
		}
	}

	return nil
}
