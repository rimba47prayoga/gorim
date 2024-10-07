package gorim

import (
	"github.com/labstack/echo/v4"
	"github.com/rimba47prayoga/gorim.git/middlewares"
)

func New() *echo.Echo {
	e := echo.New()
	e.Use(middlewares.RecoverMiddleware)
	return e
}
