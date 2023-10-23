package domain

import "github.com/labstack/echo/v4"

type Jwt struct{}

type JwtUsecase interface {
	SetJwtAdmin(g *echo.Group)
	SetJwtUser(g *echo.Group)
	SetJwtGeneral(g *echo.Group)
}
