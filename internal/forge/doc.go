// Package forge provides support for interacting with forges.
// Commands like "propose", "repo", and "ship" use this package
// to know how to perform Git Town operations on GitHub, GitLab, Bitbucket, Gitea, Codeberg.
// Implementations of connectors for particular forges conform to the Connector interface.
package forge
