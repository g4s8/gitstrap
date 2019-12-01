package gitstrap

import (
	"fmt"
)

// @todo #37:30min Move error handling to separate module and publish it
//  or find external dependency to work with errors.
type wrapper struct {
	cause error
	msg   string
}

func (w *wrapper) Error() string {
	return fmt.Sprintf("%s; caused by:\n\t%s", w.msg, w.cause)
}

// ErrCompose - add new error to origin
func ErrCompose(e error, msg string) error {
	return &wrapper{e, msg}
}
