package permissions

import "github.com/rimba47prayoga/gorim.git"

type IsAuthenticated struct {}

func (p *IsAuthenticated) HasPermission(ctx gorim.Context) bool {
	return ctx.GetBool("is_authenticated")
}
