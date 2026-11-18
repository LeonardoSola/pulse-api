package custom_error

import (
	"fmt"
)

type ErrorLevel int8

const (
	LevelDebug ErrorLevel = -4
	LevelInfo  ErrorLevel = 0
	LevelWarn  ErrorLevel = 4
	LevelError ErrorLevel = 8
)

type CustomError struct {
	// Status code of the error
	Status int `json:"-"`

	// Message to be shown to the user
	Message string `json:"error_message"`

	// Info to be used for debugging
	Info string `json:"error_info"`

	/* Level of the error:
	   -4 - Debug		(Logged only debug mode)
		0 - Info		(not logged)
		4 - Warning		(Logged)
		8 - Error		(Logged and alerted)
	*/
	Level ErrorLevel `json:"-"`

	Err error `json:"-"`
}

func (err CustomError) ToString() string {
	return fmt.Sprintf("%d; #%s; #%s", err.Status, err.Message, err.Info)
}

func (err CustomError) Error() string {
	if err.Err != nil {
		return err.ToString() + "; " + err.Err.Error()
	}

	return err.ToString()
}

// Unwrap exposes the underlying error to errors.Is / errors.As.
func (err CustomError) Unwrap() error {
	return err.Err
}

func (err CustomError) SetError(e error) CustomError {
	err.Err = e
	return err
}
