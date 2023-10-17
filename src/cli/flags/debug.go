package flags

// DryRun provides mistake-safe access to the "--debug" Cobra command-line flag.
func Debug() (AddFunc, ReadBoolFlagFunc) {
	return Bool("debug", "d", "Print all Git commands run under the hood")
}
