package fields

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// DateField is a custom date type
type DateField time.Time

const dateFormat = "2006-01-02"

// UnmarshalJSON parses the JSON string to DateField
func (df *DateField) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = str[1 : len(str)-1] // Remove the quotes

	t, err := time.Parse(dateFormat, str)
	if err != nil {
		return err
	}
	*df = DateField(t)
	return nil
}

// MarshalJSON formats the DateField to JSON string
func (df DateField) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(df).Format(dateFormat))
}

// Value converts DateField to a driver-compatible format
func (df DateField) Value() (driver.Value, error) {
	return time.Time(df), nil
}

// Scan converts driver value to DateField
func (df *DateField) Scan(value interface{}) error {
	if t, ok := value.(time.Time); ok {
		*df = DateField(t)
		return nil
	}
	return fmt.Errorf("cannot convert %v to DateField", value)
}
