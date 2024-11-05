package gorim

import (
	"github.com/labstack/echo/v4"
)

// Context is a custom context that extends Echo's Context
type Context struct {
    echo.Context
	User	interface{}
}

// NewContext creates a new Gorim context
func NewContext(echoContext echo.Context) Context {
    return Context{
        Context: echoContext,
    }
}
