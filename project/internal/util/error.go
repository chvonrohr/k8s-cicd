package util

import "log"

// FailOnError logs the given message fatally and exits if the error passed to the function is non-nil.
// This can be used to deal with functions returning a single error in one line.
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// FailOnErrorF functions similarly to FailOnError, however it accepts a function and only calls this function when
// FailOnErrorF is called. This allows the passed function to be deferred to a later point.
func FailOnErrorF(ef func() error, msg string) {
	err := ef()
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
