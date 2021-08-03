package err

import (
	"errors"
	"fmt"
)

type Err struct {
	Type    string
	Message string
	Er      error
}

func (e Err) Error() string {
	if e.Er != nil {
		return fmt.Sprintf("[%s] %s, cause: %s", e.Type, e.Message, e.Er.Error())
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

func (e Err) Unwrap() error {
	return errors.Unwrap(e.Er)
}
