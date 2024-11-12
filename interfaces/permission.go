package interfaces

import "github.com/rimba47prayoga/gorim.git"

type IPermission interface {
	HasPermission(gorim.Context) bool
}
