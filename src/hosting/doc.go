// Package hosting provides support for interacting with code hosting services.
// Commands like "new-pull-request", "repo", and "ship" use this package
// to know how to perform Git Town operations on GitHub, Gitlab, Bitbucket, etc.
// Drivers implement the CodeHostingDriver interface.
//
// A good starting point is the file core.go, which defines the main data structures used in this package.
package hosting
