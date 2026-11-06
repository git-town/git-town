package forgedomain

// GhPrompt indicates whether GH_PROMPT_DISABLED allows gh to prompt.
// Like gh, any value of GH_PROMPT_DISABLED disables prompts, even an empty one.
type GhPrompt bool

func (self GhPrompt) CanPrompt() bool {
	return bool(self)
}
