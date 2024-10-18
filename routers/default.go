package routers

import (
	"net/http"
	"reflect"

	"github.com/rimba47prayoga/gorim.git"
	"github.com/rimba47prayoga/gorim.git/views"
)


type DefaultRouter[T views.IBaseView] struct {
	RouteGroup	*gorim.Group
}

func NewDefaultRouter[T views.IBaseView](group *gorim.Group) *DefaultRouter[T] {
	return &DefaultRouter[T]{
		RouteGroup: group,
	}
}

func (r *DefaultRouter[T]) Register(handlerFunc func() T) {
	// Helper function to create and configure a handler
	createHandler := func(action string, c gorim.Context) T {
		handler := handlerFunc()
		handler.SetAction(action)
		handler.SetContext(c)
		return handler
	}

	// Helper function to handle common route logic
	handleRoute := func(method, path, action string) {

		r.RouteGroup.Add(method, path, func(c gorim.Context) error {
			handler := createHandler(action, c)
			// Get the method by name using reflection
			methodVal := reflect.ValueOf(handler).MethodByName(action)
			if !handler.HasPermission(c) {
				return c.JSON(http.StatusForbidden, gorim.Response{
					"error": "You are not authorized to access this resource",
				})
			}
			// Call the method with gorim.Context argument and capture return values
			result := methodVal.Call([]reflect.Value{reflect.ValueOf(c)})

			// Assuming the method returns an error as the last return value
			if len(result) > 0 {
				// Convert the last return value to error
				if errInterface := result[len(result)-1].Interface(); errInterface != nil {
					if err, ok := errInterface.(error); ok {
						// Return the error if present
						return err
					}
				}
			}
			return nil
		})
	}
	handleRoute(http.MethodGet, "", "List")
	handleRoute(http.MethodPost, "", "Create")
}
