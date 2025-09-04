package database

import (
	"context"
	"fmt"
	"grpc/internal/config"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func DbInit(db *config.Database) (*Database, error) {
	link := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		db.Username, db.Password, db.Host, db.Port, db.Db_name)
	pool, err := pgxpool.New(context.Background(), link)
	if err != nil {
		return nil, err
	}
	slog.Info("Database launched successfully", "Link", link)
	return &Database{Pool: pool}, nil
}
func (db *Database) Close() {
	if db.Pool != nil {
		slog.Info("Closing database connection pool...")
		db.Pool.Close()
		db.Pool = nil
		slog.Info("Database connection closed successfully")
	}
}
