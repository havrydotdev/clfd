package rest

import (
	"log"

	"github.com/clfdrive/server/domain"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func ErrorResp(code int, err error) *echo.HTTPError {
	log.Println(err)

	return echo.NewHTTPError(code, err.Error())
}

func isRequestValid(m *domain.File) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}

	return true, nil
}
