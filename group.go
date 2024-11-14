package gorim

import "github.com/labstack/echo/v4"

// Group is a custom group that extends Echo's Group
type Group struct {
    EchoGroup *echo.Group // Embed the Echo Group
}

func (g *Group) Group(prefix string, middleware ...echo.MiddlewareFunc) *Group {
	echoGroup := g.EchoGroup.Group(prefix, middleware...)
	return &Group{
		EchoGroup: echoGroup,
	}
}

func (g *Group) Use(middleware ...echo.MiddlewareFunc) {
	g.EchoGroup.Use(middleware...)
}

// Add implements `Echo#Add()` for sub-routes within the Group.
func (g *Group) Add(method string, path string, handler HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route {
	// Combine into a new slice to avoid accidentally passing the same slice for
	// multiple routes, which would lead to later add() calls overwriting the
	// middleware from earlier calls.
	return g.EchoGroup.Add(method, path, func(c echo.Context) error {
		ctx := NewContext(c)
		return handler(ctx)
	})
}
