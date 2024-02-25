package rest

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/clfdrive/server/domain"
	"github.com/labstack/echo/v4"
)

type FileService interface {
	GetFileName(ctx context.Context, fileName string) string
	Create(ctx context.Context, file *multipart.FileHeader, fileName, url string) (domain.File, error)
}

type FileHandler struct {
	Service FileService
}

func NewFileHandler(srv *echo.Echo, fileSvc FileService, r *ProtectedRouter) *echo.Echo {
	handler := &FileHandler{
		Service: fileSvc,
	}

	fileRouter := r.Group.Group("/file")
	fileRouter.POST("", handler.Create)
	fileRouter.GET("/:fileName", handler.Download)

	return srv
}

func (h *FileHandler) Create(c echo.Context) error {
	upload, err := c.FormFile("file")
	if err != nil {
		return ErrorResp(http.StatusBadRequest, err)
	}

	fileName := c.FormValue("name")

	prefix := "http"
	if c.IsTLS() {
		prefix = "https"
	}

	ctx := c.Request().Context()
	file, err := h.Service.Create(ctx, upload, fileName, fmt.Sprintf("%s://%s", prefix, c.Request().Host))
	if err != nil {
		return ErrorResp(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, file)
}

func (h *FileHandler) Download(c echo.Context) error {
	ctx := c.Request().Context()

	return c.File(h.Service.GetFileName(ctx, c.Param("fileName")))
}
