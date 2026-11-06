package envconfig

import "github.com/git-town/git-town/v24/internal/forge/forgedomain"

func GhPrompt(env EnvVars) forgedomain.GhPrompt {
	return forgedomain.GhPrompt(env.GetOpt(ghPromptDisabled).IsNone())
}
