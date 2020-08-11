package repo

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func connect(ctx context.Context, dsn string, retries int, logger pgx.Logger) (*pgxpool.Pool, error) {
	delay := time.NewTicker(1 * time.Second)
	timeout := (time.Duration(retries) * time.Second)
	defer delay.Stop()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.MinConns = int32(0)
	config.MaxConns = int32(12) // Adjust this if required
	config.ConnConfig.RuntimeParams["standard_conforming_strings"] = "on"
	config.ConnConfig.RuntimeParams["application_name"] = "Shortr"

	config.ConnConfig.Logger = logger
	config.ConnConfig.LogLevel = pgx.LogLevelError

	timeoutExceeded := time.After(timeout)
	for {
		select {
		case <-timeoutExceeded:
			return nil, errors.New("unable to connect to the database")
		case <-delay.C:
			db, err := pgxpool.ConnectConfig(ctx, config)
			if err == nil {
				return db, nil
			}
		}
	}
}
