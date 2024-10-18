package user

import (
	"github.com/rimba47prayoga/gorim.git"
	"github.com/rimba47prayoga/gorim.git/routers"
)

func RouterUser(group *gorim.Group) {
	routeGroup := group.Group("/users")
	userRoute := routers.NewDefaultRouter[*UserViewSet](routeGroup)
	userRoute.Register(NewUserViewSet)
}
