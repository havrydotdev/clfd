package middlewares

import (
	"flag"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	timeout = flag.Int("timeout", 30, "Request timeout in seconds")
)

func Use(srv *echo.Echo) *echo.Echo {
	srv.Use(middleware.ContextTimeout(time.Duration(*timeout) * time.Second))
	srv.Use(middleware.CORS())
	srv.Use(middleware.Logger())
	srv.Use(middleware.Secure())
	srv.Use(middleware.Recover())

	return srv
}
