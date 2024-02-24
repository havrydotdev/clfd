package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/clfdrive/server/file"
	"github.com/clfdrive/server/internal/repository/pg"
	"github.com/clfdrive/server/internal/rest"
	"github.com/clfdrive/server/internal/rest/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

var (
	port = flag.Int("port", 3000, "Port to start server on")
	timeout = flag.Int("timeout", 30, "Request timeout in seconds")
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to db: ", err)
	}
	defer func() {
		err := conn.Close(context.Background())
		if err != nil {
			log.Fatal("Failed to disconnect from db: ", err)
		}
	}()

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal("Failed to ping db: ", err)
	}

	srv := echo.New()

	srv.Use(middleware.SetRequestContextWithTimeout(time.Duration(*timeout) * time.Second))
	srv.Use(middleware.CORS)

	fileRepo := pg.NewFileRepository(conn)

	fileSvc := file.NewService(fileRepo)

	rest.NewFileHandler(srv, fileSvc)

	log.Fatal(srv.Start(fmt.Sprintf(":%d", *port)))
}
