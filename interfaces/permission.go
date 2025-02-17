package interfaces

import "gorim.org/gorim"

type IPermission interface {
	HasPermission(gorim.Context) bool
}
