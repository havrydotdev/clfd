package rest

import (
	"log"

	"github.com/labstack/echo/v4"
)

func ErrorResp(code int, err error) *echo.HTTPError {
	log.Println(err)

	return echo.NewHTTPError(code, err.Error())
}
