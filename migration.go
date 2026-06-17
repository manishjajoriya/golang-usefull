package db

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type GooseLogger struct {
	log zerolog.Logger
}

func (g *GooseLogger) Fatalf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	g.log.Fatal().Msg(s)
}

func (g *GooseLogger) Printf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	g.log.Info().Msg(s)
}

// RunMigration This Function is run when main.go is run, and it does migration and then close
// that connection so it can be used by pgx.pool in main.go.
func RunMigration(pool *pgxpool.Pool) error {
	connStr := pool.Config().ConnConfig.ConnString()
	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		return err
	}
	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)
	defer func(conn *sql.DB) {
		if err := conn.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close goose connection")
			return
		}
		log.Info().Msg("goose connection closed without any error")
	}(conn)

	goose.SetBaseFS(embedMigrations)

	goose.SetLogger(&GooseLogger{
		log: log.Logger,
	})

	err = goose.Up(conn, "migrations")
	if err != nil {
		return err
	}

	return nil
}
