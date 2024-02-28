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
	Delete(ctx context.Context, fileName string, userId int) error
	ReadGzip(ctx context.Context, filePath string) ([]byte, error)
	FindByUser(ctx context.Context, userId int) ([]domain.File, error)
	GetFilePath(ctx context.Context, fileName string, userId int) string
	Create(ctx context.Context, file *multipart.FileHeader, url string, userId int) (domain.File, error)
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
	fileRouter.GET("", handler.FindByUser)

	fileNameRouter := fileRouter.Group("/:fileName")
	fileNameRouter.GET("", handler.Download)
	fileNameRouter.DELETE("", handler.Delete)

	return srv
}

func (h *FileHandler) Create(c echo.Context) error {
	user := c.Get("user").(*domain.User)
	upload, err := c.FormFile("file")
	if err != nil {
		return ErrorResp(http.StatusBadRequest, err)
	}

	prefix := "http"
	if c.IsTLS() {
		prefix = "https"
	}

	ctx := c.Request().Context()
	file, err := h.Service.Create(
		ctx,
		upload,
		fmt.Sprintf("%s://%s", prefix, c.Request().Host),
		user.ID,
	)
	if err != nil {
		return ErrorResp(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, file)
}

func (h *FileHandler) FindByUser(c echo.Context) error {
	user := c.Get("user").(*domain.User)
	ctx := c.Request().Context()

	files, err := h.Service.FindByUser(ctx, user.ID)
	if err != nil {
		return ErrorResp(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"files": files,
	})
}

func (h *FileHandler) Delete(c echo.Context) error {
	user := c.Get("user").(*domain.User)
	ctx := c.Request().Context()
	fileName := c.Param("fileName")

	if err := h.Service.Delete(ctx, fileName, user.ID); err != nil {
		return ErrorResp(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
	})
}

func (h *FileHandler) Download(c echo.Context) error {
	ctx := c.Request().Context()
	filePath := h.Service.GetFilePath(
		ctx,
		c.Param("fileName"),
		c.Get("user").(*domain.User).ID,
	)

	content, err := h.Service.ReadGzip(ctx, filePath)
	if err != nil {
		return ErrorResp(http.StatusInternalServerError, err)
	}

	return c.Blob(
		http.StatusOK,
		"application/gzip",
		content,
	)
}
