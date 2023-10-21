package flags

func Version() (AddFunc, ReadBoolFlagFunc) {
	return Bool("version", "v", "Display the version number")
}
