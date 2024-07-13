package flags

// DryRun provides mistake-safe access to the "--dry-run" Cobra command-line flag.
func All() (AddFunc, ReadBoolFlagFunc) {
	return Bool("all", "a", "Sync all local branches", FlagTypeNonPersistent)
}
