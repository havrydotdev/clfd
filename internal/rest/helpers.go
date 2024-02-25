package rest

import (
	"crypto/rand"
	"io"
	"log"

	"github.com/clfdrive/server/domain"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var (
	table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
)

func ErrorResp(code int, err error) *echo.HTTPError {
	log.Println(err)

	return echo.NewHTTPError(code, err.Error())
}

func GenVerifCode() string {
	b := make([]byte, 6)
    n, err := io.ReadAtLeast(rand.Reader, b, 6)
    if n != 6 {
        panic(err)
    }

    for i := 0; i < len(b); i++ {
        b[i] = table[int(b[i])%len(table)]
    }
	
    return string(b)
}

func isRequestValid(m *domain.File) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}

	return true, nil
}