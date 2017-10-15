package exit

import (
	"log"

	"github.com/pkg/errors"
)

// If runs log.Fatal with the given error
// if the given error exists.
func If(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// IfWrap runs log.Fatal with the given error
// wrapped in the given message
// if the error exists.
func IfWrap(err error, message string) {
	if err != nil {
		log.Fatal(errors.Wrap(err, message))
	}
}

// IfWrapf runs log.Fatal with the given error
// wrapped in the given message
// if the error exists.
func IfWrapf(err error, format string, formatArgs ...interface{}) {
	if err != nil {
		log.Fatal(errors.Wrapf(err, format, formatArgs...))
	}
}
