package cli

func RenderBool(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
