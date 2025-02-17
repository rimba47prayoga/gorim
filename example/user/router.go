package user

import (
	"net/http"

	"gorim.org/gorim"
	"gorim.org/gorim/routers"
)

func RouterUser(group *gorim.Group) {
	routeGroup := group.Group("/users")
	userRoute := routers.NewDefaultRouter[*UserViewSet](routeGroup, NewUserViewSet)
	userRoute.RegisterFunc("Profile", http.MethodPost, "/<uint:pk>/profile")
}
