package flags

// Verbose provides mistake-safe access to the "--verbose" Cobra command-line flag.
func Verbose() (AddFunc, ReadBoolFlagFunc) {
	return BoolPersistent("verbose", "v", "Display all Git commands run under the hood")
}
