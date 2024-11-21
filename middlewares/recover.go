package middlewares

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/rimba47prayoga/gorim.git"
	"github.com/rimba47prayoga/gorim.git/errors"
)

type Response map[string]any

type RecoverMiddleware struct {
    BaseMiddleware
}

func (m *RecoverMiddleware) Call(c gorim.Context) error {
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
                c.JSON(http.StatusInternalServerError, Response{
                    "error": r,
                })
                // Create a buffer to hold the stack trace.
                stackTrace := make([]byte, 1024)
                n := runtime.Stack(stackTrace, false)

                // Print the stack trace.
                fmt.Printf("%s\n", stackTrace[:n])
            }
        }
    }()
    return m.Next(c)
}