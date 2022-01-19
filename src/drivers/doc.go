// Package drivers provides support for interacting with code hosting services.
// Commands like "git new-pull-request", "git-repo", and "git ship"
// use this package to know how to perform Git Town operations on GitHub, Gitlab, Bitbucket, etc.
//
// Drivers implement the CodeHostingDriver interface.
// Driver implementations are available for GitHub,
// Gitea, Bitbucket, and GitLab.
package drivers
