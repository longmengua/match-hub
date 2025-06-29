package errors

import (
	"fmt"
	"match/internal/enum"
)

var _ error = (*GrpcError)(nil)

type GrpcError struct {
	Err     error
	Code    enum.Code
	Message string
}

func (e *GrpcError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code=%d, message=%s, cause=%v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code=%d, message=%s", e.Code, e.Message)
}

func (e *GrpcError) Unwrap() error {
	return e.Err
}

// Helper to create new AppError
func New(code enum.Code, err error) *GrpcError {
	return &GrpcError{
		Code:    code,
		Message: err.Error(),
	}
}
