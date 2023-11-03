package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/m4salah/redroc/libs/storage"
	"github.com/m4salah/redroc/libs/util"
	"github.com/maragudk/migrate"
)

// release is set through the linker at build time, generally from a git sha.
// Used for logging and error reporting.
var release string

var (
	env = flag.String("env", "local", "Env")
)

type Config struct {
	DBHost     string `env:"DB_HOST,notEmpty"`
	DBPort     int    `env:"DB_PORT,notEmpty"`
	DBUser     string `env:"DB_USER,notEmpty"`
	DBPassword string `env:"DB_PASSWORD,notEmpty"`
	DBName     string `env:"DB_NAME,notEmpty"`
}

var config Config

func main() {
	// load env variables
	if err := util.LoadConfig(&config); err != nil {
		panic(err)
	}
	util.InitializeSlog(*env, release)
	os.Exit(start())
}

func start() int {

	if len(os.Args) < 2 {
		slog.Warn("Usage: migrate up|down|to")
		return 1
	}

	if os.Args[1] == "to" && len(os.Args) < 3 {
		slog.Info("Usage: migrate to <version>")
		return 1
	}

	db := storage.NewDatabase(storage.NewDatabaseOptions{
		Host:     config.DBHost,
		Port:     config.DBPort,
		User:     config.DBUser,
		Password: config.DBPassword,
		Name:     config.DBName,
	})

	if err := db.Connect(); err != nil {
		slog.Error("Error connection to database", "error", err)
		return 1
	}

	fsys := os.DirFS("../../libs/storage/migrations")
	var err error
	switch os.Args[1] {
	case "up":
		err = migrate.Up(context.Background(), db.DB.DB, fsys)
	case "down":
		err = migrate.Down(context.Background(), db.DB.DB, fsys)
	case "to":
		err = migrate.To(context.Background(), db.DB.DB, fsys, os.Args[2])
	default:
		slog.Error("Unknown command", slog.String("name", os.Args[1]))
		return 1
	}
	if err != nil {
		slog.Error("Error migrating", "error", err)
		return 1
	}
	slog.Info("Migration complete")
	return 0
}
