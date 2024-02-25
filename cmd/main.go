package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/clfdrive/server/file"
	"github.com/clfdrive/server/internal/repository"
	"github.com/clfdrive/server/internal/repository/pg"
	"github.com/clfdrive/server/internal/rest"
	"github.com/clfdrive/server/internal/rest/middlewares"
	"github.com/clfdrive/server/user"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

var (
	port = flag.Int("port", 3000, "Port to start server on")
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fx.New(
		fx.Provide(
			func() *echo.Echo {
				srv := echo.New()

				return srv
			},
			repository.NewConn,
			pg.NewFileRepository,
			file.NewService,
			pg.NewUserRepository,
			user.NewService,
			rest.NewProtectedRouter,
			rest.NewPublicRouter,
		),
		fx.Invoke(
			rest.NewFileHandler,
			rest.NewUserHandler,
			func(srv *echo.Echo) {
				srv.Logger.Fatal(srv.Start(fmt.Sprintf(":%d", *port)))
			},
		),
		fx.Decorate(
			middlewares.Use,
		),
	).Run()
}
