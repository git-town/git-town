package logs

import (
	"log"

	"github.com/pkg/errors"
)

// FatalOn runs log.Fatal with the given error
// when the given error exists.
func FatalOn(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// FatalOnWrap runs log.Fatal with the given error
// when the given error wrapped in the given message
// when the given error exists.
func FatalOnWrap(err error, message string) {
	if err != nil {
		log.Fatal(errors.Wrap(err, message))
	}
}

// FatalOnWrapf runs log.Fatal with the given error
// when the given error wrapped in the given message
// when the given error exists.
func FatalOnWrapf(err error, format string, formatArgs ...interface{}) {
	if err != nil {
		log.Fatal(errors.Wrapf(err, format, formatArgs...))
	}
}

// FatallnOn runs log.Fatalln with the given input if the given error exists
func FatallnOn(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
