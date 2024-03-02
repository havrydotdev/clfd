package rest

import (
	"context"
	"net/http"

	"github.com/clfdrive/server/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AccessTokenClaims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

type UserService interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, userId int) (domain.User, error)
	Verify(ctx context.Context, email string, verifCode string) error
	SignIn(ctx context.Context, email string, password string) (string, string, error)
	Refresh(ctx context.Context, oldRefreshToken string) (string, string, error)
}

type UserHandler struct {
	Service UserService
}

func NewUserHandler(srv *echo.Echo, svc UserService, r *PublicRouter) {
	handler := &UserHandler{
		Service: svc,
	}

	auth := r.Group.Group("/auth")

	auth.POST("/sign-up", handler.SignUp)
	auth.POST("/verify", handler.Verify)
	auth.GET("/sign-in", handler.SignIn)
	auth.GET("/refresh", handler.Refresh)
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
		"id":    input.ID,
		"email": input.Email,
	})
}

type verifyInput struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (h *UserHandler) Verify(c echo.Context) error {
	var input verifyInput
	if err := c.Bind(&input); err != nil {
		return ErrorResp(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	if err := h.Service.Verify(ctx, input.Email, input.Code); err != nil {
		return ErrorResp(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
	})
}

type signInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) SignIn(c echo.Context) error {
	var input signInInput
	if err := c.Bind(&input); err != nil {
		return ErrorResp(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	accessToken, refreshToken, err := h.Service.SignIn(ctx, input.Email, input.Password)
	if err != nil {
		return ErrorResp(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

type refreshInput struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *UserHandler) Refresh(c echo.Context) error {
	var input refreshInput
	if err := c.Bind(&input); err != nil {
		return ErrorResp(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	accessToken, refreshToken, err := h.Service.Refresh(ctx, input.RefreshToken)
	if err != nil {
		return ErrorResp(http.StatusUnauthorized, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
