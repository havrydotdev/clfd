package rest

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type PublicRouter struct {
	*echo.Group
}

func NewPublicRouter(srv *echo.Echo) *PublicRouter {
	return &PublicRouter{
		srv.Group(""),
	}
}

type ProtectedRouter struct {
	*echo.Group
}

func NewProtectedRouter(srv *echo.Echo, userSvc UserService) *ProtectedRouter {
	jwtConfig := echojwt.Config{
		BeforeFunc: func(c echo.Context) {
			fmt.Println(c.Request().Header.Get("Authorization"))
		},
		SuccessHandler: func(c echo.Context) {
			token := c.Get("token").(*jwt.Token)
			claims := token.Claims.(*AccessTokenClaims)

			user, err := userSvc.FindByID(c.Request().Context(), claims.UserId)
			if err != nil {
				log.Println(err)
			}

			c.Set("user", &user)
		},
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(AccessTokenClaims)
		},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		ContextKey: "token",
	}

	protected := srv.Group("")
	protected.Use(echojwt.WithConfig(jwtConfig))

	return &ProtectedRouter{
		protected,
	}
}
