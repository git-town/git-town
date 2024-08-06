package cmdhelpers

// Long automatically compiles the long description of Cobra commands
// out of the given short summary and description.
func Long(summary string, desc ...string) string {
	if len(desc) == 1 {
		return summary + ".\n" + desc[0]
	}
	return summary + "."
}
