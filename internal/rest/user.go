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
	Login(ctx context.Context, email string, password string) error
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
}

func (h *UserHandler) SignUp(c echo.Context) error {
	var input domain.User
	if err := c.Bind(&input); err != nil {
		return ErrorResp(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	err := h.Service.Create(ctx, &input)
	if err != nil {
		return ErrorResp(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"id": input.ID,
		"email": input.Email,
	})
}
