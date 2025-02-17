package migrations

import "gorim.org/gorim/conf"

func AddUuidExtensions() error {
	conf.DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	return nil
}
