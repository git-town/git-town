// Package hosting provides support for interacting with forges.
// Commands like "propose", "repo", and "ship" use this package
// to know how to perform Git Town operations on GitHub, GitLab, Bitbucket, etc.
// Implementations of connectors for particular forges conform to the Connector interface.
// TODO: rename this package to forges
package forges
