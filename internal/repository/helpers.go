package repository

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"go.uber.org/fx"
)

var (
	ErrNoRows = errors.New("does not exist")
)

func NewConn(lc fx.Lifecycle) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to db: ", err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal("Failed to ping db: ", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err := conn.Close(context.Background())
			if err != nil {
				log.Fatal("Failed to disconnect from db: ", err)

				return err
			}

			return nil
		},
	})

	return conn
}
