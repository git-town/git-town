package azuredevops

import "github.com/git-town/git-town/v22/internal/git/giturl"

// Detect indicates whether the current repository is hosted on Azure DevOps.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "dev.azure.com" || remoteURL.Host == "ssh.dev.azure.com"
}
