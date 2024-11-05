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
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		// If the field is the embedded BaseModel, handle it
		if field.Name == "BaseModel" || field.Anonymous {
			// Handle if the embedded field is a struct or pointer to a struct
			if fieldValue.Kind() == reflect.Struct {
				// Embedded as a struct
				baseFields, err := GetModelFields(fieldValue.Interface())
				if err != nil {
					return nil, err
				}
				fields = append(fields, baseFields...)
			} else if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
				// Embedded as a pointer to a struct
				baseFields, err := GetModelFields(fieldValue.Elem().Interface())
				if err != nil {
					return nil, err
				}
				fields = append(fields, baseFields...)
			}
			continue
		}

		// For normal fields, append the field name
		fields = append(fields, field.Name)
	}

	return fields, nil
}
