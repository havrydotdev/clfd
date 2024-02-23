package rest

import (
	"context"
	"net/http"

	"github.com/clfdrive/server/domain"
	"github.com/go-playground/validator/v10"
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

func NewFileHandler(srv *echo.Echo, svc FileService) {
	handler := &FileHandler{
		Service: svc,
	}

	// srv.POST("/file", handler.)
}

func (h *FileHandler) Create(c echo.Context) error {
	var file domain.File
	if err := c.Bind(&file); err != nil {
		return ErrorResp(http.StatusBadRequest, err)
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
