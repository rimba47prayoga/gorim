package routers

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/rimba47prayoga/gorim.git"
	"github.com/rimba47prayoga/gorim.git/interfaces"
	"github.com/rimba47prayoga/gorim.git/utils"
)


type DefaultRouter[T interfaces.IBaseView] struct {
	RouteGroup	*gorim.Group
	HandlerFunc func() T
}

func NewDefaultRouter[T interfaces.IBaseView](group *gorim.Group, handlerFunc func() T) *DefaultRouter[T] {
	router := DefaultRouter[T]{
		RouteGroup: group,
		HandlerFunc: handlerFunc,
	}
	router.AutoDiscover()
	return &router
}

func(r *DefaultRouter[T]) SetupHandler(action string, c gorim.Context) T {
	// Helper function to create and configure a handler
	handler := r.HandlerFunc()
	handler.SetAction(action)
	handler.SetContext(c)
	return handler
}

func(r *DefaultRouter[T]) RegisterFunc(name string, httpMethod string, path string) {
	r.HandleRoute(httpMethod, path, name)
}

// Helper function to handle common route logic
func(r *DefaultRouter[T]) HandleRoute(method, path, action string) {

	r.RouteGroup.Add(method, path, func(c gorim.Context) error {
		handler := r.SetupHandler(action, c)
		if !utils.HasAttr(handler, action) {
			msg := fmt.Sprintf("%s has no attribute or method %s", utils.GetStructName(handler), action)
			panic(msg)
		}
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

func (r *DefaultRouter[T]) AutoDiscover() {
	handler := r.HandlerFunc()
	if utils.HasAttr(handler, "List") {
		r.HandleRoute(http.MethodGet, "", "List")
	}
	if utils.HasAttr(handler, "Create") {
        r.HandleRoute(http.MethodPost, "", "Create")
    }
    if utils.HasAttr(handler, "Retrieve") {
        r.HandleRoute(http.MethodGet, "/:pk", "Retrieve")
    }
    if utils.HasAttr(handler, "Update") {
        r.HandleRoute(http.MethodPut, "/:pk", "Update")
    }
    if utils.HasAttr(handler, "Delete") {
        r.HandleRoute(http.MethodDelete, "/:pk", "Delete")
    }
}
