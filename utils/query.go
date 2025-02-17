package utils

import (
	"gorim.org/gorim/errors"
	"gorm.io/gorm"
)

// Generic function to get object or return 404
func GetObjectOr404[T any](queryset *gorm.DB, query string, args ...interface{}) *T {
    var result T // Create an instance of T (empty)
    
    if err := queryset.Where(query, args...).First(&result).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            errors.Raise(&errors.ObjectNotFoundError{Message: "Resource not found"})
        } else {
            errors.Raise(&errors.InternalServerError{Message: err.Error()})
        }
    }
    
    return &result // Return pointer to the result
}
