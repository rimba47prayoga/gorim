package permissions

import "github.com/labstack/echo/v4"


type IPermission interface {
	HasPermission(echo.Context) bool
}
