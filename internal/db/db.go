package db

import (
	"context"
	"fmt"
	"log"
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

func PoolWithConfig(ctx context.Context, connString string) (*PGDB, error) {
    var err error
	logger := slog.Default().With("area", "Connection Pool")

	once.Do(func() {
		cfg, parseErr := pgxpool.ParseConfig(connString)
		if err != nil {
			err = parseErr
			return
		}

		cfg.BeforeConnect = func(ctx context.Context, c *pgx.ConnConfig) error {
			log.Printf("Setting up new connection to the pool...")
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

func (p *PGDB) BeginWithTx(ctx context.Context) (pgx.Tx, error) {
    tx, err := p.Pool.Begin(ctx)
    if err != nil {
        return nil, err
    }

    defer tx.Rollback(ctx)

    return tx, tx.Commit(ctx)

}
