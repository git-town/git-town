package flags

func Version() (AddFunc, ReadBoolFlagFunc) {
	return Bool("version", "V", "Display the version number", FlagTypeNonPersistent)
}
