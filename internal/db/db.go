package db

import (
	"context"
	"log/slog"
	"sync"
	"time"
	"github.com/minguyentt/traverse/internal/tracer"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

// handle database connections with pool for concurrent queries

type PGDB struct {
	*pgxpool.Pool
	logger *slog.Logger
}

var (
	poolInstance *PGDB
	once         sync.Once
)

// TODO: implement db pool connection with Options
func NewPoolConn(
	ctx context.Context,
	connString string,
	logger *slog.Logger,
) (*PGDB, error) {
	once.Do(func() {
		cfg, parseErr := pgxpool.ParseConfig(connString)
		if parseErr != nil {
			panic(parseErr)
		}

		// setup pool configuration
		cfg.MaxConns = 15
		cfg.MinConns = 5
		cfg.MaxConnLifetime = time.Hour
		cfg.MaxConnIdleTime = 30 * time.Minute

		logger.Info("Attempting to connect to database pool...")

		// testing the output with logging
		logger.Warn(
			"connection pool config",
			"user",
			cfg.ConnConfig.User,
			"pass",
			cfg.ConnConfig.Password,
			"host",
			cfg.ConnConfig.Host,
			"database",
			cfg.ConnConfig.Database,
			"port",
			cfg.ConnConfig.Port,
		)

		traceLogger := tracer.NewLogger(logger, "INVALID_TRACE_LOG_LEVEL")

		multTcr := tracer.MultiQuery{
			Tracers: []pgx.QueryTracer{
				otelpgx.NewTracer(),
				&tracelog.TraceLog{
					Logger:   traceLogger,
					LogLevel: tracelog.LogLevelTrace,
				},
			},
		}

		cfg.ConnConfig.Tracer = &multTcr

		pool, err := pgxpool.NewWithConfig(ctx, cfg)
		if err != nil {
			panic(err)
		}

		if err := pool.Ping(ctx); err != nil {
			panic(err)
		}

		poolInstance = &PGDB{
			pool,
			logger,
		}
	})

	return poolInstance, nil
}

func (p *PGDB) GetConnection(ctx context.Context) (*pgxpool.Conn, error) {
	return p.Acquire(ctx)
}
