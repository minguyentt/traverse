package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"traverse/internal/db"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	cfg "traverse/configs"
)

const (
	DIALECT_DRIVER = "pgx"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", cfg.ENVS.MIGRATIONS.DIR, "directory with migration files")
)

func main() {
	logger := slog.Default().With("area", "migrations")
	flags.Usage = gooseUsage
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}
	command := args[0]
	migrateCtx := context.Background()

	dbString := cfg.ENVS.DB.String()
	pgx, err := db.PoolWithConfig(migrateCtx, dbString)
	if err != nil {
		log.Fatalf("migration error: %v", err)
	}

    logger.Info("acquired connection from the pool")

	defer pgx.Close()

	db := stdlib.OpenDBFromPool(pgx.Pool)
	if err := goose.SetDialect(DIALECT_DRIVER); err != nil {
		log.Fatalf("error setting goose dialect: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("pgx: failed to close DB: %v\n", err)
		}
	}()

	if err := goose.RunContext(migrateCtx, command, db, *dir, args[1:]...); err != nil {
		log.Fatalf("migration: %v: %v", command, err)
	}
}

func gooseUsage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: migrate COMMAND
Examples:
    migrate status
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations`
)
