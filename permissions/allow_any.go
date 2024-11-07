package permissions

import "github.com/rimba47prayoga/gorim.git"

type AllowAny struct {}

func (p *AllowAny) HasPermission(ctx gorim.Context) bool {
	return true
}
