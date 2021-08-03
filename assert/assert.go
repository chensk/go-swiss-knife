package assert

import "github.com/chensk/go-swiss-knife/err"

func Assert(b bool, message string) error {
	if !b {
		return err.Err{Type: "assert error", Message: message}
	}
	return nil
}

func AssertFunc(f func() bool, message string) error {
	if !f() {
		return err.Err{Type: "assert error", Message: message}
	}
	return nil
}
