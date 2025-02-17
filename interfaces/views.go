package interfaces

import "gorim.org/gorim"

type IBaseView interface {
	SetAction(string)
	SetContext(gorim.Context)
	HasPermission(gorim.Context) bool
}
