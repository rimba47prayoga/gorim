package errors

type InternalServerError struct {
    Message string
}

func (e *InternalServerError) Error() string {
    return e.Message
}
