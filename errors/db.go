package errors

// ObjectNotFoundError represents a custom error type for not found objects
type ObjectNotFoundError struct {
    Message string
}

func (e *ObjectNotFoundError) Error() string {
    return e.Message
}
