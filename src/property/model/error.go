package model

import "fmt"

type PropertyError struct {
	Message string
}

func (e *PropertyError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}
