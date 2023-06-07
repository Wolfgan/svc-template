package daemon

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"service-template/internal/config"
	"service-template/internal/daemon/handlers/auth"
	"service-template/internal/daemon/services"
	"service-template/internal/db"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"golang.org/x/exp/slices"
)

type Daemon struct {
	log     *zerolog.Logger
	cfg     *config.Config
	app     *fiber.App
	storage *db.Storage
}

// New create new daemon instance.
func New(log *zerolog.Logger, cfg *config.Config) *Daemon {
	return &Daemon{
		log: log,
		cfg: cfg,
	}
}

func (d *Daemon) Run() error {
	var err error
	d.log.Info().Msg("daemon starting")

	d.log.Info().Msg("init storages")
	if d.storage, err = db.NewStorage(d.cfg, d.log); err != nil {
		return err
	}

	d.log.Info().Msg("init HTTP server")
	d.app = d.initServerHTTP()

	d.log.Info().Msg("init HTTP handlers")
	d.initServerHandlers()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	d.log.Info().Msg("init message queue")

	go func() {
		for {
			select {
			//case errors = <-d.http.Errors():
			//	d.log.Error().Err(errors)

			case <-ctx.Done():
				stop()

				if err = d.Close(); err != nil {
					d.log.Error().Err(err)
				}

				return
			}
		}
	}()

	if !strings.HasPrefix(d.cfg.Server.Port, ":") {
		d.cfg.Server.Port = ":" + d.cfg.Server.Port
	}

	return d.app.Listen(d.cfg.Server.Port)
}

// Close закрывает соединения с БД, и глушит сервер.
func (d *Daemon) Close() error {
	var errs []error

	if d.storage != nil {
		d.log.Info().Msg("close storage")

		if err := d.storage.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if d.app != nil {
		d.log.Info().Msg("HTTP server shutdown")

		if err := d.app.Shutdown(); errs != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (d *Daemon) initServerHTTP() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Показываем пользователю ошибку, только если она в списке разрешенных
			var externalErrors = []int{fiber.StatusBadRequest}
			if slices.Contains(externalErrors, code) {
				return c.Status(code).SendString(err.Error())
			}

			return c.SendStatus(code)
		},
	})

	// логирование запросов
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger:   d.log,
		Fields:   []string{fiberzerolog.FieldStatus, fiberzerolog.FieldMethod, fiberzerolog.FieldURL, fiberzerolog.FieldError},
		Messages: []string{"Server", "Client", "Success"},
	}))

	return app
}

func (d *Daemon) initServerHandlers() {
	interactor := services.NewInteractor(d.cfg, d.storage)

	authHandler := auth.NewHandler(d.log, interactor)

	// Группа обработчиков, которые доступны неавторизованным пользователям
	publicGroup := d.app.Group("")
	publicGroup.Post("/signup", authHandler.SignUp)
	publicGroup.Post("/signin", authHandler.SignIn)
	//publicGroup.Get("/", authHandler.Root)
	//
	//// Группа обработчиков, которые требуют авторизации
	//authorizedGroup := app.Group("")
	//authorizedGroup.Use(jwtware.New(jwtware.Config{
	//	SigningKey: handlers.SecretKey,
	//	ContextKey: handlers.ContextKeyUser,
	//}))
	//authorizedGroup.Get("/profile", userHandler.Profile)

}
