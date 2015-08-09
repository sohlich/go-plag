package main

import (
	// "errors"
	"fmt"
)

//InitError should be used to popup
//if error occurs during
//initialization of aplication
type InitError struct {
	Message       string
	StandartError error
}

func (err *InitError) Error() string {
	return fmt.Sprintf("InitializationError: %s casued by %s", err.Message, err.StandartError.Error())
}
