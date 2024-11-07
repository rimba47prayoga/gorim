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

// GetBool returns the value associated with the key as a boolean.
func (c *Context) GetBool(key string) bool {
	if val := c.Get(key); val != nil {
		b, _ := val.(bool)
        return b
	}
	return false
}
