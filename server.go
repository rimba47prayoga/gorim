package gorim

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/color"
	"github.com/rimba47prayoga/gorim.git/middlewares"
)

// Server represents the Gorim server
type Server struct {
    Echo *echo.Echo
}

func (s *Server) Start(address string) error {
	return s.Echo.Start(address)
}

// Group creates a new Gorim route group
func (s *Server) Group(prefix string, middleware ...echo.MiddlewareFunc) *Group {
    g := s.Echo.Group(prefix, middleware...)
    return &Group{
		EchoGroup: g,
	} // Return a new Gorim group
}

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

// Context is a custom context that extends Echo's Context
type Context struct {
    echo.Context
}

// HandlerFunc is a custom function type for handling requests
type HandlerFunc func(Context) error

// NewContext creates a new Gorim context
func NewContext(echoContext echo.Context) Context {
    return Context{
        Context: echoContext,
    }
}

// AddRoute registers a new route with the specified method, path, and handler
func (s *Server) AddRoute(method string, path string, handler HandlerFunc) {
    s.Echo.Add(method, path, func(c echo.Context) error {
        ctx := NewContext(c) // Convert to your custom context
        return handler(ctx)   // Call the handler with your custom context
    })
}

// Override HTTP methods
func (s *Server) GET(path string, handler HandlerFunc) {
    s.AddRoute(echo.GET, path, handler)
}

func (s *Server) POST(path string, handler HandlerFunc) {
    s.AddRoute(echo.POST, path, handler)
}

func (s *Server) PUT(path string, handler HandlerFunc) {
    s.AddRoute(echo.PUT, path, handler)
}

func (s *Server) DELETE(path string, handler HandlerFunc) {
    s.AddRoute(echo.DELETE, path, handler)
}

func (s *Server) PATCH(path string, handler HandlerFunc) {
    s.AddRoute(echo.PATCH, path, handler)
}

func (s *Server) OPTIONS(path string, handler HandlerFunc) {
    s.AddRoute(echo.OPTIONS, path, handler)
}

func New() *Server {
	e := echo.New()
	e.HideBanner = true

	server := Server{
		Echo: e,
	}
	versionNumber := "v1.1"
	version := color.Red(versionNumber)
	powered := "Powered by echo.labstack.com"
	inspired := color.Green("~ Inspired By Django")

	customBanner := `
  _____ ____   ___   ____ __  ___
 / ___// __ \ / _ \ /  _//  |/  /
/ (_ // /_/ // , _/_/ / / /|_/ / 
\___/ \____//_/|_|/___//_/  /_/    %s

The Go Rest Framework for perfectionists with deadlines.
%s
%s
`
	banner := fmt.Sprintf(customBanner, version, inspired, powered)
	println(banner)
	e.Use(middlewares.RecoverMiddleware)
	return &server
}
