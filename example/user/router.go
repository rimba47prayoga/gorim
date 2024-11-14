package user

import (
	"net/http"

	"github.com/rimba47prayoga/gorim.git"
	"github.com/rimba47prayoga/gorim.git/routers"
)

func RouterUser(group *gorim.Group) {
	routeGroup := group.Group("/users")
	userRoute := routers.NewDefaultRouter[*UserViewSet](routeGroup, NewUserViewSet)
	userRoute.RegisterFunc("UpdateProfile", http.MethodPost, "/profile")
}
