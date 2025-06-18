package db

import (
	"context"
	"log/slog"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

func NewPoolConn(
	ctx context.Context,
	connString string,
	logger *slog.Logger,
) (*PGDB, error) {
	var err error
	once.Do(func() {
		cfg, parseErr := pgxpool.ParseConfig(connString)
		if err != nil {
			err = parseErr
			return
		}

		logger.Info("Attempting to connect to database pool...")

		cfg.BeforeConnect = func(ctx context.Context, c *pgx.ConnConfig) error {
			//TODO: could be better than running the env vars in main
			// 		have the pgx pool configured in here instead maybe?

				// Fetch credentials dynamically
				// user, password, host, database, port, err := getCredentials()
				// if err != nil {
				// 	return err
				// }
				//
				// // Update the connection config
				// connConfig.User = user
				// connConfig.Password = password
				// connConfig.Host = host
				// connConfig.Database = database
				// connConfig.Port = port
				//
				// return nil

			return nil
		}

		//TODO: im not sure if i need this...
		cfg.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
			// current user_name connections
			connInfo := conn.PgConn().ParameterStatus("session_authorization")

			// logger.Info("checking for existing connections...", "user", connInfo)

			// handle not in transaction
			if inTx := conn.PgConn().TxStatus() != 'I'; inTx {
				logger.Info("Connection in transaction, unable to acquire", "user", connInfo)
				return false
			}
			// ping database
			err = conn.Ping(ctx)
			if err != nil {
				logger.Error("connection failed to ping database", "error", err)
			}

			return true
		}

		pool, cfgErr := pgxpool.NewWithConfig(ctx, cfg)
		if err != nil {
			err = cfgErr
			return
		}

		poolInstance = &PGDB{
			pool,
			logger,
		}
	})

	return poolInstance, nil
}

// return existing pool instance
func GetPool() *PGDB {
	return poolInstance
}

func (p *PGDB) GetConnection(ctx context.Context) (*pgxpool.Conn, error) {
	return p.Acquire(ctx)
}
