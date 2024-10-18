package user

import (
	"example/db"

	"github.com/rimba47prayoga/gorim.git/models"
	"github.com/rimba47prayoga/gorim.git/views"
)


type UserViewSet struct {
	views.ModelViewSet[models.User]
}

func NewUserViewSet() *UserViewSet {
	var model models.User
	var serializer UserSerializer
	queryset := db.DB.Model(model)
	modelViewSet := views.NewModelViewSet[models.User](
		&model,
		queryset,
		&serializer,
		nil,
	)
	return &UserViewSet{
		ModelViewSet: *modelViewSet,
	}
}


// func (h *UserViewSet) Create(c echo.Context) error {
// 	serializer := UserSerializer{}
// 	serializer.SetContext(c)
// 	serializer.SetMeta(serializer.Meta())
// 	c.Bind(&serializer)
// 	validate := validator.New()
// 	err := validate.Struct(&serializer)
// 	if err != nil {
// 		fmt.Println("FALSE")
// 		return c.JSON(http.StatusBadRequest, gorim.Response{
// 			"error": err.Error(),
// 		})
// 	}
// 	return c.JSON(http.StatusCreated, gorim.Response{
// 		"msg": "success",
// 	})
// }
