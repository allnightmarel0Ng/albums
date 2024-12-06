package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/allnightmarel0Ng/albums/internal/domain/repository"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/redis"
)

var (
	ErrUnexpected  = errors.New("unexpected error")
	ErrJWTNotFound = errors.New("jwt not found")
)

type AuthorizationRepository interface {
	GetIDPasswordHash(ctx context.Context, email string) (int, string, bool, error)
	AddNewUser(ctx context.Context, email, password_hash string, isAdmin bool, nickname, imageURL string) error
	FindUserByEmail(ctx context.Context, email string) (bool, error)
	AddJWT(ctx context.Context, jwt string, expirationSeconds int) error
	FindJWT(ctx context.Context, jwt string) error
	DelJWT(ctx context.Context, jwt string) error
}

type authorizationRepository struct {
	users repository.UserRepository
	redis redis.Client
}

func NewAuthorizationRepository(users repository.UserRepository, redis redis.Client) AuthorizationRepository {
	return &authorizationRepository{
		users: users,
		redis: redis,
	}
}

func (a *authorizationRepository) GetIDPasswordHash(ctx context.Context, email string) (int, string, bool, error) {
	select {
	case <-ctx.Done():
		return 0, "", false, ctx.Err()
	default:
		return a.users.GetIDPasswordHash(ctx, email)
	}
}

func (a *authorizationRepository) AddNewUser(ctx context.Context, email, password_hash string, isAdmin bool, nickname, imageURL string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return a.users.AddNewUser(ctx, email, password_hash, isAdmin, nickname, imageURL)
	}
}

func (a *authorizationRepository) AddJWT(ctx context.Context, jwt string, expirationSeconds int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		log.Println(time.Duration(expirationSeconds) * time.Second)
		err := a.redis.Set(ctx, jwt, "", time.Duration(expirationSeconds)*time.Second)
		if err != nil {
			return ErrUnexpected
		}
		return nil
	}
}

func (a *authorizationRepository) FindJWT(ctx context.Context, jwt string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, err := a.redis.Get(ctx, jwt)
		if err != nil {
			if err == redis.ErrNotFound {
				return ErrJWTNotFound
			} else {
				return ErrUnexpected
			}
		}
		return nil
	}
}

func (a *authorizationRepository) DelJWT(ctx context.Context, jwt string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		err := a.redis.Del(ctx, jwt)
		if err != nil {
			if err == redis.ErrNotFound {
				return ErrJWTNotFound
			} else {
				return ErrUnexpected
			}
		}
		return nil
	}
}

func (a *authorizationRepository) FindUserByEmail(ctx context.Context, email string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		return a.users.FindUserByEmail(ctx, email)
	}
}
