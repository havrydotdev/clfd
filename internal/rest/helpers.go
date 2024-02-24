package rest

import (
	"crypto/rand"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"

	"github.com/labstack/echo/v4"
)

var (
	currDir, _ = os.Getwd()
	driveDir = path.Join(currDir, "drive")

	table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
)

func ErrorResp(code int, err error) *echo.HTTPError {
	log.Println(err)

	return echo.NewHTTPError(code, err.Error())
}

func SaveFile(file *multipart.FileHeader, fileName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	loc := path.Join(driveDir, fileName)

	dst, err := os.Create(loc)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return loc, nil
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
