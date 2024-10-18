package utils

import (
	"fmt"
	"reflect"
)

func GetStructValue(instance interface{}, field string) (interface{}, error) {
	val := reflect.ValueOf(instance)
	if !val.IsValid() {
		return nil, fmt.Errorf("invalid reflect.Value")
	}

	// Check if instance is a pointer, if so, get the value it points to
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, nil // Return nil if the pointer is nil
		}
		val = val.Elem()
	}

	// Check if val is a struct
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct but got %s", val.Kind())
	}

	fieldVal := val.FieldByName(field)
	if !fieldVal.IsValid() {
		return nil, fmt.Errorf("no such field: %s", field)
	}

	// Handle nil pointers and interfaces
	if fieldVal.Kind() == reflect.Ptr || fieldVal.Kind() == reflect.Interface {
		if fieldVal.IsNil() {
			return nil, nil // Return nil if the field is a nil pointer or interface
		}
	}

	return fieldVal.Interface(), nil
}


// SetStructValue sets the value of a struct field
func SetStructValue(instance interface{}, field string, value interface{}) error {
	val := reflect.ValueOf(instance)
	if !val.IsValid() {
		return fmt.Errorf("invalid reflect.Value")
	}

	// Check if instance is a pointer, if so, get the value it points to
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return fmt.Errorf("nil pointer dereference")
		}
		val = val.Elem()
	}

	// Check if val is a struct
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct but got %s", val.Kind())
	}

	fieldVal := val.FieldByName(field)
	if !fieldVal.IsValid() {
		return fmt.Errorf("no such field: %s", field)
	}
	if !fieldVal.CanSet() {
		return fmt.Errorf("field %s cannot be set", field)
	}

	// Handle setting nil for pointer, slice, map, channel, and interface types
	if value == nil {
		fieldType := fieldVal.Type()
		if fieldType.Kind() == reflect.Ptr || fieldType.Kind() == reflect.Slice || fieldType.Kind() == reflect.Map || fieldType.Kind() == reflect.Chan || fieldType.Kind() == reflect.Interface {
			// Set the field to nil
			fieldVal.Set(reflect.Zero(fieldType))
			return nil
		}
		return fmt.Errorf("cannot set nil to non-nilable type %s", fieldType)
	}

	newValueVal := reflect.ValueOf(value)
	if fieldVal.Type() != newValueVal.Type() {
		return fmt.Errorf("provided value type %s doesn't match field type %s", newValueVal.Type(), fieldVal.Type())
	}

	fieldVal.Set(newValueVal)
	return nil
}


func PrintStructName(data interface{}) {
	t := reflect.TypeOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Struct {
		fmt.Println("Struct name:", t.Name())
	} else {
		fmt.Println("Not a struct")
	}
}
