package flags

// DryRun provides mistake-safe access to the "--dry-run" Cobra command-line flag.
func DryRun() (AddFunc, ReadBoolFlagFunc) {
	return BoolPersistent("dry-run", "", "Print but do not run the Git commands")
}
