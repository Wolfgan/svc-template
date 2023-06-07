package main

import (
	"os"
	"service-template/internal/config"
	"service-template/internal/daemon"
	"service-template/pkg/migrator"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

const (
	// ServiceName содержит имя сервиса. Выводится в логах и при вызове help.
	ServiceName = "template"

	// ServiceUsage содержит краткое описание сервиса.
	ServiceUsage = "template service"

	// ServiceDescription Отображает полное описание сервиса.
	ServiceDescription = "template service"
)

var (
	// ServiceVersion содержит номер версии приложения: `-ldflags "-X main.ServiceVersion=${VERSION}"`.
	ServiceVersion = "0.0.0-develop"
)

///go:embed migrations/*.sql
//var sqlMigrations embed.FS

func main() {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().
		Timestamp().Logger().Level(zerolog.DebugLevel)

	var err error
	var cfg *config.Config

	app := &cli.App{
		Name:        ServiceName,
		Usage:       ServiceUsage,
		Version:     ServiceVersion,
		Description: ServiceDescription,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Usage:    "Load configuration from `FILE`",
				Aliases:  []string{"c"},
				Required: true,
			},
		},

		Commands: []*cli.Command{
			migrator.MigrateCommands(),
		},

		// Перед выполнением action`s инициализируем параметры
		Before: func(c *cli.Context) error {
			if cfg, err = config.New(c.String("config")); err != nil {
				return err
			}

			if err = cfg.Validate(); err != nil {
				return err
			}

			log.Level(cfg.Logger.Level)

			if !cfg.Logger.Console {
				logger = log.With().
					Str("service", ServiceName).
					Str("version", ServiceVersion).Logger().Output(os.Stderr)
			}

			if cfg.Logger.WithCaller {
				logger = log.With().Caller().Logger()
			}

			return nil
		},

		Action: func(c *cli.Context) error {
			dmn := daemon.New(&logger, cfg)

			return dmn.Run()
		},
	}

	if err = app.Run(os.Args); err != nil {
		logger.Fatal().Msg(err.Error())
	}

	//app := fiber.New()
	//
	//authStorage := cache.NewCache[string, model.User]()
	//authHandler := &handlers.Auth{Storage: authStorage}
	//userHandler := &handlers.User{Storage: authStorage}
	//
	//// Группа обработчиков, которые доступны неавторизованным пользователям
	//publicGroup := app.Group("")
	//publicGroup.Post("/register", authHandler.Register)
	//publicGroup.Post("/login", authHandler.Login)
	//publicGroup.Get("/", authHandler.Root)
	//
	//// Группа обработчиков, которые требуют авторизации
	//authorizedGroup := app.Group("")
	//authorizedGroup.Use(jwtware.New(jwtware.Config{
	//	SigningKey: handlers.SecretKey,
	//	ContextKey: handlers.ContextKeyUser,
	//}))
	//authorizedGroup.Get("/profile", userHandler.Profile)
	//
	//logrus.Fatal(app.Listen(":8080"))
}
