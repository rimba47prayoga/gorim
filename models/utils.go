package models

import (
	"fmt"
	"reflect"
)

func GetModelFields(instance interface{}) ([]string, error) {
	val := reflect.ValueOf(instance)
	// If val is a pointer, get the element it points to
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Check if the value is a struct after dereferencing
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("instance must be a struct or a pointer to a struct")
	}
	if !val.IsValid() {
		return nil, fmt.Errorf("invalid reflect.Value")
	}

	fields := []string{}
	for i := 0; i < val.NumField(); i++ {
		if val.Type().Field(i).Name == "BaseModel" {
			baseFields, err := GetModelFields(val.Type().Field(i).Type.Elem())
			if err != nil {
				// Handle or return the error
				return nil, err // Assuming the outer function returns ([]string, error)
			}
			fields = append(fields, baseFields...)
		}
		field := val.Type().Field(i).Name
		fields = append(fields, field)
	}
	return fields, nil
}
