package settings

import (
	"github.com/rimba47prayoga/gorim.git/conf"
)

var DATABASE = conf.Database{
	Name: "example_gorim",
	Host: "localhost",
	Port: 5432,
	User: "rimbaprayoga",
	Password: "qweqweqwe",
}
