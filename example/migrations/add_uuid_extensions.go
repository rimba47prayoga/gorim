package migrations

import "github.com/rimba47prayoga/gorim.git/conf"

func AddUuidExtensions() error {
	conf.DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	return nil
}
