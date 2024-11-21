package filters

import (
	"reflect"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)


type IFilterSet interface {
	ApplyFilters(interface{}, echo.Context, *gorm.DB) *gorm.DB
}

type FilterSet struct{}

// FilteredFields returns a map of field names to their values, excluding fields that are nil, empty, or named "FilterSet".
func (fs FilterSet) FilteredFields(filter interface{}) map[string]reflect.Value {
	val := reflect.ValueOf(filter)

	// Handle if filter is a pointer to a struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	filteredFields := make(map[string]reflect.Value)

	// Ensure that val is a struct
	if val.Kind() != reflect.Struct {
		return filteredFields
	}

	// Iterate through the struct fields
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		// Skip fields specifically named "FilterSet"
		if fieldType.Type.Name() == "FilterSet" {
			continue
		}

		// Skip fields that are nil or empty
		if (fieldVal.Kind() == reflect.Ptr && fieldVal.IsNil()) || (fieldVal.Kind() == reflect.String && fieldVal.String() == "") {
			continue
		}

		// If it's not nil or empty, add to the filteredFields map
		filteredFields[fieldType.Name] = fieldVal
	}

	return filteredFields
}

// ApplyFilters applies the filters dynamically using reflection, but only for fields returned by FilteredFields.
func (fs FilterSet) ApplyFilters(filter interface{}, ctx echo.Context, query *gorm.DB) *gorm.DB {
	filteredFields := fs.FilteredFields(filter)
	if len(filteredFields) == 0 {
		// If no fields are set, return the original query without filtering
		return query
	}

	val := reflect.ValueOf(filter)
	typ := val.Elem().Type()

	for fieldName, fieldVal := range filteredFields {
		fieldType, _ := typ.FieldByName(fieldName)

		// Check if the field has a custom method tag
		methodName := fieldType.Tag.Get("method")
		if methodName != "" {
			// Find the method by name and call it if it exists
			method := val.MethodByName(methodName)
			if method.IsValid() {
				// Call the method dynamically with gin.Context and query as arguments
				results := method.Call([]reflect.Value{
					reflect.ValueOf(ctx),
					reflect.ValueOf(query),
				})
				query = results[0].Interface().(*gorm.DB)
			}
			continue
		}

		// Standard filtering logic based on the "operator" tag
		dbName := fieldType.Tag.Get("db")
		operator := fieldType.Tag.Get("operator")

		switch operator {
			case "in":
				query = query.Where(dbName+" IN (?)", fieldVal.Interface())
			case "eq":
				query = query.Where(dbName+" = ?", fieldVal.Interface())
			case "gte":
				query = query.Where(dbName+" >= ?", fieldVal.Interface())
			case "lte":
				query = query.Where(dbName+" <= ?", fieldVal.Interface())
			case "like":
				query = query.Where(dbName+" LIKE ?", "%"+fieldVal.String()+"%")
			case "ilike":
				query = query.Where(dbName+" ILIKE ?", "%"+fieldVal.String()+"%")
		}
	}

	return query
}

