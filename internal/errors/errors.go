package errors

import "fmt"

type RequestError struct {
	StatusCode int
	Message    string
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Message)
}
