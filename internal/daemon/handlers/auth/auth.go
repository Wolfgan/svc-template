package auth

import (
	"errors"
	"fmt"
	"service-template/internal/daemon/services"
	"service-template/internal/daemon/services/auth"
	"service-template/internal/daemon/services/auth/request"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Handler struct {
	log        *zerolog.Logger
	interactor *services.Interactor
}

func NewHandler(log *zerolog.Logger, interactor *services.Interactor) *Handler {
	return &Handler{
		log:        log,
		interactor: interactor,
	}
}

// SignUp Обработчик HTTP-запросов на вход в аккаунт пользователя.
func (h *Handler) SignUp(c *fiber.Ctx) error {
	ctx := h.log.WithContext(c.Context())

	signup := request.SignUp{}
	if err := c.BodyParser(&signup); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Errorf("body parser: %w", err).Error())
	}

	if err := signup.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response, err := h.interactor.Auth.SignUp(ctx, &signup)
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response)
}

// SignIn Обработчик HTTP-запросов на вход в аккаунт пользователя.
func (h *Handler) SignIn(c *fiber.Ctx) error {
	ctx := h.log.WithContext(c.Context())

	signin := request.SignIn{}
	if err := c.BodyParser(&signin); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Errorf("body parser: %w", err).Error())
	}

	if err := signin.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response, err := h.interactor.Auth.SignIn(ctx, &signin)
	if err != nil {
		if errors.Is(err, auth.ErrWrongUsernameOrPassword) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response)
}
