package logs

import (
	"log"
	"strings"

	"github.com/pkg/errors"
)

// FatalOn runs log.Fatal with the given error
// when the given error exists.
// If an optional message is provided,
// the error is wrapped in it.
func FatalOn(err error, messages ...string) {
	if err == nil {
		return
	}
	if len(messages) == 0 {
		log.Fatal(err)
	} else {
		log.Fatal(errors.Wrap(err, strings.Join(messages, " ")))
	}
}

// FatalfOn runs log.Fatalf with the given error
// when the given error exists.
// If an optional message is provided,
// the error is wrapped in it.
func FatalfOn(err error, format string, formatArgs ...interface{}) {
	if err == nil {
		return
	}
	log.Fatal(errors.Wrapf(err, format, formatArgs...))
}

// FatallnOn runs log.Fatalln with the given input if the given error exists
func FatallnOn(err error, messages ...string) {
	if err == nil {
		return
	}
	if len(messages) == 0 {
		log.Fatalln(err)
	} else {
		log.Fatalln(errors.Wrap(err, strings.Join(messages, " ")))
	}
}
