package views

import (
	"github.com/rimba47prayoga/gorim.git"
)

type IBaseView interface {
	SetAction(string)
	SetContext(gorim.Context)
	HasPermission(gorim.Context) bool
}
