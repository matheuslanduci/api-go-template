package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/storage/redis/v3"
	"github.com/gofiber/utils/v2"
	"github.com/jmoiron/sqlx"
	"github.com/nrednav/cuid2"
	"golang.org/x/crypto/bcrypt"
	"matheuslanduci.com/api-fiber/database/models"
	"matheuslanduci.com/api-fiber/dto"
)

type AuthenticationService struct {
	db    *sqlx.DB
	redis *redis.Storage
}

type SessionResponse struct {
	SessionToken  string
	RememberToken *string
}

func NewAuthenticationService(db *sqlx.DB, redis *redis.Storage) *AuthenticationService {
	return &AuthenticationService{
		db:    db,
		redis: redis,
	}
}

const (
	ErrUserNotFound    = "user_not_found"
	ErrInvalidPassword = "invalid_password"
	ErrInvalidSession  = "invalid_session"
)

const SessionLifetime = time.Hour * 1
const RememberLifetime = time.Hour * 24 * 7

func GenerateSessionToken() string {
	token := utils.UUIDv4()
	cuid := cuid2.Generate()

	return fmt.Sprintf("@app-%s:%s", token, cuid)
}

func GenerateRememberToken() string {
	token := utils.UUIDv4()
	cuid := cuid2.Generate()

	return fmt.Sprintf("@app-remember-%s:%s", token, cuid)
}

func (service *AuthenticationService) CreateSessionWithPassword(request *dto.CreateSessionWithPasswordRequest) (*SessionResponse, error) {
	user := models.User{}

	err := service.db.Get(&user, `
		SELECT id, first_name, last_name, email, password, created_at, updated_at
		FROM users WHERE email = $1
	`, request.Email)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))

	if err != nil {
		return nil, errors.New(ErrInvalidPassword)
	}

	sessionToken := GenerateSessionToken()

	formattedUser, err := json.Marshal(
		&models.User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	)

	if err != nil {
		return nil, err
	}

	err = service.redis.Set(sessionToken, formattedUser, SessionLifetime)

	if err != nil {
		return nil, err
	}

	var rememberToken *string

	if request.Remember {
		rememberToken = new(string)

		*rememberToken = GenerateRememberToken()

		err = service.redis.Set(*rememberToken, formattedUser, RememberLifetime)

		if err != nil {
			return nil, err
		}
	}

	return &SessionResponse{
		SessionToken:  sessionToken,
		RememberToken: rememberToken,
	}, nil
}

func (service *AuthenticationService) GetSession(sessionToken string) (*models.User, error) {
	formattedUser, err := service.redis.Get(sessionToken)

	if err != nil {
		return nil, err
	}

	if formattedUser == nil {
		return nil, errors.New(ErrInvalidSession)
	}

	user := models.User{}

	err = json.Unmarshal([]byte(formattedUser), &user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (service *AuthenticationService) DeleteSession(sessionToken string, rememberToken *string) error {
	if rememberToken != nil {
		err := service.redis.Delete(*rememberToken)

		if err != nil {
			return err
		}
	}

	return service.redis.Delete(sessionToken)
}

func (service *AuthenticationService) CreateSessionWithRememberToken(rememberToken string) (*SessionResponse, error) {
	formattedUser, err := service.redis.Get(rememberToken)

	if err != nil {
		return nil, err
	}

	user := models.User{}

	err = json.Unmarshal([]byte(formattedUser), &user)

	if err != nil {
		return nil, err
	}

	newRememberToken := GenerateRememberToken()

	err = service.redis.Set(newRememberToken, formattedUser, time.Hour*24*7)

	if err != nil {
		return nil, err
	}

	err = service.redis.Delete(rememberToken)

	if err != nil {
		return nil, err
	}

	sessionToken := GenerateSessionToken()

	formattedUser, err = json.Marshal(
		&models.User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	)

	if err != nil {
		return nil, err
	}

	err = service.redis.Set(sessionToken, formattedUser, SessionLifetime)

	if err != nil {
		return nil, err
	}

	return &SessionResponse{
		SessionToken:  sessionToken,
		RememberToken: &newRememberToken,
	}, nil
}
