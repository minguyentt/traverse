package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// handle database connections with pool for concurrent queries

type PGDB struct {
	*pgxpool.Pool
	logger *zap.SugaredLogger
}

var (
	poolInstance *PGDB
	once         sync.Once
)

func NewPoolConn(
	ctx context.Context,
	connString string,
	logger *zap.SugaredLogger,
) (*PGDB, error) {
	var err error
	cpl := logger.Named("Connection Pool")
	once.Do(func() {
		cfg, parseErr := pgxpool.ParseConfig(connString)
		if err != nil {
			err = parseErr
			return
		}

		cfg.BeforeConnect = func(ctx context.Context, c *pgx.ConnConfig) error {
			cpl.Info("Setting up new connection to the pool...")
			c.RuntimeParams["application_name"] = "TraverseApp"

			// set connection timeout?

			if c.Database == "" {
				return fmt.Errorf("database name is required")
			}

			return nil
		}

		// beforeAcquire hook
		cfg.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
			// try logger later too?

			// current user_name connections
			connInfo := conn.PgConn().ParameterStatus("session_authorization")

			// logger.Info("checking for existing connections...", "user", connInfo)

			// handle not in transaction
			if inTx := conn.PgConn().TxStatus() != 'I'; inTx {
				cpl.Infow("Connection in transaction, unable to acquire", "user", connInfo)
				return false
			}
			// ping database
			err = conn.Ping(ctx)
			if err != nil {
				cpl.Error("connection failed to ping database", "error", err)
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
			cpl,
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
