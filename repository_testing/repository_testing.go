package repository_testing

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"
)

const localEnvFileName = "../../.env.local"

func InitDB(ctx context.Context, t *testing.T, dbURL string) *pg.DB {
	t.Helper()

	opts, err := pg.ParseURL(dbURL)
	if err != nil {
		slog.With("error", err).ErrorContext(ctx, "error parsing db url")
		t.FailNow()
	}

	db := pg.Connect(opts)
	if err = db.Ping(ctx); err != nil {
		slog.With("error", err).ErrorContext(ctx, "error pinging db")
		t.FailNow()
	}

	return db
}

func InitTestingConfig(ctx context.Context, t *testing.T) {
	t.Helper()

	if err := godotenv.Load(localEnvFileName); err != nil {
		slog.With(
			"fileName", localEnvFileName,
			"error", err.Error(),
		).ErrorContext(ctx, "error loading env config")

		t.FailNow()
	}
}

func MustGetEnv(ctx context.Context, t *testing.T, alias string) string {
	t.Helper()

	value := os.Getenv(alias)
	if value == "" {
		slog.With("alias", alias).ErrorContext(ctx, "could not get env variable")
		t.FailNow()
	}

	return value
}
