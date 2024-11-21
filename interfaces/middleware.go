package interfaces

import (
	"github.com/labstack/echo/v4"
	"github.com/rimba47prayoga/gorim.git"
)

type IMiddleware interface {
	Call(gorim.Context) error
	SetNextFunc(echo.HandlerFunc)
}
