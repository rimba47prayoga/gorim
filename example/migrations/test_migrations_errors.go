package migrations

import "errors"

func TestMigrationError() error {
	return errors.New("chill, its just test error")
}
