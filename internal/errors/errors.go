package errors

import "fmt"

type ValidationError struct{
	Field interface{}
	Message string
}

func NewValidationError(field interface{} ,message string) *ValidationError{
	return &ValidationError{
		Field: field,
		Message: message,
	}
}

func (v *ValidationError) Error()string{
	return fmt.Sprintf("validation failed for field %s:,%s",v.Field,v.Message)
}


type NotFoundError struct{
	Val interface{}
	Resource string
}

func (e *NotFoundError) Error()string{
	return (fmt.Sprintf("%s not found with id %d",e.Resource,e.Val))
}

func NewNotFoundError(val interface{},resource string) *NotFoundError{
	return &NotFoundError{
		Val: val,
		Resource: resource,
	}
}

type DuplicateError struct{
	Resource interface{}
	Value string
}
func(e *DuplicateError) Error() string{
	return (fmt.Sprintf("%s already exists: %s",e.Resource,e.Value))
}
func NewDuplicateError(resource interface{},value string)*DuplicateError{
	return &DuplicateError{
		Resource: resource,
		Value: value,
	}
}