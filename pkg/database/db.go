package database

import (
	"context"
	"fmt"
	"grpc/internal/config"
	g_serv "grpc/pkg/proto"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseCG interface {
	CreatePerson(ctx context.Context, req *g_serv.PostRequest) (int64, error)
	GetPerson(ctx context.Context, id int64) (*Person, error)
}

type Person struct {
	Id      int
	User_id int
	Name    string
	Admin   bool
}

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

func (d *Database) CreatePerson(ctx context.Context, req *g_serv.PostRequest) (int64, error) {
	var person Person
	err := d.Pool.QueryRow(ctx,
		"INSERT INTO people (User_id, Name, Admin) VALUES ($1, $2, $3) RETURNING id",
		req.Info.UserId, req.Info.Name, req.Info.Admin).Scan(&person.Id)
	if err != nil {
		slog.Error("Error adding field", "ERROR", err)
		return 0, err
	}

	return int64(person.Id), nil
}
func (d *Database) GetPerson(ctx context.Context, id int64) (*Person, error) {
	var person Person
	err := d.Pool.QueryRow(ctx, "SELECT * FROM people WHERE ID=$1", id).
		Scan(&person.Id, &person.Name, &person.Admin, &person.User_id)
	if err != nil {
		slog.Error("Error while receiving data", "ERROR", err)
		return nil, err
	}
	return &person, nil
}
