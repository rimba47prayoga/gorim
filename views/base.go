package views

import "github.com/labstack/echo/v4"

type IBaseView interface {
	SetAction(string)
	SetContext(echo.Context)
	HasPermission(echo.Context) bool
}
