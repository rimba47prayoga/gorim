package gorim

import (
	"github.com/labstack/echo/v4"
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


// HandlerFunc is a custom function type for handling requests
type HandlerFunc func(Context) error

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

func (s *Server) Use(middleware ...echo.MiddlewareFunc) {
	s.Echo.Use(middleware...)
}

func New() *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	server := Server{
		Echo: e,
	}
	e.Use(middlewares.RecoverMiddleware)
	return &server
}
