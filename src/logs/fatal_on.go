package logs

import "log"

// FatalOn runs log.Fatal when the given error is not nil
func FatalOn(err error, v ...interface{}) {
	if err != nil {
		if len(v) == 0 {
			log.Fatal(err)
		} else {
			log.Fatal(v...)
		}
	}
}

// FatalfOn runs log.Fatalf when the given error is not nil
func FatalfOn(err error, format string, v ...interface{}) {
	if err != nil {
		log.Fatalf(format, v...)
	}
}

// FatallnOn runs log.Fatalln with the given input if the given error exists
func FatallnOn(err error, v ...interface{}) {
	if err != nil {
		if len(v) == 0 {
			log.Fatal(err)
		} else {
			log.Fatal(v...)
		}
	}
}
