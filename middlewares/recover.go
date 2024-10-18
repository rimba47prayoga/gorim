package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rimba47prayoga/gorim.git/errors"
)

type Response map[string]any

// RecoverMiddleware is the middleware that recovers from panics
func RecoverMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) (err error) {
        defer func() {
            if r := recover(); r != nil {
                // Check if the panic is of type ObjectNotFoundError
                if notFoundErr, ok := r.(*errors.ObjectNotFoundError); ok {
                    // Return 404 for ObjectNotFoundError
                    c.JSON(http.StatusNotFound, Response{
                        "error": notFoundErr.Error(),
                    })
                } else if internalServerErr, ok := r.(*errors.InternalServerError); ok {
                    c.JSON(http.StatusInternalServerError, Response{
                        "error": internalServerErr.Error(),
                    })
                } else {
                    // For other panics, return a generic 500 error
                    panic(r)
                }
            }
        }()
        return next(c)
    }
}