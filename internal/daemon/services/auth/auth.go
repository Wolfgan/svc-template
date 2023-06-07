package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"service-template/internal/config"
	"service-template/internal/daemon/services/auth/request"
	"service-template/internal/daemon/services/auth/response"
	"service-template/internal/db"
	"service-template/internal/db/users"
	"service-template/internal/utils"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrWrongUsernameOrPassword = errors.New("wrong username or password")
)

type Service struct {
	cfg     *config.Config
	storage *db.Storage
}

func NewService(cfg *config.Config, storage *db.Storage) *Service {
	return &Service{
		cfg:     cfg,
		storage: storage,
	}
}

// SignUp регистрация пользователя.
func (s *Service) SignUp(ctx context.Context, signup *request.SignUp) (*response.SignUp, error) {
	if bytes, err := bcrypt.GenerateFromPassword([]byte(signup.Password), bcrypt.DefaultCost); err != nil {
		return nil, fmt.Errorf("user bcrypt: %w", err)
	} else {
		signup.Password = string(bytes)
	}

	if exists, err := s.storage.Users.Exists(ctx, signup.ToModel()); err != nil {
		return nil, fmt.Errorf("exists user: %w", err)
	} else if exists {
		return nil, ErrUserAlreadyExists
	}

	user, err := s.storage.Users.Create(ctx, signup.ToModel())
	if err != nil {
		return nil, fmt.Errorf("user create: %w", err)
	}

	result := response.SignUp{
		ID:    user.ID,
		Email: user.Email,
		Phone: user.Phone,
	}

	return &result, nil
}

// SignIn ход в аккаунт пользователя.
func (s *Service) SignIn(ctx context.Context, signin *request.SignIn) (*response.SignIn, error) {
	// Получаем пользователя по email или телефону
	user, err := s.storage.Users.Get(ctx, signin.ToModel())
	if err != nil {
		if errors.Is(err, users.ErrNotExists) {
			return nil, ErrWrongUsernameOrPassword
		}

		return nil, fmt.Errorf("user get: %w", err)
	}

	// Проверяем совпадение пароля
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signin.Password)); err != nil {
		return nil, ErrWrongUsernameOrPassword
	}

	// Время жизни токена
	expiration := time.Now().Add(s.cfg.Server.Auth.AccessExpire).Unix()

	// Генерируем полезные данные, которые будут храниться в токене
	payload := jwt.MapClaims{
		"exp":   expiration,
		"sub":   user.ID,
		"email": user.Email,
		"phone": user.Phone,
	}

	// Создаем новый JWT-токен и подписываем его по алгоритму HS256
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(s.cfg.Server.Auth.TokenSecret))
	if err != nil {
		return nil, fmt.Errorf("user token: %w", err)
	}

	fingerprint := fmt.Sprintf("%s:%s:%s", user.Email, user.Phone, s.cfg.Server.Auth.TokenSecret)

	result := response.SignIn{
		AccessToken:  token,
		RefreshToken: utils.SHA256([]byte(fingerprint)),
		Expiration:   expiration,
	}

	return &result, nil
}
