package migrator

import (
	"fmt"
	"os"
	"strings"

	"service-template/internal/config"
	"service-template/pkg/drivers/postgres"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
)

// MigrateCommands возвращает команду для работы с миграциями.
func MigrateCommands() *cli.Command {
	var migrations *migrate.Migrations
	var cfg *config.Config

	return &cli.Command{
		Name:  "db",
		Usage: "database migrations",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "migrations",
				Usage:    "Load migrations from `DIR`",
				Aliases:  []string{"m"},
				Required: true,
				Action: func(c *cli.Context, path cli.Path) error {
					var err error
					if cfg, err = config.New(c.String("config")); err != nil {
						return err
					}

					if err = cfg.Validate(); err != nil {
						return err
					}

					migrations = migrate.NewMigrations(migrate.WithMigrationsDirectory(c.String("migrations")))
					return migrations.Discover(os.DirFS(path))
				},
			},
		},
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					db, err := postgres.NewPostgresDB(cfg.Postgres)
					if err != nil {
						return err
					}
					defer db.Close()

					migrator := migrate.NewMigrator(db, migrations)

					return migrator.Init(c.Context)
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					db, err := postgres.NewPostgresDB(cfg.Postgres)
					if err != nil {
						return err
					}
					defer db.Close()

					migrator := migrate.NewMigrator(db, migrations)

					var group *migrate.MigrationGroup
					if group, err = migrator.Migrate(c.Context); err != nil {
						return err
					}

					if group.IsZero() {
						fmt.Println("there are no new migrations to run (database is up to date)")

						return nil
					}

					log.Info().Msgf("migrated to %s", group)

					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					db, err := postgres.NewPostgresDB(cfg.Postgres)
					if err != nil {
						return err
					}
					defer db.Close()

					migrator := migrate.NewMigrator(db, migrations)

					var group *migrate.MigrationGroup
					if group, err = migrator.Rollback(c.Context); err != nil {
						return err
					}

					if group.IsZero() {
						fmt.Println("there are no groups to roll back")

						return nil
					}

					fmt.Printf("rolled back %s\n", group)

					return nil
				},
			},
			{
				Name:  "lock",
				Usage: "lock migrations",
				Action: func(c *cli.Context) error {
					db, err := postgres.NewPostgresDB(cfg.Postgres)
					if err != nil {
						return err
					}
					defer db.Close()

					migrator := migrate.NewMigrator(db, migrations)

					return migrator.Lock(c.Context)
				},
			},
			{
				Name:  "unlock",
				Usage: "unlock migrations",
				Action: func(c *cli.Context) error {
					db, err := postgres.NewPostgresDB(cfg.Postgres)
					if err != nil {
						return err
					}
					defer db.Close()

					migrator := migrate.NewMigrator(db, migrations)

					return migrator.Unlock(c.Context)
				},
			},
			{
				Name:  "create_go",
				Usage: "create Go migration",
				Action: func(c *cli.Context) error {
					db, err := postgres.NewPostgresDB(cfg.Postgres)
					if err != nil {
						return err
					}
					defer db.Close()

					migrator := migrate.NewMigrator(db, migrations)

					name := strings.Join(c.Args().Slice(), "_")

					var mf *migrate.MigrationFile
					if mf, err = migrator.CreateGoMigration(c.Context, name); err != nil {
						return err
					}

					fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)

					return nil
				},
			},
			{
				Name:  "create_sql",
				Usage: "create up and down SQL migrations",
				Action: func(c *cli.Context) error {
					db, err := postgres.NewPostgresDB(cfg.Postgres)
					if err != nil {
						return err
					}
					defer db.Close()

					migrator := migrate.NewMigrator(db, migrations)

					name := strings.Join(c.Args().Slice(), "_")
					var files []*migrate.MigrationFile
					if files, err = migrator.CreateSQLMigrations(c.Context, name); err != nil {
						return err
					}

					for _, mf := range files {
						fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
					}

					return nil
				},
			},
			{
				Name:  "status",
				Usage: "print migrations status",
				Action: func(c *cli.Context) error {
					db, err := postgres.NewPostgresDB(cfg.Postgres)
					if err != nil {
						return err
					}
					defer db.Close()

					migrator := migrate.NewMigrator(db, migrations)

					var ms migrate.MigrationSlice
					if ms, err = migrator.MigrationsWithStatus(c.Context); err != nil {
						return err
					}

					fmt.Printf("migrations: %s\n", ms)
					fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
					fmt.Printf("last migration group: %s\n", ms.LastGroup())

					return nil
				},
			},
			{
				Name:  "mark_applied",
				Usage: "mark migrations as applied without actually running them",
				Action: func(c *cli.Context) error {
					db, err := postgres.NewPostgresDB(cfg.Postgres)
					if err != nil {
						return err
					}
					defer db.Close()

					migrator := migrate.NewMigrator(db, migrations)

					var group *migrate.MigrationGroup
					if group, err = migrator.Migrate(c.Context, migrate.WithNopMigration()); err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Println("there are no new migrations to mark as applied")

						return nil
					}

					fmt.Printf("marked as applied %s\n", group)

					return nil
				},
			},
		},
	}
}
