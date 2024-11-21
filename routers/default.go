package routers

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"

	"github.com/rimba47prayoga/gorim.git"
	"github.com/rimba47prayoga/gorim.git/interfaces"
	"github.com/rimba47prayoga/gorim.git/utils"
)


type ParamPath struct {
	Name	string
	Type	string
}


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

func (r *DefaultRouter[T]) PathConverter(path string) (string, []ParamPath) {
	var paramPath []ParamPath
	// Map of supported types and their regex patterns
	typePatterns := map[string]string{
		"int":    `^[0-9]+$`,
		"uint":   `^[0-9]+$`,
		"slug":   `[a-zA-Z0-9-_]+`,
		"uuid":   `[a-fA-F0-9-]+`,
	}

	// Regex pattern to find "<type:param>" in the path
	pattern := regexp.MustCompile(`<(\w+):(\w+)>`)

	// Replace all matches with Echo-style parameters
	convertedPath := pattern.ReplaceAllStringFunc(path, func(match string) string {
		// Extract the type and parameter name
		matches := pattern.FindStringSubmatch(match)
		paramType := matches[1]
		paramName := matches[2]

		// Check if the type is supported
		if _, exists := typePatterns[paramType]; !exists {
			panic(fmt.Sprintf("Unsupported parameter type: %s", paramType))
		}

		paramPath = append(paramPath, ParamPath{
			Name: paramName,
			Type: paramType,
		})
		// Replace with Echo-style parameter (e.g., ":paramName")
		return ":" + paramName
	})

	return convertedPath, paramPath
}

// ValidateParam validates the value of a route parameter based on its type.
func (r *DefaultRouter[T]) ValidateParam(paramValue string, paramPath ParamPath) error {
	switch paramPath.Type {
		case "int":
			// Check if the value is a valid integer
			if _, err := strconv.Atoi(paramValue); err != nil {
				return fmt.Errorf("invalid int value: %s", paramValue)
			}
		case "uint":
			// Check if the value is a valid unsigned integer
			if val, err := strconv.ParseUint(paramValue, 10, 64); err != nil || val < 0 {
				return fmt.Errorf("invalid uint value: %s", paramValue)
			}
		case "uuid":
			// Regex to validate UUID format
			uuidRegex := regexp.MustCompile(`^[a-fA-F0-9\-]{36}$`)
			if !uuidRegex.MatchString(paramValue) {
				return fmt.Errorf("invalid uuid value: %s", paramValue)
			}
		case "slug":
			// Regex to validate slug format
			slugRegex := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
			if !slugRegex.MatchString(paramValue) {
				return fmt.Errorf("invalid slug value: %s", paramValue)
			}
		default:
			return fmt.Errorf("unsupported param type: %s", paramPath.Type)
	}
	return nil
}

// Helper function to handle common route logic
func(r *DefaultRouter[T]) HandleRoute(method, path, action string) {

	path, paramPaths := r.PathConverter(path)
	r.RouteGroup.Add(method, path, func(c gorim.Context) error {
		if len(paramPaths) > 0 {
			for _, paramPath := range paramPaths {
				err := r.ValidateParam(c.Param(paramPath.Name), paramPath)
				if err != nil {
					return c.JSON(http.StatusBadRequest, gorim.Response{
						"error": err.Error(),
					})
				}
			}
		}
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
	// TODO: user handler.GetPKField() to make dynamic parameter name
	handler := r.HandlerFunc()
	if utils.HasAttr(handler, "List") {
		r.HandleRoute(http.MethodGet, "", "List")
	}
	if utils.HasAttr(handler, "Create") {
        r.HandleRoute(http.MethodPost, "", "Create")
    }
    if utils.HasAttr(handler, "Retrieve") {
        r.HandleRoute(http.MethodGet, "/<uint:pk>", "Retrieve")
    }
    if utils.HasAttr(handler, "Update") {
        r.HandleRoute(http.MethodPut, "/<uint:pk>", "Update")
    }
    if utils.HasAttr(handler, "Delete") {
        r.HandleRoute(http.MethodDelete, "/<uint:pk>", "Delete")
    }
}
