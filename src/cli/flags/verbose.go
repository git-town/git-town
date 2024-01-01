package flags

// Verbose provides mistake-safe access to the "--verbose" Cobra command-line flag.
func Verbose() (AddFunc, ReadBoolFlagFunc) {
	return Bool("verbose", "v", "Display all Git commands run under the hood", FlagTypePersistent)
}
