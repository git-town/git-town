// Package hosting provides support for interacting with code hosting platforms.
// Commands like "propose", "repo", and "ship" use this package
// to know how to perform Git Town operations on GitHub, GitLab, Bitbucket, etc.
// Implementations of connectors for particular code hosting platforms conform to the Connector interface.
package hosting
