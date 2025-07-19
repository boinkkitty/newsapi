package main

import (
	"fmt"
	"github.com/boinkkitty/newsapi/internal/migration"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
)

func main() {
	//db, err := postgres.NewDB(&postgres.Config{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//db.AddQueryHook(bundebug.NewQueryHook(
	//	bundebug.WithEnabled(false),
	//	bundebug.FromEnv(),
	//))

	app := &cli.App{
		Name: "migrate",
		Commands: []*cli.Command{
			newMigrationCommand(migrate.NewMigrator(nil, migration.New(), migrate.WithMarkAppliedOnSuccess(true))),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func newMigrationCommand(m *migrate.Migrator) *cli.Command {
	return &cli.Command{
		Name:  "migrate",
		Usage: "database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migrations table",
				Action: func(ctx *cli.Context) error {
					return m.Init(ctx.Context)
				},
			},
			{
				Name:  "up",
				Usage: "run migrations up",
				Action: func(ctx *cli.Context) error {
					if err := m.Lock(ctx.Context); err != nil {
						return err
					}
					defer m.Unlock(ctx.Context)

					group, err := m.Migrate(ctx.Context)
					if err != nil {
						return err
					}

					if group.IsZero() {
						fmt.Println("no migrations available, database is up to date\n")
						return nil
					}
					fmt.Printf("migrated to %s\n", group)
					return nil
				},
			},
			{
				Name:  "down",
				Usage: "run migrations down",
				Action: func(ctx *cli.Context) error {
					if err := m.Lock(ctx.Context); err != nil {
						return err
					}
					defer m.Unlock(ctx.Context)

					group, err := m.Rollback(ctx.Context)
					if err != nil {
						return err
					}

					if group.IsZero() {
						fmt.Println("no groups to rollback\n")
						return nil
					}
					fmt.Printf("rolled back to %s\n", group)
					return nil
				},
			},
			{
				Name:  "create",
				Usage: "create up and down sql migrations",
				Action: func(ctx *cli.Context) error {
					name := strings.Join(ctx.Args().Slice(), "_")
					files, err := m.CreateTxSQLMigrations(ctx.Context, name)
					if err != nil {
						return err
					}
					for _, f := range files {
						fmt.Printf("created migrations %s (%s)\n", f.Name, f.Path)
					}
					return nil
				},
			},
			{
				Name:  "status",
				Usage: "show status of migrations",
				Action: func(ctx *cli.Context) error {
					ms, err := m.MigrationsWithStatus(ctx.Context)
					if err != nil {
						return err
					}
					fmt.Printf("migrations: %s\n", ms)
					fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
					fmt.Printf("last migration group: %s\n", ms.LastGroup())
					return nil
				},
			},
		},
	}
}
