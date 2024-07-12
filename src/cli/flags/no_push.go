package flags

// DryRun provides mistake-safe access to the "--dry-run" Cobra command-line flag.
func NoPush() (AddFunc, ReadBoolFlagFunc) {
	return Bool("no-push", "", "Do not push local branches", FlagTypePersistent)
}
