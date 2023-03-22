package flags

// DryRun provides mistake-safe access to the "--dry-run" Cobra command-line flag.
func DryRun() (AddFunc, readBoolFlagFunc) {
	return Bool("dry-run", "", "Print but do not run the Git commands")
}
