package rest

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/clfdrive/server/domain"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type FileService interface {
	Create(ctx context.Context, file *domain.File) error
}

type FileHandler struct {
	Service FileService
}

func isRequestValid(m *domain.File) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}

	return true, nil
}

func NewFileHandler(srv *echo.Echo, svc FileService) *echo.Echo {
	handler := &FileHandler{
		Service: svc,
	}

	srv.POST("/file", handler.Create)
	srv.GET("/file/:fileName", handler.Download)

	return srv
}

func (h *FileHandler) Create(c echo.Context) error {
	upload, err := c.FormFile("file")
	if err != nil {
		return ErrorResp(http.StatusBadRequest, err)
	}

	fileName := c.FormValue("name")
	if fileName == "" {
		parts := strings.Split(upload.Filename, ".")

		fileName = fmt.Sprintf("%s.%s", uuid.NewString(), parts[len(parts) - 1])
	}
	
	_, err = SaveFile(upload, fileName)
	if err != nil {
		return ErrorResp(http.StatusBadRequest, err)
	}

	prefix := "http"
	if c.IsTLS() {
		prefix = "https"
	}

	file := domain.File{
		Name: fileName,
		Location: fmt.Sprintf("%s://%s/file/%s", prefix, c.Request().Host, fileName),
	}

	if ok, err := isRequestValid(&file); !ok {
		return ErrorResp(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	if err := h.Service.Create(ctx, &file); err != nil {
		return ErrorResp(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, file)
}

func (h *FileHandler) Download(c echo.Context) error {
	return c.File(path.Join(driveDir, c.Param("fileName")))
}
