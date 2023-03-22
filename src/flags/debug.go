package flags

// Debug provides access to the `--debug` flag for Cobra commands
// in a way where Go's usage checker (which produces compilation errors for unused variables)
// enforces that the programmer didn't forget to define or read the flag.
func Debug() (AddFunc, readBoolFlagFunc) {
	return Bool("debug", "d", "Print all Git commands run under the hood")
}
