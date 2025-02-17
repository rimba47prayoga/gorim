package permissions

import "gorim.org/gorim"

type IsAuthenticated struct {
	Message		string
	Code		int
}

func (p *IsAuthenticated) HasPermission(ctx gorim.Context) bool {
	return ctx.GetBool("is_authenticated")
}

// do response un authorized 401
func (p *IsAuthenticated) Response() {
	
}
