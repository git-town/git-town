package cli

// RenderBool converts the given bool into either "yes" or "no".
func RenderBool(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
