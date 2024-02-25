package rest

import (
	"context"
	"net/http"

	"github.com/clfdrive/server/domain"
	"github.com/labstack/echo/v4"
)

type UserService interface {
	Create(ctx context.Context, user *domain.User) error
	Verify(ctx context.Context, email string, verifCode string) error
	SignIn(ctx context.Context, email string, password string) error
	UpdateRefreshToken(ctx context.Context, userId int, refreshToken string) error 
}

type UserHandler struct {
	Service UserService
}

func NewUserHandler(srv *echo.Echo, svc UserService) {
	handler := &UserHandler{
		Service: svc,
	}

	srv.POST("/sign-up", handler.SignUp)
	srv.POST("/verify/:code", handler.Verify)

	// srv.GET("/sign-in", handler.SignIn)
	// srv.GET("/refresh", handler.Refresh)
}

func (h *UserHandler) SignUp(c echo.Context) error {
	var input domain.User
	if err := c.Bind(&input); err != nil {
		return ErrorResp(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	if err := h.Service.Create(ctx, &input); err != nil {
		return ErrorResp(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"id": input.ID,
		"email": input.Email,
	})
}

func (h *UserHandler) Verify(c echo.Context) error {
	email := c.QueryParam("email")
	verifCode := c.Param("code")

	ctx := c.Request().Context()
	if err := h.Service.Verify(ctx, email, verifCode); err != nil {
		return ErrorResp(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
	})
}

// func (h *UserHandler) SignIn(c echo.Context) error {

// }

// func (h *UserHandler) Refresh(c echo.Context) error {
	
// }