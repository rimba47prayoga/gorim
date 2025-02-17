package interfaces

import (
	"github.com/labstack/echo/v4"
	"gorim.org/gorim"
)

type IMiddleware interface {
	Call(gorim.Context) error
	SetNextFunc(echo.HandlerFunc)
}
