package exit

import (
	"log"

	"github.com/pkg/errors"
)

// On runs log.Fatal with the given error
// if the given error exists.
func On(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// OnWrap runs log.Fatal with the given error
// wrapped in the given message
// if the error exists.
func OnWrap(err error, message string) {
	if err != nil {
		log.Fatal(errors.Wrap(err, message))
	}
}

// OnWrapf runs log.Fatal with the given error
// wrapped in the given message
// if the error exists.
func OnWrapf(err error, format string, formatArgs ...interface{}) {
	if err != nil {
		log.Fatal(errors.Wrapf(err, format, formatArgs...))
	}
}
