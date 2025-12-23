package custom_error

import (
	"fmt"
	"log"
)

type ErrorLevel int8

const (
	Debug   ErrorLevel = -1 // Log only debug mode
	Info    ErrorLevel = 0  // Not log
	Warning ErrorLevel = 1  // Log
	Error   ErrorLevel = 2  // Log and alert
)

type CustomError struct {
	// Status code of the error
	Status int `json:"-"`

	// Message to be shown to the user
	Message string `json:"error_message"`

	// Info to be used for debugging
	Info string `json:"error_info"`

	/* Level of the error:
	   -1 - Debug		(Logged only debug mode)
		0 - Info		(not logged)
		1 - Warning		(Logged)
		2 - Error		(Logged and alerted)
	*/
	Level ErrorLevel `json:"-"`

	Err error `json:"-"`
}

func (err CustomError) ToString() string {
	return fmt.Sprintf("%d; #%s; #%s", err.Status, err.Message, err.Info)
}

func (err CustomError) Error() string {
	errMessage := fmt.Sprintf("Error: %s", err.ToString())

	if err.Err != nil {
		errMessage = "; " + err.Err.Error()
	}

	log.Println(errMessage)

	return err.ToString()
}

func (err CustomError) SetError(e error) CustomError {
	err.Err = e
	return err
}
