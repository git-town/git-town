package flags

// dryRunFlag provides access to the `--dry-run` flag for Cobra commands
// in a way where Go's usage checker (which produces compilation errors for unused variables)
// enforces that the programmer didn't forget to define or read the flag.
func DryRun() (AddFunc, readBoolFlagFunc) {
	return Bool("dry-run", "", "Print but do not run the Git commands")
}
