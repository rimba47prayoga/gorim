package permissions

import "gorim.org/gorim"

type AllowAny struct {}

func (p *AllowAny) HasPermission(ctx gorim.Context) bool {
	return true
}
